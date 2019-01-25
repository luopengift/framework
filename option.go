package framework

import (
	"flag"
	"os"
	"reflect"

	"github.com/luopengift/types"
)

// Option option 配置优先级(级别递增, 简单的理解就是后面的配置,会覆盖前面的配置)
// 1. emptyOption
// 2. defaultOption
// 3. New函数设置的OptionList
// 4. 用户Config中同名的配置项
// 5. 系统环境变量
// 6. 命令行配置
// 7. API或者其他方式修改的配置项(TODO)
type Option struct {
	Version        bool   `json:"version" yaml:"version" env:"VERSION"`                            // 程序信息
	Debug          bool   `json:"debug" yaml:"debug" env:"DEBUG"`                                  // 控制台日志, Level=Debug
	Tz             string `json:"tz" yaml:"tz" env:"TZ"`                                           // 时区, 默认Asia/Shanghai
	PprofPath      string `json:"pprof_path" yaml:"pprof_path" env:"PPROF_PATH"`                   // 性能分析路径
	ConfigPath     string `json:"config_path" yaml:"config_path" env:"CONFIG_PATH"`                // 配置文件路径
	LogPath        string `json:"log_path" yaml:"log_path" env:"LOG_PATH"`                         // 日志文件路径
	LogTextFormat  string `json:"log_text_format" yaml:"log_text_format" env:"LOG_TEXT_FORMAT"`    //日志消息格式
	LogMode        int    `json:"log_mode" yaml:"log_mode" env:"LOG_MODE"`                         //日志级别颜色, 1: 无颜色, 0:有颜色
	MaxBytes       int    `json:"max_bytes" yaml:"max_bytes" env:"MAX_BYTES"`                      // 日志文件大小
	MaxBackupIndex int    `json:"max_backup_index" yaml:"max_backup_index" env:"MAX_BACKUP_INDEX"` // 日志文件数量
	ReportURL      string `json:"report_url" yaml:"report_url" env:"REPORT_URL"`                   // 数据上报地址
	Httpd          string `json:"httpd" yaml:"httpd" env:"HTTPD"`                                  // http监听地址
}

// LoadEnv load system env
func (opt *Option) _LoadEnv() error {
	rv := reflect.Indirect(reflect.ValueOf(opt))
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		v, ok := rt.Field(i).Tag.Lookup("env")
		if !ok {
			continue
		}
		env := os.Getenv(v)
		if env != "" {
			field := rv.Field(i)
			switch field.Kind() {
			case reflect.String:
				field.SetString(env)
			case reflect.Bool:
				value, err := types.StringToBool(env)
				if err != nil {
					return err
				}
				field.SetBool(value)
			case reflect.Int:
				value, err := types.StringToInt64(env)
				if err != nil {
					return err
				}
				field.SetInt(value)
			case reflect.Slice:
				// value := strings.Split(env, ",")
				// field.SetSlice(value)
				// //fmt.Println(field.Kind())
			}
		}
		//fmt.Println(rv.Field(i).Kind())
	}
	return nil
}

func (opt *Option) mergeIn(o *Option) {
	if o.Version != emptyOption.Version {
		opt.Version = o.Version
	}
	if o.Debug != emptyOption.Debug {
		opt.Debug = o.Debug
	}
	if o.Tz != emptyOption.Tz {
		opt.Tz = o.Tz
	}
	if o.PprofPath != emptyOption.PprofPath {
		opt.PprofPath = o.PprofPath
	}
	if o.ConfigPath != emptyOption.ConfigPath {
		opt.ConfigPath = o.ConfigPath
	}
	if o.LogPath != emptyOption.LogPath {
		opt.LogPath = o.LogPath
	}
	if o.LogTextFormat != emptyOption.LogTextFormat {
		opt.LogTextFormat = o.LogTextFormat
	}
	if o.LogMode != emptyOption.LogMode {
		opt.LogMode = o.LogMode
	}
	if o.MaxBytes != emptyOption.MaxBytes {
		opt.MaxBytes = o.MaxBytes
	}
	if o.MaxBackupIndex != emptyOption.MaxBackupIndex {
		opt.MaxBackupIndex = o.MaxBackupIndex
	}
	if o.ReportURL != emptyOption.ReportURL {
		opt.ReportURL = o.ReportURL
	}
	if o.Httpd != emptyOption.Httpd {
		opt.Httpd = o.Httpd
	}
}

// Merge merge
func (opt *Option) Merge(opts ...*Option) {
	for _, o := range opts {
		opt.mergeIn(o)
	}
}

var (
	emptyOption   = &Option{}
	argsOption    = &Option{}
	defaultOption = &Option{
		Tz:             "Asia/Shanghai",
		LogPath:        "logs/%Y-%M-%D.log",
		LogTextFormat:  "TIME [LEVEL] FILE:LINE MESSAGE",
		LogMode:        2,
		MaxBytes:       200 * 1024 * 1024, //200M
		MaxBackupIndex: 50,
	}
)

func newArgsOpt() *Option {
	if !flag.Parsed() {
		flag.StringVar(&argsOption.ConfigPath, "conf", emptyOption.ConfigPath, "(conf)配置文件")
		flag.BoolVar(&argsOption.Debug, "debug", emptyOption.Debug, "(debug)调试模式")
		flag.StringVar(&argsOption.LogPath, "log", emptyOption.LogPath, "(log)日志文件")
		flag.StringVar(&argsOption.Tz, "tz", emptyOption.Tz, "(timezone)时区")
		flag.StringVar(&argsOption.PprofPath, "pprof", emptyOption.PprofPath, "(pprof)性能分析路径")
		flag.BoolVar(&argsOption.Version, "version", emptyOption.Version, "(version)版本")
		flag.StringVar(&argsOption.Httpd, "httpd", emptyOption.Httpd, "(httpd)IP:端口")
		flag.Parse()
	}
	return argsOption
}

// LoadEnv load system env
func newEnvOpt() (*Option, error) {
	var err error
	opt := &Option{}
	debug := os.Getenv("DEBUG")
	if debug != "" {
		if opt.Debug, err = types.StringToBool(debug); err != nil {
			return opt, err
		}
	}
	opt.ConfigPath = os.Getenv("CONFIG_PATH")
	opt.Tz = os.Getenv("TZ")
	opt.LogPath = os.Getenv("LOG_PATH")
	opt.LogTextFormat = os.Getenv("LOG_TEXT_FORMAT")
	logMode := os.Getenv("LOG_MODE")
	if logMode != "" {
		if opt.LogMode, err = types.StringToInt(logMode); err != nil {
			return opt, err
		}
	}

	maxBytes := os.Getenv("MAX_BYTES")
	if maxBytes != "" {
		if opt.MaxBytes, err = types.StringToInt(maxBytes); err != nil {
			return opt, err
		}
	}
	maxBackupIndex := os.Getenv("MAX_BACKUP_INDEX")
	if maxBackupIndex != "" {
		if opt.MaxBackupIndex, err = types.StringToInt(maxBytes); err != nil {
			return opt, err
		}
	}
	opt.ReportURL = os.Getenv("REPORT_URL")
	opt.Httpd = os.Getenv("HTTPD")
	return opt, err
}
