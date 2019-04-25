package framework

import (
	"os"
	"path/filepath"

	"github.com/luopengift/log"
)

// Logger log
type Logger struct {
	Debug          bool   `json:"debug" yaml:"debug" env:"DEBUG"`                                  // 控制台日志, Level=Debug
	Path           string `json:"path" yaml:"path" env:"PATH"`                                     // 日志文件路径
	TextFormat     string `json:"text_format" yaml:"text_format" env:"TEXT_FORMAT"`                // 日志消息格式
	TimeFormat     string `json:"time_format" yaml:"time_format" env:"TIME_FORMAT"`                // 时间格式
	MaxBytes       int    `json:"max_bytes" yaml:"max_bytes" env:"MAX_BYTES"`                      // 日志文件大小
	MaxBackupIndex int    `json:"max_backup_index" yaml:"max_backup_index" env:"MAX_BACKUP_INDEX"` // 日志文件数量
	Depth          int    `json:"depth" yaml:"deph" env:"DEPTH"`                                   // 日志深度
	*log.Log
}

func (opt *Logger) mergeIn(o *Logger) {
	if o == nil {
		return
	}
	opt.Debug = o.Debug
	if o.Path != "" {
		opt.Path = o.Path
	}
	if o.TextFormat != "" {
		opt.TextFormat = o.TextFormat
	}
	if o.TimeFormat != "" {
		opt.TimeFormat = o.TimeFormat
	}
	opt.MaxBytes = o.MaxBytes
	opt.MaxBackupIndex = o.MaxBackupIndex
	if o.Depth != 0 {
		opt.Depth = o.Depth
	}
}

// NewLogOpt new Log opt
func NewLogOpt() *Logger {
	return &Logger{
		Path:       "logs/%Y-%M-%D.log",
		TextFormat: "TIME [LEVEL] FILE:LINE MESSAGE",
		TimeFormat: "2006-01-02 15:04:05.000",
		Depth:      3,
	}
}

// Init init log
func (opt *Logger) Init() error {
	if err := os.MkdirAll(filepath.Dir(opt.Path), 0755); err != nil {
		return err
	}

	w := log.NewFile(opt.Path)
	w.SetMaxBytes(opt.MaxBytes)
	w.SetMaxIndex(opt.MaxBackupIndex)
	opt.Log = log.NewLog("framework", w)

	if opt.Debug {
		opt.Log.SetLevel(log.DEBUG)
		opt.Log.SetOutput(w, os.Stderr)
	} else {
		opt.Log.SetLevel(log.INFO)
	}
	opt.Log.SetTextFormat(opt.TextFormat, 0) //1: 无颜色, 0:有颜色
	opt.Log.SetTimeFormat(opt.TimeFormat)
	opt.Log.SetCallDepth(opt.Depth)
	return nil
}
