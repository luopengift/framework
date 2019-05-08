package framework

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/luopengift/version"
)

// Option option
type Option struct {
	Level      string `json:"level,omitempty" yaml:"level" env:"LEVEL"`                   // 控制台日志, Level=debug, info, warn, error, fatal
	Tz         string `json:"tz,omitempty" yaml:"tz" env:"TZ"`                            // 时区, 默认Asia/Shanghai
	PprofPath  string `json:"pprof_path,omitempty" yaml:"pprof_path" env:"PPROF_PATH"`    // 性能分析路径
	ConfigPath string `json:"config_path,omitempty" yaml:"config_path" env:"CONFIG_PATH"` // 配置文件路径
	Httpd      string `json:"httpd,omitempty" yaml:"httpd" env:"HTTPD"`                   // http监听地址
}

func (opt *Option) mergeIn(o *Option) {
	if o.Level != emptyOption.Level {
		opt.Level = o.Level
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
		Tz: "Asia/Shanghai",
	}
)

func newArgsOpt() *Option {
	if !flag.Parsed() {
		v := flag.Bool("version", false, "(version)版本")
		flag.StringVar(&argsOption.ConfigPath, "conf", emptyOption.ConfigPath, "(conf)配置文件")
		flag.StringVar(&argsOption.Level, "level", emptyOption.Level, "(level)日志级别, 支持级别[debug, info, warn, error, fatal]")
		flag.StringVar(&argsOption.Tz, "tz", emptyOption.Tz, "(timezone)时区")
		flag.StringVar(&argsOption.PprofPath, "pprof", emptyOption.PprofPath, "(pprof)性能分析路径")
		flag.StringVar(&argsOption.Httpd, "httpd", emptyOption.Httpd, "(httpd)IP:端口")
		flag.Parse()
		if *v {
			fmt.Println(version.String())
			os.Exit(0)
		}
	}
	return argsOption
}

// LoadEnv load system env
func newEnvOpt() *Option {
	opt := &Option{}
	opt.Level = os.Getenv("LEVEL")
	opt.ConfigPath = os.Getenv("CONFIG_PATH")
	opt.Tz = os.Getenv("TZ")
	opt.Httpd = os.Getenv("HTTPD")
	return opt
}

// Format  dest by dest
func Format(dest, src interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dest)
}

// UpdateTo update to
func UpdateTo(src interface{}, dests ...interface{}) error {
	for _, dest := range dests {
		if err := Format(dest, src); err != nil {
			return err
		}
	}
	return nil
}

// UpdateFrom update form
func UpdateFrom(dest interface{}, srcs ...interface{}) error {
	fmt.Println("update", dest, srcs)
	for _, src := range srcs {
		if err := Format(dest, src); err != nil {
			return err
		}
	}
	return nil
}
