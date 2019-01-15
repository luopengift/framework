package framework

import (
	"flag"
	"os"

	"github.com/luopengift/types"
)

func newArgsOpt() *Option {
	opt := &Option{}
	flag.StringVar(&opt.ConfigPath, "conf", defaultOption.ConfigPath, "(conf)配置文件")
	flag.BoolVar(&opt.Debug, "debug", defaultOption.Debug, "(debug)调试模式")
	flag.StringVar(&opt.LogPath, "log", defaultOption.LogPath, "(log)日志文件")
	flag.StringVar(&opt.Tz, "tz", defaultOption.Tz, "(timezone)时区")
	flag.StringVar(&opt.PprofPath, "pprof", defaultOption.PprofPath, "(pprof)性能分析路径")
	flag.BoolVar(&opt.Version, "version", defaultOption.Version, "(version)版本")
	flag.StringVar(&opt.Httpd, "httpd", defaultOption.Httpd, "(httpd)IP:端口")
	flag.Parse()
	return opt
}

// LoadEnv load system env
func newEnvOpt() (*Option, error) {
	var err error
	opt := &Option{}
	debug := os.Getenv("DEBUG")
	if debug != "" {
		if opt.Debug, err = types.StringToBool(debug); err != nil {
			return nil, err
		}
	}
	opt.Tz = os.Getenv("TZ")
	opt.LogPath = os.Getenv("LOG_PATH")
	maxBytes := os.Getenv("MAX_BYTES")
	if maxBytes != "" {
		if opt.MaxBytes, err = types.StringToInt(maxBytes); err != nil {
			return nil, err
		}
	}
	maxBackupIndex := os.Getenv("MAX_BACKUP_INDEX")
	if maxBackupIndex != "" {
		if opt.MaxBackupIndex, err = types.StringToInt(maxBytes); err != nil {
			return nil, err
		}
	}
	opt.ReportURL = os.Getenv("REPORT_URL")
	opt.Httpd = os.Getenv("HTTPD")
	return opt, err
}
