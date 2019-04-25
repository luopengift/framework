package log

import (
	"fmt"
	"io"
)

// Log log
type Log struct {
	out    io.Writer
	Labels map[string]string
}

// func (o *Log) output(s string, v ...interface{}) string {
// 	text := fmt.Sprintf(s, v...)

// }

// Debugf debugf
func (o *Log) Debugf(s string, v ...interface{}) {
	text := fmt.Sprintf(perfix, "D", s)
	fmt.Fprintf(o.out, text, v...)
}

// Infof infof
func (o *Log) Infof(s string, v ...interface{}) {
	text := fmt.Sprintf(perfix, "I", s)
	fmt.Fprintf(o.out, text, v...)
}

// Warnf warnf
func (o *Log) Warnf(s string, v ...interface{}) {
	text := fmt.Sprintf(perfix, "W", s)
	fmt.Fprintf(o.out, text, v...)
}

// Errorf errorf
func (o *Log) Errorf(s string, v ...interface{}) {
	text := fmt.Sprintf(perfix, "E", s)
	fmt.Fprintf(o.out, text, v...)
}

// Fatalf fatalf
func (o *Log) Fatalf(s string, v ...interface{}) {
	text := fmt.Sprintf(perfix, "F", s)
	fmt.Fprintf(o.out, text, v...)
}
