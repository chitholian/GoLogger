package main

import (
	log "github.com/chitholian/GoLogger"
	"os"
)

func main() {
	// Create a logger with debug level and colorization feature.
	l := log.New(log.LevelDebug, "COLORED", os.Stderr, log.FlagColorMode)
	l.Println(log.LevelWarn, "This warning entry is from instanced logger")

	// Clone the logger
	clone := l.Clone()
	// Change the prefix of cloned logger.
	clone.SetPrefix(clone.GetPrefix() + ":CLONE")

	// Both can be used independently.
	l.Println(log.LevelError, "Error entry from COLORED")
	clone.Println(log.LevelError, "Error entry from COLORED:CLONE")
}
