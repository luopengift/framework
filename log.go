package framework

// Log appFramework logger manager.
type Log struct {
	LogProvider
}

// LogProvider interface
type LogProvider interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

// SetLogger 设置日志模块
func (app *App) SetLogger(provider LogProvider) {
	app.Log = &Log{
		LogProvider: provider,
	}
}

// SetLogProvider 设置日志接口
func (app *App) SetLogProvider(provider LogProvider, setLogger ...bool) {
	if len(setLogger) == 1 && setLogger[0] {
		app.SetLogger(provider)
	}
	app.Log.LogProvider = provider
}

// GetLogProvider 获取日志接口
func (app *App) GetLogProvider() LogProvider {
	return app.Log.LogProvider
}
