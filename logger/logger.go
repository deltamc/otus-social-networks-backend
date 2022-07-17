package logger

import (
	"flag"
	"os"
	"log"
)

var (
	Log      *log.Logger
)


func init() {
	// set location of log file
	var logpath = os.Getenv("LOG_FILE")

	flag.Parse()
	var file, err1 = os.Create(logpath)

	if err1 != nil {
		panic(err1)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
	Log.Println("LogFile : " + logpath)
}