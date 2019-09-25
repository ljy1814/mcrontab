package main

import (
	"context"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorhill/cronexpr"
)

type Scheduler struct {
	//
	jobEventChan chan *JobEvent
	// job执行计划表
	jobEventPlanMap map[string]*SchedulePlan
	// 正在执行的job表
	jobExecMap map[string]*SchedulePlan
}

func (s *Scheduler) scheduleLoop() {
	var (
		timeAfter      = s.trySchedule()
		schedulerTimer = time.NewTimer(timeAfter)
	)
	for {
		select {
		case v := <-s.jobEventChan:
			logrus.Info(v)
			// 任务更新
		case <-schedulerTimer.C:
			// 休眠结束
		case result := <-GExecutor.ScheduleResultChan:
			// 任务执行完毕
			logrus.Info(result)
		}

		timeAfter = s.trySchedule()
		schedulerTimer.Reset(timeAfter)
	}
}

func (s *Scheduler) trySchedule() (timeAfter time.Duration) {
	if len(s.jobEventPlanMap) == 0 {
		return time.Second
	}

	now := time.Now()
	var nearestTime *time.Time
	for _, v := range s.jobEventPlanMap {
		if v.NextTime.Before(now) || v.NextTime.Equal(now) {
			s.tryStartJob(v)
			v.NextTime = v.Expr.Next(now)
		}

		if nearestTime == nil || v.NextTime.Before(*nearestTime) {
			nearestTime = &v.NextTime
		}
	}

	return nearestTime.Sub(now)
}

func (s *Scheduler) Push(je *JobEvent) {
	s.jobEventChan <- je
}

func (s *Scheduler) tryStartJob(plan *SchedulePlan) {
	fun := "Scheduler.tryStartJob -->"

	// 任务还在执行则跳过本次调度
	if _, ok := s.jobExecMap[plan.Job.Name]; ok {
		logrus.Warnf("%s job:%s executeing", fun, plan.Job)
		return
	}

	// 加入执行表
	s.jobExecMap[plan.Job.Name] = plan

}

var (
	GScheduler *Scheduler
)

func InitScheduler() (err error) {
	GScheduler = &Scheduler{}

	go GScheduler.scheduleLoop()

	return
}

type SchedulePlan struct {
	Job        *Job
	Expr       *cronexpr.Expression
	NextTime   time.Time
	ctx        context.Context
	cancelFunc context.CancelFunc
}
