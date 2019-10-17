package main

import (
	"time"
)

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

func ExtractCreateJobName(key string) string {
	if len(key) <= len(JOB_PREFIX_PUT) {
		return key
	}
	return key[len(JOB_PREFIX_PUT):]
}

func ExtractKillJobName(key string) string {
	if len(key) <= len(JOB_PREFIX_KILL) {
		return key
	}
	return key[len(JOB_PREFIX_KILL):]
}
