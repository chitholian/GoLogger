package main

import (
	log "github.com/chitholian/GoLogger"
	"os"
)

func main() {
	// Create a logger with debug level and colorization feature.
	l := log.New(log.LevelDebug, "COLORED", os.Stderr, log.FlagColorMode)
	l.Println(log.LevelWarn, "This warning entry is from instanced logger")

	// Also default logger can be used simultaneously.
	log.Println(log.LevelError, "This error entry is from default logger")
}
