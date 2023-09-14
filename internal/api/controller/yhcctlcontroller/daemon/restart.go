package daemon

type restartCmd struct {
	Force bool `name:"force" short:"f" help:"Use kill -9 to stop daeomn, then restart daemon."`
}

// [Interface Func]
func (c restartCmd) Run() error {
	return nil
}
