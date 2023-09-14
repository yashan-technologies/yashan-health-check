package strategy

type updateCmd struct {
	DisableInteractive bool   `name:"disable-interactive" default:"false" help:"Disable interactive mode."`
	Type               string `name:"type" help:"If interactive mode is disabled, the type parameter is used to set the cycle type of the strategy"`
	Time               string `name:"time" help:"If interactive mode is disabled, the type parameter is used to set the startTime of the strategy"`
	PeriodDays         string `name:"period-days" help:"If the interactive mode is disabled and the type parameter is set to weekly or monthly, the PeriodDay parameter can set the day of execution, which can be set to 1-7 or 1-31, split by ','."`
}

// [Interface Func]
func (c updateCmd) Run() error {
	return nil
}
