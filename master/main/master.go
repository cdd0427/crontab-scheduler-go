package main

import (
	"crontab-scheduler-go/master"
	"flag"
	"fmt"
	"runtime"
)

var (
	confFile string
)

//parsing command line args
func initArgs() {
	//master -config ./master.json
	flag.StringVar(&confFile, "config", "./master.json", "config file path")
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
	if err = master.InitConig(confFile); err != nil {
		goto ERR
	}
	//start up api http service
	if err = master.InitAPiServer(); err != nil {
		goto ERR
	}

	return

ERR:
	fmt.Println(err)
}
