package util

import (
	"fmt"
	"time"
)

func FormatDate(date time.Time) string {
	return fmt.Sprintf("%d-%d-%d:%d:%d:%d", date.Month(),
		date.Day(), date.Year(), date.Hour(), date.Minute(), date.Second())
}
