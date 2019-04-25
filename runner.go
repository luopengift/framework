package framework

// Default default
type run struct {
	onPrepare  PrepareProvider
	onInit     InitProvider
	onMain     MainProvider
	onThreads  []ThreadProvider
	goroutines []*Goroutine
	onExit     ExitProvider
}

func (r *run) SetFunc(opts ...interface{}) {
	for _, opt := range opts {
		switch v := opt.(type) {
		case PrepareProvider:
			app.SetPrepareProvider(v)
		case PrepareFunc:
			app.SetPrepareFunc(v)
		case InitProvider:
			app.SetInitProvider(v)
		case InitFunc:
			app.SetInitFunc(v)
		case MainProvider:
			app.SetMainProvider(v)
		case MainFunc:
			app.SetMainFunc(v)
		case ExitProvider:
			app.SetExitProvider(v)
		case ExitFunc:
			app.SetExitFunc(v)
		default:
			panic("unkonw type")
		}
	}
}

func (r *run) SetPrepareProvider(provider PrepareProvider) {
	r.onPrepare = provider
}
func (r *run) SetPrepareFunc(f PrepareFunc) {
	r.onPrepare = f
}

func (r *run) SetInitProvider(provider InitProvider) {
	r.onInit = provider
}
func (r *run) SetInitFunc(f InitFunc) {
	r.onInit = f
}

func (r *run) SetThreadProvider(provider ThreadProvider) {
	r.onThreads = append(r.onThreads, provider)
}

func (r *run) SetThreadFunc(f ThreadFunc) {
	r.onThreads = append(r.onThreads, f)
}

func (r *run) SetMainProvider(provider MainProvider) {
	r.onMain = provider
}

func (r *run) SetMainFunc(f MainFunc) {
	r.onMain = f
}

func (r *run) SetExitProvider(provider ExitProvider) {
	r.onExit = provider
}
func (r *run) SetExitFunc(f ExitFunc) {
	r.onExit = f
}
