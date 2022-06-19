package logger

import (
	"bytes"
	"fmt"
	"testing"
)

var testLevels = map[Level]string{
	// Skip LevelFatal because it will call os.Exit(1)
	LevelError: "E",
	LevelWarn:  "W",
	LevelInfo:  "I",
	LevelDebug: "D",
	LevelTrace: "T",
}

func TestAll(t *testing.T) {
	for k, v := range testLevels {
		l := GetDefault()
		l.SetLevel(LevelTrace) // Maximum output level.

		// Test for log.Println()
		buf := new(bytes.Buffer)
		l.SetOutput(buf)
		l.Println(k, "Using Println level", k, v)
		out := buf.String()
		out = out[13 : len(out)-1]
		expected := fmt.Sprintf("Using Println level %d %s", k, v)
		t.Log("Got: ", out)
		if out != expected {
			t.Errorf("Pattern mismatch,\n\texpected: %s\n\tgot: %s", expected, out)
		}

		// Test for log.Printf()
		buf = new(bytes.Buffer)
		l.SetOutput(buf)
		l.Printf(k, "Using %s level %d %s", "Printf", k, v)
		out = buf.String()
		out = out[13 : len(out)-1]
		expected = fmt.Sprintf("Using Printf level %d %s", k, v)
		t.Log("Got: ", out)
		if out != expected {
			t.Errorf("Pattern mismatch,\n\texpected: %s\n\tgot: %s", expected, out)
		}

		// Test for log.Print()
		buf = new(bytes.Buffer)
		l.SetOutput(buf)
		l.Print(k, "Using Print level name = ", k, v)
		out = buf.String()
		out = out[13 : len(out)-1]
		expected = fmt.Sprintf("Using Print level name = %d%s", k, v)
		t.Log("Got: ", out)
		if out != expected {
			t.Errorf("Pattern mismatch,\n\texpected: %s\n\tgot: %s", expected, out)
		}
	}
}
