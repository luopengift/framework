package framework

import (
	"os"
	"path/filepath"

	"github.com/luopengift/log"
)

// Logger log
type Logger struct {
	Debug          bool   `json:"debug,omitempty" yaml:"debug,omitempty" env:"DEBUG"`                                  // 控制台日志, Level=Debug
	Path           string `json:"path,omitempty" yaml:"path,omitempty" env:"PATH"`                                     // 日志文件路径
	TextFormat     string `json:"text_format,omitempty" yaml:"text_format,omitempty" env:"TEXT_FORMAT"`                // 日志消息格式
	TimeFormat     string `json:"time_format,omitempty" yaml:"time_format,omitempty" env:"TIME_FORMAT"`                // 时间格式
	MaxBytes       int    `json:"max_bytes,omitempty" yaml:"max_bytes,omitempty" env:"MAX_BYTES"`                      // 日志文件大小
	MaxBackupIndex int    `json:"max_backup_index,omitempty" yaml:"max_backup_index,omitempty" env:"MAX_BACKUP_INDEX"` // 日志文件数量
	Depth          int    `json:"depth,omitempty" yaml:"deph,omitempty" env:"DEPTH"`                                   // 日志深度
	*log.Log       `json:"-" yaml:"-"`
}

// func (opt *Logger) mergeIn(o *Logger) {
// 	if o == nil {
// 		return
// 	}
// 	opt.Debug = o.Debug
// 	if o.Path != "" {
// 		opt.Path = o.Path
// 	}
// 	if o.TextFormat != "" {
// 		opt.TextFormat = o.TextFormat
// 	}
// 	if o.TimeFormat != "" {
// 		opt.TimeFormat = o.TimeFormat
// 	}
// 	opt.MaxBytes = o.MaxBytes
// 	opt.MaxBackupIndex = o.MaxBackupIndex
// 	if o.Depth != 0 {
// 		opt.Depth = o.Depth
// 	}
// }

// NewLog new Log opt
func NewLog() *Logger {
	return &Logger{
		Path:       "logs/%Y-%M-%D.log",
		TextFormat: "TIME [LEVEL] FILE:LINE MESSAGE",
		TimeFormat: "2006-01-02 15:04:05.000",
		Depth:      3,
	}
}

// Init init log
func (logger *Logger) Init() error {
	if err := os.MkdirAll(filepath.Dir(logger.Path), 0755); err != nil {
		return err
	}

	w := log.NewFile(logger.Path)
	w.SetMaxBytes(logger.MaxBytes)
	w.SetMaxIndex(logger.MaxBackupIndex)
	logger.Log = log.NewLog("framework", w)

	if logger.Debug {
		logger.Log.SetLevel(log.DEBUG)
		logger.Log.SetOutput(w, os.Stderr)
	} else {
		logger.Log.SetLevel(log.INFO)
	}
	logger.Log.SetTextFormat(logger.TextFormat, 0) //1: 无颜色, 0:有颜色
	logger.Log.SetTimeFormat(logger.TimeFormat)
	logger.Log.SetCallDepth(logger.Depth)
	return nil
}
