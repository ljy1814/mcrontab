package main

import "time"

const (
	TimeFormat = "20060102 15:04:05"
)

func GetNowString() string {
	return time.Now().Format(TimeFormat)
}

func GetTimeString(t time.Time) string {
	return t.Format(TimeFormat)
}

func GetJobLockKey(name string) string {
	return JOB_PREFIX_LOCK + name
}

func GetJobCreateKey(name string) string {
	return JOB_PREFIX_PUT + name
}
