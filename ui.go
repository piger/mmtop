package sqlgo

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"strconv"
	"time"
)

func DisplayResults(resc chan ProcessList, control chan bool, logc chan LogMsg) {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	// termbox.Close() must be called before we signal the main thread to exit.
	defer func() {
		control <- true
	}()
	defer termbox.Close()

	var uiQuitc chan bool = make(chan bool)
	var logs []LogMsg
	var maxLogs int = 4

	go func() {
		for {
			event := termbox.PollEvent()
			if event.Type == termbox.EventKey {
				if event.Ch == 'q' || event.Key == termbox.KeyEsc {
					uiQuitc <- true
				}
			}
		}
	}()

	// timer to refresh display
	t := time.NewTimer(5 * time.Second)

	// keep track of existing process lists
	var status map[string]ProcessList = make(map[string]ProcessList)

	for {
		select {
		case <-t.C:
			DrawProcessList(status, logs)
			t.Reset(5 * time.Second)
		case result := <-resc:
			status[result.Conn.data.Name] = result
			DrawProcessList(status, logs)
		case msg := <-logc:
			logs = append(logs, msg)
			if len(logs) > maxLogs {
				d := len(logs) - maxLogs
				logs = logs[d:]
			}
			DrawProcessList(status, logs)
		case <-uiQuitc:
			return
		}
	}
}

func DrawWord(s string, x, y int, fg, bg termbox.Attribute) (count int) {
	maxWidth, _ := termbox.Size()
	for _, c := range s {
		termbox.SetCell(x, y, c, fg, bg)
		x += 1
		count += 1
		if x >= maxWidth {
			return
		}
	}
	return
}

// Draw a box line, like "+----------------+"
func DrawBoxLine(x, y, length int, fg, bg termbox.Attribute) {
	termbox.SetCell(x, y, '+', fg, bg)
	for i := x + 1; i < length-2; i++ {
		termbox.SetCell(i, y, '-', fg, bg)
	}
	termbox.SetCell(x+length-1, y, '+', fg, bg)
}

func DrawTableCell(s string, maxlength, x, y int, fg, bg termbox.Attribute) int {
	termbox.SetCell(x, y, '|', fg, bg)
	x += 2
	x += DrawWord(s, x, y, fg, bg)
	x += (maxlength - len(s)) + 1
	return x
}

func DrawProcessList(status map[string]ProcessList, logs []LogMsg) {
	fg := termbox.ColorDefault
	bg := termbox.ColorDefault

	// clear the screen first
	termbox.Clear(fg, bg)

	// calculate header lengths
	sName := len("Name")
	sId := len("Id")
	sUser := len("User")
	sHost := len("Host")
	sDb := len("Db")
	sCommand := len("Command")
	sTime := len("Time")
	sState := len("State")
	for name, v := range status {
		sName = max(len(name), sName)
		for _, process := range v.Processes {
			id := strconv.Itoa(process.Id)
			sId = max(len(id), sId)
			sUser = max(len(process.User), sUser)
			sHost = max(len(process.Host), sHost)
			sDb = max(len(process.Db), sDb)
			sCommand = max(len(process.Command), sCommand)
			timestamp := strconv.Itoa(process.Timestamp)
			sTime = max(len(timestamp), sTime)
			sState = max(len(process.State), sState)
		}
	}

	// start drawing
	width, height := termbox.Size()
	x := 0
	y := 0

	// draw first header box line
	DrawBoxLine(x, y, width, fg, bg)

	x = 0
	y += 1
	x = DrawTableCell("Hostname", sName, x, y, fg, bg)
	x = DrawTableCell("Id", sId, x, y, fg, bg)
	x = DrawTableCell("User", sUser, x, y, fg, bg)
	x = DrawTableCell("Host", sHost, x, y, fg, bg)
	x = DrawTableCell("Db", sDb, x, y, fg, bg)
	x = DrawTableCell("Command", sCommand, x, y, fg, bg)
	x = DrawTableCell("Time", sTime, x, y, fg, bg)
	x = DrawTableCell("State", sState, x, y, fg, bg)
	x = DrawTableCell("Info", len("Info"), x, y, fg, bg)

	// draw the last '|' of the table
	termbox.SetCell(width-1, y, '|', fg, bg)

	// draw the bottom part of the header
	y += 1
	x = 0
	DrawBoxLine(x, y, width, fg, bg)

	for hostname, plist := range status {
		for _, process := range plist.Processes {
			x = 0
			y += 1

			x = DrawTableCell(hostname, sName, x, y, fg, bg)
			id := strconv.Itoa(process.Id)
			x = DrawTableCell(id, sId, x, y, fg, bg)
			x = DrawTableCell(process.User, sUser, x, y, fg, bg)
			x = DrawTableCell(process.Host, sHost, x, y, fg, bg)
			x = DrawTableCell(process.Db, sDb, x, y, fg, bg)
			x = DrawTableCell(process.Command, sCommand, x, y, fg, bg)
			t := strconv.Itoa(process.Timestamp)
			x = DrawTableCell(t, sTime, x, y, fg, bg)
			x = DrawTableCell(process.State, sState, x, y, fg, bg)
			x = DrawTableCell(process.Info, width-x-1, x, y, fg, bg)

			// draw the last '|' of the table
			termbox.SetCell(width-1, y, '|', fg, bg)
		}
	}

	x = 0
	y += 1
	DrawBoxLine(x, y, width, fg, bg)

	// draw log box
	if len(logs) > 0 {
		x = 0
		y = height - 5
		DrawBoxLine(x, y, width, fg, bg)

		x = 0
		y += 1
		for _, logMsg := range logs {
			line := fmt.Sprintf("%s %s", logMsg.t.Format("15:04:05"), logMsg.m)
			DrawTableCell(line, width-2, x, y, fg, bg)
			termbox.SetCell(width-1, y, '|', fg, bg)
			y += 1
		}
	}

	termbox.Flush()
}
