package log

import "time"

func Now() string {
	return time.Now().Format("[15:04:05]")
}