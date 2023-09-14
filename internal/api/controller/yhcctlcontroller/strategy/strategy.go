package strategy

type StrategyCmd struct {
	Show    showCmd    `cmd:"show"    name:"show"    help:"Show current strategy."`
	Update  updateCmd  `cmd:"update"  name:"update"  help:"Update current strategy."`
	Replace replaceCmd `cmd:"replace" name:"replace" help:"Replace current strategy."`
	Export  exportCmd  `cmd:"export"  name:"export"  help:"Export current strategy."`
}

// [Interface Func]
func (c StrategyCmd) Run() error {
	return nil
}
