package main

import "time"

const (
	TimeFormat = "20060102 15:04:05"
)

func GetNowString() string {
	return time.Now().Format(TimeFormat)
}
