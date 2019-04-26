package framework

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/luopengift/framework/pkg/path/pathutil"
	gcfg "gopkg.in/gcfg.v1"
	yaml "gopkg.in/yaml.v2"
)

// ConfigProvider config provider interface
type ConfigProvider interface {
}

// LoadConfig loading config step by step
// 配置优先级(级别递增, 简单的理解就是后面的配置,会覆盖前面的配置)
// 1. emptyOption
// 2. defaultOption
// 3. New函数设置的OptionList
// 4. 用户Config中同名的配置项
// 5. 系统环境变量
// 6. 命令行参数配置
// 7. API或者其他方式修改的配置项(TODO)
// 配置分为两类：框架配置(option)和用户自定义(config)配置
func (app *App) LoadConfig() error {
	envOpt := newEnvOpt()
	argsOpt := newArgsOpt()
	app.Option.Merge(envOpt, argsOpt) // 仅为了合并提取ConfigPath字段

	ok, err := pathutil.Exist(app.Option.ConfigPath)
	if err != nil {
		return err
	}
	if ok {
		for _, v := range []interface{}{app.ConfigProvider, app.Option, app} {
			if err = parseConfigFile(v, app.Option.ConfigPath); err != nil {
				return err
			}
		}
		app.Option.Merge(envOpt, argsOpt) // 修改被配置文件改掉配置
	}
	app.TimeZone, err = time.LoadLocation(app.Option.Tz)
	return err
}

// ParseConfigFile parse config file
func parseConfigFile(v interface{}, file string) error {
	filepath := strings.Replace(file, "~", os.Getenv("HOME"), -1)
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	switch suffix := path.Ext(filepath); suffix {
	case ".json", ".js":
		return json.Unmarshal(b, v)
	case ".xml":
		return xml.Unmarshal(b, v)
	case ".ini":
		return gcfg.ReadStringInto(v, string(b))
	case ".yaml", ".yml":
		return yaml.Unmarshal(b, v)
	case ".conf":
		fallthrough
	default:
		return fmt.Errorf("unknown suffix: %s", suffix)
	}
}
