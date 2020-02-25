package worker

import (
	"crontab-scheduler-go/common"
	"fmt"
	"time"
	common2 "源码/go-crontab/source/crontab/common"
)

//job scheduler
type Scheduler struct {
	jobEventChan      chan *common.JobEvent
	jobPlanTable      map[string]*common.JobSchedulePlan
	jobExecutingTable map[string]*common.JobExecuteInfo
	jobResultChan     chan *common.JobExecuteResult
}

var (
	G_scheduler *Scheduler
)

func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		jobSchedulePlan *common.JobSchedulePlan
		jobExecuteInfo  *common.JobExecuteInfo
		jobExisted      bool
		jobExecuted     bool
		err             error
	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:
		//jobEvent.Job
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
	case common.JOB_EVENT_DELETE:
		if jobSchedulePlan, jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	case common.JOB_EVENT_KILL:
		if jobExecuteInfo, jobExecuted = scheduler.jobExecutingTable[jobEvent.Job.Name]; jobExecuted {
			jobExecuteInfo.CancelFunc()
		}
	}
}

//try execute job
func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulePlan) {
	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting   bool
	)
	//if the job is existed in job executing table,it's excuting,skip this schedule
	if _, jobExecuting = scheduler.jobExecutingTable[jobPlan.Job.Name]; jobExecuting {
		fmt.Println("job is executing,skip this schedule")
		return
	}
	jobExecuteInfo = common.BuildJobExecuteInfo(jobPlan)
	scheduler.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo
	//execute bash command
	fmt.Println("executing：", jobExecuteInfo.Job.Name, jobExecuteInfo.PlanTime, jobExecuteInfo.RealTime)
	G_executor.ExecuteJob(jobExecuteInfo)
}

//compute job scheduling status
func (scheduler *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		jobPlan  *common.JobSchedulePlan
		now      time.Time
		nearTime *time.Time
	)
	//if schedule plan table is null, sleep random time
	if len(scheduler.jobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}

	now = time.Now()

	for _, jobPlan = range scheduler.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			//todo:try to excute
			scheduler.TryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now) //update next excute time
		}
		//count the execution time of the most recent task to expire
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}
	//next scheduling interval = most recent time - now
	scheduleAfter = (*nearTime).Sub(now)
	return
}

func (scheduler *Scheduler) handleJobResult(result *common.JobExecuteResult) {
	var (
		jobLog *common.JobLog
	)
	//delete execute status
	delete(scheduler.jobExecutingTable, result.ExecuteInfo.Job.Name)
	//generate execute log
	if result.Err != common2.ERR_LOCK_ALREADY_REQUIRED {
		jobLog = &common.JobLog{
			JobName:      result.ExecuteInfo.Job.Name,
			Command:      result.ExecuteInfo.Job.Command,
			Output:       string(result.Output),
			PlanTime:     result.ExecuteInfo.PlanTime.UnixNano() / 1000 / 1000,
			ScheduleTime: result.ExecuteInfo.RealTime.UnixNano() / 1000 / 1000,
			StartTime:    result.StartTime.UnixNano() / 1000 / 1000,
			EndTime:      result.EndTime.UnixNano() / 1000 / 1000,
		}
		if result.Err != nil {
			jobLog.Err = result.Err.Error()
		} else {
			jobLog.Err = ""
		}
	}
	G_logSink.Append(jobLog)
}

//schedule coroutine
func (scheduler *Scheduler) scheduleLoop() {
	var (
		jobEvent      *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
		jobResult     *common.JobExecuteResult
	)
	scheduleAfter = scheduler.TrySchedule()

	//delay timer
	scheduleTimer = time.NewTimer(scheduleAfter)

	//watch for job event
	for {
		select {
		case jobEvent = <-scheduler.jobEventChan:
			scheduler.handleJobEvent(jobEvent)
		case <-scheduleTimer.C: //recent job expire
		case jobResult = <-scheduler.jobResultChan:
			scheduler.handleJobResult(jobResult)
		}
		scheduleAfter = scheduler.TrySchedule()
		scheduleTimer.Reset(scheduleAfter)
	}
}

//push job update event
func (scheduler Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

//scheduler init
func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		jobEventChan:      make(chan *common.JobEvent, 1000),
		jobPlanTable:      make(map[string]*common.JobSchedulePlan),
		jobExecutingTable: make(map[string]*common.JobExecuteInfo),
		jobResultChan:     make(chan *common.JobExecuteResult, 1000),
	}

	go G_scheduler.scheduleLoop()
	return
}

func (scheduler *Scheduler) PushJobResult(jobResult *common.JobExecuteResult) {
	scheduler.jobResultChan <- jobResult
}
