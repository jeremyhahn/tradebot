package util

import "time"

func Every(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}

func EveryIf(d time.Duration, f func(), condition bool) {
	for range time.Tick(d) {
		if condition {
			f()
		} else {
			return
		}
	}
}
