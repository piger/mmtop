package sqlgo

import (
	"fmt"
	"math"
	"time"
)

// Log messages to be displayed in the UI
type LogMsg struct {
	t time.Time
	m string
}

func NewLog(format string, args ...interface{}) LogMsg {
	return LogMsg{time.Now(), fmt.Sprintf(format, args...)}
}

func max(x, y int) int {
	return int(math.Max(float64(x), float64(y)))
}
