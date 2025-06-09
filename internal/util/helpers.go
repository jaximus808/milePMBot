package util

import "time"

type HandleReport struct {
	success   bool
	timestamp time.Time
	info      string
}

func CreateHandleReport(success bool, info string) *HandleReport {
	return &HandleReport{
		success:   success,
		timestamp: time.Now(),
		info:      info,
	}
}
