package main

import (
	"runtime"
	"strconv"
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

func GetJobKillKey(name string) string {
	return JOB_PREFIX_KILL + name
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

func funcName(skip int) (name string) {
	if _, file, lineNo, ok := runtime.Caller(skip); ok {
		return file + ":" + strconv.Itoa(lineNo)
	}

	return "unknown:0"
}

func errIncr(lv Level, source string) {
	if lv == _errorLevel {
		// 统计
	}
}
