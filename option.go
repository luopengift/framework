package framework

// Option option 配置优先级(级别递增, 简单的理解就是后面的配置,会覆盖前面的配置)
// 1. emptyOption
// 2. defaultOption
// 3. New函数设置的&Option
// 4. 用户Config中同名的配置项
// 5. 命令行配置
// 6. API或者其他方式修改的配置项(TODO)
type Option struct {
	Debug          bool   `json:"debug" yaml:"debug"`                       // 控制台日志, Level=Debug
	LogPath        string `json:"log_path" yaml:"log_path"`                 // 日志文件路径
	MaxBytes       int    `json:"max_bytes" yaml:"max_bytes"`               // 日志文件大小
	MaxBackupIndex int    `json:"max_backup_index" yaml:"max_backup_index"` // 日志文件数量
	ReportURL      string `json:"report_url" yaml:"report_url"`             // 数据上报地址
}

func (opt *Option) mrgreIn(o *Option) {
	if o.Debug != emptyOption.Debug {
		opt.Debug = o.Debug
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
}

// Merge merge
func (opt *Option) Merge(opts ...*Option) {
	for _, o := range opts {
		opt.mrgreIn(o)
	}
}

var (
	emptyOption   = &Option{}
	defaultOption = &Option{
		Debug:          false,
		LogPath:        "logs/%Y-%M-%D.log",
		MaxBytes:       200 * 1024 * 1024, //200M
		MaxBackupIndex: 50,
	}
)
