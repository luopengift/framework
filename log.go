package framework

import (
	"os"
	"path/filepath"

	"github.com/luopengift/log"
)

// Logger interface
type Logger interface {
	Debug(v interface{})
	Info(v interface{})
	Warn(v interface{})
	Error(v interface{})
}

// Log is framework log module
var Log Logger

// InitLog init log
func (app *App) InitLog() error {
	if err := os.MkdirAll(filepath.Dir(app.Option.LogPath), 0755); err != nil {
		return err
	}
	file := log.NewFile(app.Option.LogPath)
	file.SetMaxBytes(app.Option.LogMaxBytes)
	file.SetMaxIndex(app.Option.LogMaxBackupIndex)
	if app.Option.Debug {
		log.SetLevel(log.DEBUG)
		log.SetOutput(file, os.Stderr)
	} else {
		log.SetLevel(log.INFO)
		log.SetOutput(file)
	}
	log.SetTextFormat(app.LogTextFormat, app.LogMode)
	log.SetTimeFormat("2006-01-02 15:04:05.000")
	return nil
}
