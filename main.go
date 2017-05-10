//project of a parking lot billing system
package main

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
