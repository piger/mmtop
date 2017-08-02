package sqlgo

import (
	"github.com/nsf/termbox-go"
	"math"
	"strconv"
	"time"
)

func max(x, y int) int {
	return int(math.Max(float64(x), float64(y)))
}

func DisplayResults(resc chan ProcessList, control chan bool) {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	// termbox.Close() must be called before we signal the main thread to exit.
	defer func() {
		control <- true
	}()
	defer termbox.Close()

	var uiQuitc chan bool = make(chan bool)

	go func() {
		for {
			event := termbox.PollEvent()
			if event.Type == termbox.EventKey {
				if event.Ch == 'q' {
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
			DrawProcessList(status)
			t.Reset(5 * time.Second)
		case result := <-resc:
			status[result.Conn.data.Name] = result
			DrawProcessList(status)
		case <-uiQuitc:
			return
		}
	}
}

func DisplayWord(s string, x, y int) (count int) {
	maxWidth, _ := termbox.Size()
	for _, c := range s {
		termbox.SetCell(x, y, c, termbox.ColorDefault, termbox.ColorBlack)
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

func DrawProcessList(status map[string]ProcessList) {
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
	width, _ := termbox.Size()
	x := 0
	y := 0

	// draw first header box line
	DrawBoxLine(x, y, width, fg, bg)

	x = 0
	y += 1
	termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
	x += 2
	x += DisplayWord("Hostname", x, y)
	x += (sName - len("Hostname")) + 1

	termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
	x += 2
	x += DisplayWord("Id", x, y)
	x += (sId - len("Id")) + 1

	termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
	x += 2
	x += DisplayWord("User", x, y)
	x += (sUser - len("User")) + 1

	termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
	x += 2
	x += DisplayWord("Host", x, y)
	x += (sHost - len("Host")) + 1

	termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
	x += 2
	x += DisplayWord("Db", x, y)
	x += (sDb - len("Db")) + 1

	termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
	x += 2
	x += DisplayWord("Command", x, y)
	x += (sCommand - len("Command")) + 1

	termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
	x += 2
	x += DisplayWord("Time", x, y)
	x += (sTime - len("Time")) + 1

	termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
	x += 2
	x += DisplayWord("State", x, y)
	x += (sState - len("State")) + 1

	termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
	x += 2
	x += DisplayWord("Info", x, y)

	termbox.SetCell(width-1, y, '|', termbox.ColorDefault, termbox.ColorBlack)

	// draw the bottom part of the header
	y += 1
	x = 0
	DrawBoxLine(x, y, width, fg, bg)

	for hostname, plist := range status {
		for _, process := range plist.Processes {
			x = 0
			y += 1

			termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
			x += 2
			x += DisplayWord(hostname, x, y)
			x += (sName - len(hostname)) + 1

			termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
			x += 2
			id := strconv.Itoa(process.Id)
			x += DisplayWord(id, x, y)
			x += (sId - len(id)) + 1

			termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
			x += 2
			x += DisplayWord(process.User, x, y)
			x += (sUser - len(process.User)) + 1

			termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
			x += 2
			x += DisplayWord(process.Host, x, y)
			x += (sHost - len(process.Host)) + 1

			termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
			x += 2
			x += DisplayWord(process.Db, x, y)
			x += (sDb - len(process.Db)) + 1

			termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
			x += 2
			x += DisplayWord(process.Command, x, y)
			x += (sCommand - len(process.Command)) + 1

			termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
			x += 2
			t := strconv.Itoa(process.Timestamp)
			x += DisplayWord(t, x, y)
			x += (sTime - len(t)) + 1

			termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
			x += 2
			x += DisplayWord(process.State, x, y)
			x += (sState - len(process.State)) + 1

			termbox.SetCell(x, y, '|', termbox.ColorDefault, termbox.ColorBlack)
			x += 2
			x += DisplayWord(process.Info, x, y)

			termbox.SetCell(width-1, y, '|', termbox.ColorDefault, termbox.ColorBlack)
		}
	}

	x = 0
	y += 1
	DrawBoxLine(x, y, width, fg, bg)

	termbox.Flush()
}
