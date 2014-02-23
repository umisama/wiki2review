package main

import (
	"fmt"
	"github.com/stretchr/goweb"
	log "github.com/umisama/golog"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var logger log.Logger

const LISTEN_ADDR = ":8080"

func GetStringsFromFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf(err.Error())
		return ""
	}

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf(err.Error())
		return ""
	}

	return string(buf)
}

func main() {
	logger, _ := log.NewLogger(os.Stdout, log.TIME_FORMAT_MILLISEC, log.LOG_FORMAT_POWERFUL, log.LogLevel_Debug)

	// go to localmode if arg is none.
	if len(os.Args) <= 1 {
		logger.Info("start servermode")

		err := HttpListenAndServe()
		logger.Info("server was end with err = ", err)
		return
	} else {
		src := GetStringsFromFile(os.Args[1])
		cvtr := NewConverter(src)
		if cvtr == nil {
			return
		}

		err := cvtr.DoConvert()
		if err != nil {
			logger.Debug(err)
			return
		}
		fmt.Print(cvtr.GetResult())
	}
}

func HttpListenAndServe() (err error) {
	s := http.Server{
		Addr:           LISTEN_ADDR,
		Handler:        goweb.DefaultHttpHandler(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	goweb.Map("/cnvt", ConvertHandler)
	err = s.ListenAndServe()
	return
}
