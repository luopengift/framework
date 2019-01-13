package framework

import (
	"flag"
	"os"

	"github.com/luopengift/types"
)

func newArgsOpt() *Option {
	opt := &Option{}
	flag.StringVar(&opt.ConfigPath, "conf", "conf.yml", "(conf)配置文件")
	flag.BoolVar(&opt.Debug, "debug", false, "(debug)调试模式")
	flag.StringVar(&opt.Tz, "tz", "Asia/Shanghai", "(timezone)时区")
	flag.StringVar(&opt.PprofPath, "pprof", "", "(pprof)性能分析路径")
	flag.BoolVar(&opt.Version, "version", false, "(version)版本")
	flag.StringVar(&opt.Httpd, "httpd", "", "(httpd)IP:端口")
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
