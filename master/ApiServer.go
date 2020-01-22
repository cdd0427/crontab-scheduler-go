package master

import (
	"crontab-scheduler-go/common"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

//http api
type ApiServer struct {
	httpServer *http.Server
}

var (
	//singleton object
	G_apiServer *ApiServer
)

//Job save handler
//POST job:{"name":"job1","command":"echo xxx","cron_expr":"* * * * *"}
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		resp    []byte
	)
	//save jobs to etcd
	//parsing post form
	if err = r.ParseForm(); err != nil {
		goto ERR
	}
	//get job from json file and unmarshall it into job struct
	postJob = r.Form.Get("job")
	fmt.Println(postJob)
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	//save job to etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	//normal resp ({"errno":0,"msg":"",data:{...}})
	if resp, err = common.BuildResponse(0, "success", oldJob); err == nil {
		_, _ = w.Write(resp)
	}
	return
ERR:
	if resp, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		_, _ = w.Write(resp)
	}
}

//Job Delete handler
//POST /job/delete name=job1
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		name   string
		oldJob *common.Job
		bytes  []byte
	)
	if err = r.ParseForm(); err != nil {
		goto ERR
	}
	name = r.Form.Get("name")
	//delete job
	if oldJob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}
	//normal resp ({"errno":0,"msg":"",data:{...}})
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		_, _ = w.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		_, _ = w.Write(bytes)
	}
}

//job list retriver
//GET
func handleJobList(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		jobList []*common.Job
		bytes   []byte
	)
	//get job list
	if jobList, err = G_jobMgr.ListJobs(); err != nil {
		goto ERR
	}
	//return normal response
	if bytes, err = common.BuildResponse(0, "success", jobList); err == nil {
		_, _ = w.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		_, _ = w.Write(bytes)
	}
}

//kill job
// POST /job/kill name=job1
func handleJobKill(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		name  string
		bytes []byte
	)
	if err = r.ParseForm(); err != nil {
		goto ERR
	}
	name = r.Form.Get("name")
	if err = G_jobMgr.KillJob(name); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
		_, _ = w.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		_, _ = w.Write(bytes)
	}
}

//serveice init
func InitAPiServer() error {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
		err        error
	)
	//router config
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)

	//start tcp listening
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return err
	}

	//http service
	httpServer = &http.Server{
		Handler:      mux,
		ReadTimeout:  time.Duration(G_config.ApiReadTimeOut) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeOut) * time.Millisecond,
	}

	//singleton init
	G_apiServer = &ApiServer{httpServer: httpServer}

	//start service
	go httpServer.Serve(listener)

	return nil
}
