package daemon

type stopCmd struct {
	Force bool `name:"force" short:"f" help:"Use kill -9 to stop daeomn."`
}

// [Interface Func]
func (c stopCmd) Run() error {
	return nil
}
