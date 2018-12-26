package framework

import (
	"flag"
	"os"

	"github.com/luopengift/types"
)

func newArgsOpt() *Option {
	conf := flag.String("conf", "conf.yml", "(conf)配置文件")
	debug := flag.Bool("debug", false, "(debug)调试模式")
	pprof := flag.String("pprof", "", "(pprof)性能分析路径")
	version := flag.Bool("version", false, "(version)版本")
	httpd := flag.String("httpd", "", "(httpd)IP:端口")
	flag.Parse()
	return &Option{
		Version:    *version,
		Debug:      *debug,
		PprofPath:  *pprof,
		ConfigPath: *conf,
		Httpd:      *httpd,
	}
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
