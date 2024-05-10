package utils

import "time"

func NanoUnixToTime(nano int64) time.Time {
	sec := nano / 1e9
	nsec := nano % 1e9
	return time.Unix(sec, nsec)
}

