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

//Job saving handler
//POST job={"name":"job1","command":"echo xxx","cron_expr":"* * * * *"}
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
		fmt.Println(job)
		goto ERR
	}
	//save job to etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	//normal resp ({"errno":0,"msg":"",data:{...}})
	if resp, err = common.BuildResponse(0, "success", oldJob); err == nil {
		w.Write(resp)
	}
	return
ERR:
	if resp, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		w.Write(resp)
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
