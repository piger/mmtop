// mysql> describe INFORMATION_SCHEMA.PROCESSLIST;
// +---------+---------------------+------+-----+---------+-------+
// | Field   | Type                | Null | Key | Default | Extra |
// +---------+---------------------+------+-----+---------+-------+
// | ID      | bigint(21) unsigned | NO   |     | 0       |       |
// | USER    | varchar(32)         | NO   |     |         |       |
// | HOST    | varchar(64)         | NO   |     |         |       |
// | DB      | varchar(64)         | YES  |     | NULL    |       |
// | COMMAND | varchar(16)         | NO   |     |         |       |
// | TIME    | int(7)              | NO   |     | 0       |       |
// | STATE   | varchar(64)         | YES  |     | NULL    |       |
// | INFO    | longtext            | YES  |     | NULL    |       |
// +---------+---------------------+------+-----+---------+-------+
// 8 rows in set (0.00 sec)

// mysql> SELECT ID, USER, HOST, DB, COMMAND, TIME, STATE, INFO FROM INFORMATION_SCHEMA.PROCESSLIST;
// +----+------+-----------+-------+---------+-------+-----------+-------------------------------------------------------------------------------------------+
// | ID | USER | HOST      | DB    | COMMAND | TIME  | STATE     | INFO                                                                                      |
// +----+------+-----------+-------+---------+-------+-----------+-------------------------------------------------------------------------------------------+
// | 26 | root | localhost | mysql | Query   |     0 | executing | SELECT ID, USER, HOST, DB, COMMAND, TIME, STATE, INFO FROM INFORMATION_SCHEMA.PROCESSLIST |
// |  4 | root | localhost | NULL  | Sleep   | 33363 |           | NULL                                                                                      |
// +----+------+-----------+-------+---------+-------+-----------+-------------------------------------------------------------------------------------------+
// 2 rows in set (0.00 sec)

package sqlgo

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
	"time"
)

type Connection struct {
	data    DbConnectionData
	db      *sql.DB
	in      chan string
	result  chan ProcessList
	control chan bool
}

type ProcessList struct {
	Conn      Connection
	Processes []Process
}

type Process struct {
	Id        int
	User      string
	Host      string
	Db        string
	Command   string
	Timestamp int
	State     string
	Info      string
}

func MayEmptyString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func (p Process) String() string {
	return fmt.Sprintf("id=%d, user=%s, host=%s, db=%s, command=%s, timestamp=%d, state=%s, info=%s", p.Id, p.User, p.Host, p.Db, p.Command, p.Timestamp, p.State, p.Info)
}

func QueryRunner(conn Connection, logc chan string) {
	t := time.NewTimer(5 * time.Second)

	RunProcessList(conn, logc)
	for {
		select {
		case <-conn.control:
			return
		case <-t.C:
			RunProcessList(conn, logc)
			t.Reset(5 * time.Second)
		}
	}
}

func RunProcessList(conn Connection, logc chan string) {
	rows, err := conn.db.Query("SHOW FULL PROCESSLIST")
	if err != nil {
		logc <- fmt.Sprintf("ps query failed on %s: %s", conn.data.Name, err)
	} else {
		var processes []Process
		defer rows.Close()
		for rows.Next() {
			var id, timestamp int
			var user, host, command string
			var db, state, info sql.NullString
			if err := rows.Scan(&id, &user, &host, &db, &command, &timestamp, &state, &info); err != nil {
				logc <- fmt.Sprintf("Failed parsing of result row from %s: %s", conn.data.Name, err)
			} else {
				// do something with data
				p := Process{id, user, host, MayEmptyString(db), command, timestamp, MayEmptyString(state), MayEmptyString(info)}
				logc <- p.String()
				processes = append(processes, p)
			}
		}
		if err := rows.Err(); err != nil {
			logc <- fmt.Sprintf("Error after reading rows from %s: %s", conn.data.Name, err)
		}
		plist := ProcessList{conn, processes}
		conn.result <- plist
	}
}

func RunClient(configs []DbConnectionData) {
	var clients map[string]Connection = make(map[string]Connection)

	var inConns chan DbConnectionData = make(chan DbConnectionData)
	var outConns chan Connection = make(chan Connection)
	var resc chan ProcessList = make(chan ProcessList)
	var logc chan string = make(chan string)
	var cmdc chan string = make(chan string)
	var quitc chan bool = make(chan bool)

	go Connector(inConns, outConns, resc, logc)

	fmt.Printf("Feeding Connector with connections...\n")
	for _, config := range configs {
		fmt.Printf("Feeding Connector with %s\n", config)
		inConns <- config
	}

	// go PromptReader(cmdc, logc, quitc)
	go LogPrinter(logc)
	go ResultDisplayer(resc, quitc)

	for {
		select {
		case client := <-outConns:
			clients[client.data.Name] = client
			go QueryRunner(client, logc)
		case <-cmdc:
			// nop
		case <-quitc:
			fmt.Printf("Exiting...\n")
			return
		}
	}
}

func Connector(in chan DbConnectionData, out chan Connection, resc chan ProcessList, logc chan string) {
	for {
		select {
		case d := <-in:
			go func(data DbConnectionData) {
				for {
					db, err := sql.Open("mysql", fmt.Sprintf("%s/mysql", data))
					if err != nil {
						logc <- fmt.Sprintf("Error creating DB connection object: %s", err)
						return
					}
					if err := db.Ping(); err == nil {
						conn := Connection{
							data:    data,
							db:      db,
							in:      make(chan string),
							result:  resc,
							control: make(chan bool),
						}
						out <- conn
						return
					}
					logc <- fmt.Sprintf("Connection to %s failed; sleeping 5 seconds", data.Name)
					time.Sleep(time.Duration(5 * time.Second))
				}
			}(d)
		}
	}
}

func PromptReader(cmdc chan string, logc chan string, quitc chan bool) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			// EOF ?
			logc <- fmt.Sprintf("Error reading stdin: %s", err)
			quitc <- true
			return
		} else {
			if strings.Contains(text, "quit") || strings.Contains(text, "exit") {
				quitc <- true
				return
			}
			cmdc <- text
		}
	}
}

func LogPrinter(logc chan string) {
	for {
		select {
		// case msg := <-logc:
		// 	fmt.Printf("log: %s\n", msg)
		case <-logc:
		}
	}
}

func ResultDisplayer(resc chan ProcessList, quitc chan bool) {
	var status map[string]ProcessList = make(map[string]ProcessList)
	t := time.NewTimer(5 * time.Second)

	var in chan map[string]ProcessList = make(chan map[string]ProcessList)
	go DisplayResults(in, quitc)

	for {
		select {
		case result := <-resc:
			_, notNew := status[result.Conn.data.Name]
			status[result.Conn.data.Name] = result
			if !notNew {
				in <- status
			}
		case <-t.C:
			in <- status
			t.Reset(5 * time.Second)
		}
	}
}
