package utils

import "time"

func CallInterval(f func(), interval time.Duration) {
	for {
		f()
		time.Sleep(interval)
	}
}
