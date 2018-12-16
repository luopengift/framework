package framework

import (
	"os"
	"path/filepath"

	"github.com/luopengift/log"
)

func (app *App) initLog() error {
	file := log.NewFile(app.Option.LogPath)
	file.SetMaxBytes(app.Option.MaxBytes)
	file.SetMaxIndex(app.Option.MaxBackupIndex)
	if app.Option.Debug {
		log.SetLevel(log.DEBUG)
		log.SetOutput(file, os.Stderr)
	} else {
		log.SetLevel(log.INFO)
		log.SetOutput(file)
	}
	return os.MkdirAll(filepath.Dir(app.Option.LogPath), 0755)
}
