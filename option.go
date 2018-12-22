package framework

import (
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
	PprofPath      string `json:"pprof_path" yaml:"pprof_path" env:"PPROF_PATH"`                   // 性能分析路径
	ConfigPath     string `json:"config_path" yaml:"config_path" env:"CONFIG_PATH"`                // 配置文件路径
	LogPath        string `json:"log_path" yaml:"log_path" env:"LOG_PATH"`                         // 日志文件路径
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
	if o.PprofPath != emptyOption.PprofPath {
		opt.PprofPath = o.PprofPath
	}
	if o.ConfigPath != emptyOption.ConfigPath {
		opt.ConfigPath = o.ConfigPath
	}
	if o.LogPath != emptyOption.LogPath {
		opt.LogPath = o.LogPath
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
	defaultOption = &Option{
		LogPath:        "logs/%Y-%M-%D.log",
		MaxBytes:       200 * 1024 * 1024, //200M
		MaxBackupIndex: 50,
	}
)
