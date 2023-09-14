package strategy

type replaceCmd struct {
	Config string `name:"config" help:"Configuration file path to use instead of the default configuration file. "`
}

// [Interface Func]
func (c replaceCmd) Run() error {
	return nil
}
