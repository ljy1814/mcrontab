package main

const (
	JOB_EVENT_SAVE   = 1
	JOB_EVENT_DELETE = 2
	JOB_EVENT_KILL   = 4
)

const (
	JOB_PREFIX_PUT  = "/jobs/put"
	JOB_PREFIX_KILL = "/jobs/kill/"
	JOB_PREFIX_LOCK = "/jobs/lock/"
)
