package main

import (
	"fmt"
	log "github.com/umisama/golog"
	"io/ioutil"
	"os"
)

var logger log.Logger

func outputUsage() {
	fmt.Printf("usage : %s [input-file]", os.Args[0])
	return
}

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
	logger, _ := log.NewLogger(os.Stdout, log.TIME_FORMAT_MILLISEC, log.LOG_FORMAT_POWERFUL, log.LogLevel_Debug )
	if len(os.Args) <= 1 {
		outputUsage()
		return
	}

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
