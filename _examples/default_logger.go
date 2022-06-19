package main

import log "github.com/chitholian/GoLogger"

func main() {
	// Increase default log level.
	l := log.GetDefault()
	l.SetLevel(log.LevelTrace)
	log.Println(log.LevelInfo, "This is", "an", "info level log")
	log.Printf(log.LevelInfo, "Log level: %d (%s) using formatter string", log.LevelInfo, "LevelInfo")
}
