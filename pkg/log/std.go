package log

import (
	"fmt"
	"os"
)

const (
	perfix = "%s[2006/01/02 15:04:05] %s\n"
)

// StdLog simple consile log
type StdLog struct {
}

// NewStdLog default log print into stderr
func NewStdLog(opt ...interface{}) *StdLog {
	return &StdLog{}
}

// Init init
func (*StdLog) Init() error {
	return nil
}

// Debugf debugf
func (*StdLog) Debugf(s string, v ...interface{}) {
	text := fmt.Sprintf(perfix, "D", s)
	fmt.Fprintf(os.Stderr, text, v...)
}

// Infof infof
func (*StdLog) Infof(s string, v ...interface{}) {
	text := fmt.Sprintf(perfix, "I", s)
	fmt.Fprintf(os.Stderr, text, v...)
}

// Warnf warnf
func (*StdLog) Warnf(s string, v ...interface{}) {
	text := fmt.Sprintf(perfix, "W", s)
	fmt.Fprintf(os.Stderr, text, v...)
}

// Errorf errorf
func (*StdLog) Errorf(s string, v ...interface{}) {
	text := fmt.Sprintf(perfix, "E", s)
	fmt.Fprintf(os.Stderr, text, v...)
}

// Fatalf fatalf
func (*StdLog) Fatalf(s string, v ...interface{}) {
	text := fmt.Sprintf(perfix, "F", s)
	fmt.Fprintf(os.Stderr, text, v...)
}
