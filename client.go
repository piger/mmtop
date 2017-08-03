package sqlgo

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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

func NewProcess(id int, user, host string, db sql.NullString, command string, timestamp int, state, info sql.NullString) Process {
	return Process{
		Id:        id,
		User:      user,
		Host:      host,
		Db:        MayEmptyString(db),
		Command:   command,
		Timestamp: timestamp,
		State:     MayEmptyString(state),
		Info:      MayEmptyString(info),
	}
}

func MayEmptyString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

// Run "show full processlist" every 5 seconds
func QueryRunner(conn Connection, logc chan LogMsg) {
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

func RunProcessList(conn Connection, logc chan LogMsg) {
	rows, err := conn.db.Query("SHOW FULL PROCESSLIST")
	if err != nil {
		logc <- NewLog("ps query failed on %s: %s", conn.data.Name, err)
	} else {
		var processes []Process
		defer rows.Close()
		for rows.Next() {
			var id, timestamp int
			var user, host, command string
			var db, state, info sql.NullString
			if err := rows.Scan(&id, &user, &host, &db, &command, &timestamp, &state, &info); err != nil {
				logc <- NewLog("Failed parsing of result row from %s: %s", conn.data.Name, err)
			} else {
				// do something with data
				p := NewProcess(id, user, host, db, command, timestamp, state, info)
				processes = append(processes, p)
			}
		}
		if err := rows.Err(); err != nil {
			logc <- NewLog("Error after reading rows from %s: %s", conn.data.Name, err)
		}
		plist := ProcessList{conn, processes}
		conn.result <- plist
	}
}

// Main entry point
func RunClient(configs []DbConnectionData) {
	var clients map[string]Connection = make(map[string]Connection)

	var inConns chan DbConnectionData = make(chan DbConnectionData)
	var outConns chan Connection = make(chan Connection)
	var resc chan ProcessList = make(chan ProcessList)
	var logc chan LogMsg = make(chan LogMsg)
	var quitc chan bool = make(chan bool)

	// feed the Connector with connection requests
	go Connector(inConns, outConns, resc, logc)
	for _, config := range configs {
		inConns <- config
	}

	go DisplayResults(resc, quitc, logc)

	for {
		select {
		// Listens for established connections from the Connector (via 'outConns') and starts a QueryRunner
		// for that Connection.
		case client := <-outConns:
			clients[client.data.Name] = client
			go QueryRunner(client, logc)
		case <-quitc:
			fmt.Printf("Exiting...\n")
			return
		}
	}
}

// Listens on 'in' for incoming "connection requests" (as DbConnectionData objects), attempts to
// connect to the database and in case of success write a Connection object into 'out'.
// Will log via 'logc'.
func Connector(in chan DbConnectionData, out chan Connection, resc chan ProcessList, logc chan LogMsg) {
	for {
		select {
		case d := <-in:
			go func(data DbConnectionData) {
				for {
					db, err := sql.Open("mysql", fmt.Sprintf("%s/mysql", data))
					if err != nil {
						logc <- NewLog("Error creating DB connection object: %s", err)
						return
					}
					if err := db.Ping(); err != nil {
						logc <- NewLog("Connection to %s failed; sleeping 5 seconds", data.Name)
						time.Sleep(time.Duration(5 * time.Second))
					} else {
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
				}
			}(d)
		}
	}
}
