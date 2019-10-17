package main

import (
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"time"

	"github.com/Sirupsen/logrus"
)

type Executor struct {
	ScheduleResultChan chan *ExecResult
}

func (e *Executor) PushJobResult(result *ExecResult) {
	e.ScheduleResultChan <- result
}

type ExecResult struct {
	Err          error
	Output       []byte
	JobPlan      *SchedulePlan
	StartTime    time.Time
	EndTime      time.Time
	ScheduleTime time.Time
}

func (e ExecResult) String() string {
	return fmt.Sprintf("{err:'%v',output:'%s',jobPlan:%s,startTime:'%s',endTime:'%s',scheduleTime:'%s'}",
		e.Err, e.Output, e.JobPlan,
		GetTimeString(e.StartTime),
		GetTimeString(e.EndTime),
		GetTimeString(e.ScheduleTime),
	)
}

func (e *Executor) ExecJob(plan *SchedulePlan) error {
	return e.execJob(plan)
}

func (e *Executor) execJob(plan *SchedulePlan) error {
	fun := "Executor.execJob -->"

	// 获取锁
	jobLock := GJobMgr.NewJobLock(plan.Job.Name)
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)

	ctx := context.Background()
	err := jobLock.TryLock(ctx, GetJobLockKey(plan.Job.Name))

	// 执行完释放锁
	defer jobLock.UnLock(ctx)

	now := time.Now()
	result := &ExecResult{
		Err:     err,
		EndTime: now,
		JobPlan: plan,
	}

	if err != nil {
		logrus.Errorf("%s TryLock failed plan:%s err:%v", fun, plan, err)
	} else {
		logrus.Infof("%s TryLock successfully plan:%s now:%s err:%v", fun, plan, GetNowString(), err)

		cmd := exec.CommandContext(plan.ctx, "/bin/bash", "-c", plan.Job.Command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			logrus.Warnf("%s exec failed plan:%s err:%v", fun, plan, err)
		}

		result.Err = err
		result.Output = output
		result.StartTime = now
		result.EndTime = time.Now()
		result.ScheduleTime = plan.NextTime
	}

	e.PushJobResult(result)

	return nil
}
