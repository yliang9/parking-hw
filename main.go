package main

func main() {
    logFile := setLog("./parking.log")
    defer logFile.Close()
    a := App{}
    a.Init()
    a.Run(":8080")
}
