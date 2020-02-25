package worker

import (
	"crontab-scheduler-go/common"
	"math/rand"
	"os/exec"
	"time"
)

type Executor struct {
}

var (
	G_executor *Executor
)

func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		var (
			cmd     *exec.Cmd
			output  []byte
			result  *common.JobExecuteResult
			jobLock *JobLock
			err     error
		)
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}
		//get distributed locks
		jobLock = G_jobMgr.CreateJobLock(info.Job.Name)
		result.StartTime = time.Now()

		//random sleep
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		if err = jobLock.TryLock(); err != nil {
			result.Err = err
			result.EndTime = time.Now()
			return
		}
		result.StartTime = time.Now()
		cmd = exec.CommandContext(info.CancelCtx, "/bin/bash", "-c", info.Job.Command)
		output, err = cmd.CombinedOutput()
		result.EndTime = time.Now()
		result.Output = output
		result.Err = err
		//After job execution is completed,return the result to the scheduler
		G_scheduler.PushJobResult(result)
		jobLock.Unlock()
	}()
}

//init executor
func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}
