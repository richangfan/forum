package tool

import "time"

func GetCurrentDateTime() string {
	return time.Now().String()[0:19]
}
