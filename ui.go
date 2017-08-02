package sqlgo

import (
	"github.com/nsf/termbox-go"
	"strconv"
)

func DisplayResults(in chan map[string]ProcessList, control chan bool) {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()
	defer func() {
		control <- true
	}()

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

	for {
		select {
		case status := <-in:
			// display stuff
			termbox.Clear(termbox.ColorDefault, termbox.ColorBlack)

			// calculate header lengths
			sName := len("Name")
			sId := len("Id")
			sUser := len("User")
			sHost := len("Host")
			sDb := len("Db")
			sCommand := len("Command")
			sTime := len("Time")
			sState := len("State")
			for k, v := range status {
				var l int
				l = len(k)
				if l > sName {
					sName = l
				}
				for _, process := range v.Processes {
					id := strconv.Itoa(process.Id)
					l = len(id)
					if l > sId {
						sId = l
					}

					l = len(process.User)
					if l > sUser {
						sUser = l
					}

					l = len(process.Host)
					if l > sHost {
						sHost = l
					}

					l = len(process.Db)
					if l > sDb {
						sDb = l
					}

					l = len(process.Command)
					if l > sCommand {
						sCommand = l
					}

					timestamp := strconv.Itoa(process.Timestamp)
					l = len(timestamp)
					if l > sTime {
						sTime = l
					}

					l = len(process.State)
					if l > sState {
						sState = l
					}
				}
			}

			// draw header
			width, _ := termbox.Size()
			x := 0
			y := 0
			termbox.SetCell(x, y, '+', termbox.ColorDefault, termbox.ColorBlack)

			for x = 1; x < width-1; x++ {
				termbox.SetCell(x, y, '-', termbox.ColorDefault, termbox.ColorBlack)
			}

			termbox.SetCell(width-1, y, '+', termbox.ColorDefault, termbox.ColorBlack)

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
			termbox.SetCell(x, y, '+', termbox.ColorDefault, termbox.ColorBlack)
			for x = 1; x < width-1; x++ {
				termbox.SetCell(x, y, '-', termbox.ColorDefault, termbox.ColorBlack)
			}
			termbox.SetCell(width-1, y, '+', termbox.ColorDefault, termbox.ColorBlack)

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

			termbox.Flush()

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
