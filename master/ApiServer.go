package master

import (
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
func handleJobSave(w http.ResponseWriter, r *http.Request) {

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
