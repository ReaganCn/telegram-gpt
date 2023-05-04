package utils

import "time"

/* Function that returns the current time in milliseconds */
func GetTimeInMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
