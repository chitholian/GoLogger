// Package logger provides an interface for logger with different levels like Warn, Error, Info etc.
// It also provides a basic implementation usable as is out of the box.
// The basic implementation is similar to GO standard log package.

package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Level is an alias to int. It indicates the log level.
type Level int

const (
	// LevelQuiet is to be used in ILogger.SetLevel function. It indicates no log should be written to output.
	LevelQuiet Level = iota
	// LevelFatal is similar to LevelError except it tells the logger to call os.Exit too.
	LevelFatal
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

const (
	// FlagColorMode indicates logs should be colorized based on their levels, e.g. red for LevelError.
	FlagColorMode = 1 << iota
)

// These prefix characters are to be prepended to every log entries.
var levelPrefixes = map[Level]string{
	LevelFatal: "F",
	LevelError: "E",
	LevelWarn:  "W",
	LevelInfo:  "I",
	LevelDebug: "D",
	LevelTrace: "T",
}

// These colors will be used to colorize logs when FlagColorMode is set.
var levelColors = map[Level][]byte{
	LevelFatal: []byte("\033[31;1m"),
	LevelError: []byte("\033[31m"),
	LevelWarn:  []byte("\033[33;1m"),
	LevelInfo:  []byte("\033[32m"),
	LevelDebug: []byte("\033[34;1m"),
	LevelTrace: []byte("\033[36m"),
}

// ILogger is an interface for simple and easy logging system.
type ILogger interface {
	// SetLevel sets the maximum Level to current instance.
	// Entries greater than this Level should not be printed to output.
	SetLevel(level Level)
	// GetLevel returns current Level of the logger instance.
	GetLevel() Level

	// SetFlags sets flags to the logger instance.
	SetFlags(flags int)
	// GetFlags returns current flags of the logger instance.
	GetFlags() int

	// SetPrefix sets a prefix to be used with every log entries.
	SetPrefix(prefix string)
	// GetPrefix returns the prefix currently set.
	GetPrefix() string

	// SetOutput sets an io.Writer as target where logs should be printed.
	// For example os.Stderr can be used to log to console.
	SetOutput(out io.Writer)
	// GetOutput returns an io.Writer where logs are to be written currently.
	GetOutput() io.Writer

	// Print writes a log entry to the output. Behaves like fmt.Print standard function.
	// It should return immediately (writing nothing) if current log level is smaller than the passed Level.
	// But if the passed Level is LevelFatal, then os.Exit should be called before return.
	Print(level Level, v ...any)

	// Println writes a log entry to the output. Behaves like fmt.Println standard function.
	// It should return immediately (writing nothing) if current log level is smaller than the passed Level.
	// But if the passed Level is LevelFatal, then os.Exit should be called before return.
	Println(level Level, v ...any)

	// Printf writes a log entry to the output. Behaves like fmt.Printf standard function.
	// It should return immediately (writing nothing) if current log level is smaller than the passed Level.
	// But if the passed Level is LevelFatal, then os.Exit should be called before return.
	Printf(level Level, format string, v ...any)

	// Clone returns an identical copy of the current log instance.
	// It is useful when you need to create multiple loggers with similar configuration.
	Clone() ILogger
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func iToA(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// logger is a simple implementation of ILogger to be used out of the box.
type logger struct {
	level  int32
	prefix string
	flags  int
	out    io.Writer
	buf    []byte
	sync.Mutex
}

// SetLevel sets max log level. Passing LevelQuiet will make the logger write nothing to the output.
func (l *logger) SetLevel(level Level) {
	if level < LevelQuiet {
		level = LevelQuiet
	} else if level > LevelTrace {
		level = LevelTrace
	}
	atomic.StoreInt32(&l.level, int32(level))
}

func (l *logger) GetLevel() Level {
	return Level(atomic.LoadInt32(&l.level))
}

func (l *logger) SetFlags(flags int) {
	l.Lock()
	defer l.Unlock()
	l.flags = flags
}

func (l *logger) GetFlags() int {
	l.Lock()
	defer l.Unlock()
	return l.flags
}

func (l *logger) SetPrefix(prefix string) {
	l.Lock()
	defer l.Unlock()
	l.prefix = prefix
}

func (l *logger) GetPrefix() string {
	l.Lock()
	defer l.Unlock()
	return l.prefix
}

func (l *logger) SetOutput(out io.Writer) {
	l.Lock()
	defer l.Unlock()
	l.out = out
}

func (l *logger) GetOutput() io.Writer {
	l.Lock()
	defer l.Unlock()
	return l.out
}

func (l *logger) buildHeader(level Level, buf *[]byte, t time.Time) {
	pref, _ := levelPrefixes[level]
	*buf = append(*buf, pref...)
	*buf = append(*buf, '/')
	hour, min, sec := t.Clock()
	iToA(buf, hour, 2)
	*buf = append(*buf, ':')
	iToA(buf, min, 2)
	*buf = append(*buf, ':')
	iToA(buf, sec, 2)
	*buf = append(*buf, ' ')
	*buf = append(*buf, l.prefix...)
	*buf = append(*buf, ": "...)
}

func (l *logger) printOut(level Level, s string) error {
	now := time.Now()
	l.Lock()
	defer l.Unlock()
	l.buf = l.buf[:0]
	color, hasColor := levelColors[level]
	if hasColor = hasColor && l.flags&FlagColorMode != 0; hasColor {
		l.buf = append(l.buf, color...)
	}
	l.buildHeader(level, &l.buf, now)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	if hasColor {
		l.buf = append(l.buf, "\033[0m"...)
	}
	_, e := l.out.Write(l.buf)
	return e
}

func (l *logger) Print(level Level, v ...any) {
	if atomic.LoadInt32(&l.level) < int32(level) {
		if level == LevelFatal {
			os.Exit(1)
		}
		return
	}
	_ = l.printOut(level, fmt.Sprint(v...))
	if level == LevelFatal {
		os.Exit(1)
	}
}

func (l *logger) Println(level Level, v ...any) {
	if atomic.LoadInt32(&l.level) < int32(level) {
		if level == LevelFatal {
			os.Exit(1)
		}
		return
	}
	_ = l.printOut(level, fmt.Sprintln(v...))
	if level == LevelFatal {
		os.Exit(1)
	}
}

func (l *logger) Printf(level Level, format string, v ...any) {
	if atomic.LoadInt32(&l.level) < int32(level) {
		if level <= LevelFatal {
			os.Exit(1)
		}
		return
	}
	_ = l.printOut(level, fmt.Sprintf(format, v...))
	if level <= LevelFatal {
		os.Exit(1)
	}
}

func (l *logger) Clone() ILogger {
	l.Lock()
	defer l.Unlock()
	newLog := New(Level(atomic.LoadInt32(&l.level)), l.prefix, l.out, l.flags)
	return newLog
}

func New(level Level, prefix string, out io.Writer, flags int) ILogger {
	l := logger{
		prefix: prefix,
		level:  int32(level),
		flags:  flags,
		out:    out,
	}
	return &l
}

// std is the default instance created to be used out of the box.
var std = New(LevelWarn, "", os.Stderr, 0)

// GetDefault returns a simple implementation of ILogger.
// It is used when you call logger.Print etc. functions without creating an instance.
func GetDefault() ILogger {
	return std
}

// NewDefault returns a clone of the default instance.
func NewDefault() ILogger {
	return std.Clone()
}

// Print writes a log entry to the output using default instance. Behaves like fmt.Print standard function.
// It returns immediately (writing nothing) if current log level is smaller than the passed Level.
// But if the passed Level is LevelFatal, then os.Exit will be called before return.
func Print(level Level, v ...any) {
	std.Print(level, v...)
}

// Printf writes a log entry to the output using default instance. Behaves like fmt.Printf standard function.
// It returns immediately (writing nothing) if current log level is smaller than the passed Level.
// But if the passed Level is LevelFatal, then os.Exit will be called before return.
func Printf(level Level, format string, v ...any) {
	std.Printf(level, format, v...)
}

// Println writes a log entry to the output using default instance. Behaves like fmt.Println standard function.
// It returns immediately (writing nothing) if current log level is smaller than the passed Level.
// But if the passed Level is LevelFatal, then os.Exit will be called before return.
func Println(level Level, format string, v ...any) {
	std.Printf(level, format, v...)
}
