//project of a parking lot billing system
//     Schemes: http, https
//     Host: localhost
//     Version: 0.0.1
//
// swagger:meta

package main

import (
	"log"
	"os"
)

var Log *log.Logger

//main function
//after started, it listening on the port 8080 for incoming request
func main() {
	port := ":8080"
	logFile := setLog("./parking.log")
	defer logFile.Close()
	a := App{}
	a.Init()
	a.Run(port)
}

//set log file
func setLog(logfile string) (file *os.File) {
	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
	return file
}
