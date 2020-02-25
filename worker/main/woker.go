package main

import (
	"crontab-scheduler-go/worker"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var (
	confFile string
)

//parsing command line args
func initArgs() {
	//worker -config ./worker.json
	flag.StringVar(&confFile, "config", "./worker.json", "worker.json")
	flag.Parse()
}

func initEnv() {
	//get the num of cpu in the work machine
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)
	//init command args
	initArgs()
	//thread initialization
	initEnv()

	//load config
	if err = worker.InitConig(confFile); err != nil {
		goto ERR
	}
	if err = worker.InitScheduler(); err != nil {
		goto ERR
	}
	if err = worker.InitJobMgr(); err != nil {
		goto ERR
	}
	if err = worker.InitExecutor(); err != nil {
		goto ERR
	}
	if err = worker.InitLogSink(); err != nil {
		goto ERR
	}
	for {
		time.Sleep(1 * time.Second)
	}
	return

ERR:
	fmt.Println(err)
}
