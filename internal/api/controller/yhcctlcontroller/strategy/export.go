package strategy

type exportCmd struct {
	Output string `name:"output" help:"Export default configuration file to file. "`
}

// [Interface Func]
func (c exportCmd) Run() error {
	return nil
}
