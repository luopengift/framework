package framework

// Option option
type Option struct {
	Debug          bool
	LogPath        string
	MaxBytes       int // 日志文件大小
	MaxBackupIndex int // 日志文件数量
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
		Debug:          true,
		LogPath:        "logs/%Y-%M-%D.log",
		MaxBytes:       200 * 1024 * 1024, //200M
		MaxBackupIndex: 50,
	}
)
