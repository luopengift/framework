package framework

// Prepare Preparer interface
func (app *App) Prepare(prepare Function) {
	app.onPrepare = prepare
}

// PrepareFunc ParpareFunc
func (app *App) PrepareFunc(f Func) {
	app.onPrepare = f
}

// Init Initer interface
func (app *App) Init(init Function) {
	app.onInit = init
}

// InitFunc init func program global var
func (app *App) InitFunc(f Func) {
	app.onInit = f
}

// Main interface
func (app *App) Main(main Function) {
	app.onMain = main
}

// MainFunc handle main func
func (app *App) MainFunc(f Func) {
	app.onMain = f
}

// Exit Exiter interface
func (app *App) Exit(exit Function) {
	app.onExit = exit
}

// ExitFunc init func program global var
func (app *App) ExitFunc(f Func) {
	app.onExit = f
}
