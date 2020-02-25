package common

import (
	"context"
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cron_expr"`
}

type JobSchedulePlan struct {
	Job      *Job                 //information on scheduled tasks
	Expr     *cronexpr.Expression //parsed cron expression
	NextTime time.Time            //next scheduling time
}

//job execute status
type JobExecuteInfo struct {
	Job        *Job
	PlanTime   time.Time //theoretical schedule time
	RealTime   time.Time //real schedule time
	CancelCtx  context.Context
	CancelFunc context.CancelFunc
}

//http response
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

type JobEvent struct {
	EventType int //SAVE,DELETE
	Job       *Job
}

type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo
	Output      []byte
	Err         error
	StartTime   time.Time
	EndTime     time.Time
}

//job execute log
type JobLog struct {
	JobName      string `bson:"jobName"`
	Command      string `bson:"command"`
	Err          string `bson:"err"`
	Output       string `bson:"output"`
	PlanTime     int64  `bson:"planTime"`
	ScheduleTime int64  `bson:"scheduleTime"`
	StartTime    int64  `bson:"startTime"`
	EndTime      int64  `bson:"endTime"`
}

type LogBatch struct {
	Logs []interface{}
}

//response methond
func BuildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {
	var (
		response Response
	)
	response = Response{
		Errno: errno,
		Msg:   msg,
		Data:  data,
	}
	resp, err = json.Marshal(response)
	return
}

//unserialize job
func UnpackJob(value []byte) (res *Job, err error) {
	var (
		job *Job
	)
	job = &Job{}
	if err = json.Unmarshal(value, job); err != nil {
		return
	}
	res = job
	return
}

//从etcd key中提取任务名
func ExtratJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

func ExtratKillerName(killerKey string) string {
	return strings.TrimPrefix(killerKey, JOB_KILLER_DIR)
}

func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

func BuildJobSchedulePlan(job *Job) (jobSchedulePlan *JobSchedulePlan, err error) {
	var (
		expr *cronexpr.Expression
	)
	//parse cron expression,return err if parse failed
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}

	jobSchedulePlan = &JobSchedulePlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}
	return
}

func BuildJobExecuteInfo(jobSchedulePlan *JobSchedulePlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job:      jobSchedulePlan.Job,
		PlanTime: jobSchedulePlan.NextTime,
		RealTime: time.Now(),
	}
	jobExecuteInfo.CancelCtx, jobExecuteInfo.CancelFunc = context.WithCancel(context.TODO())
	return
}
