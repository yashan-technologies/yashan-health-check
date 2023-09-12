package checkcontroller

type CheckGlobal struct {
	Range  string `name:"range"  short:"r" help:"The time range of the check, such as '1M', '1d', '1h', '1m'. If <range> is given, <start> and <end> will be discard."`
	Start  string `name:"start"  short:"s" help:"The start datetime of the check, such as 'yyyy-MM-dd', 'yyyy-MM-dd-hh', 'yyyy-MM-dd-hh-mm'"`
	End    string `name:"end"    short:"e" help:"The end timestamp of the check, such as 'yyyy-MM-dd', 'yyyy-MM-dd-hh', 'yyyy-MM-dd-hh-mm', default value is current datetime."`
	Output string `name:"output" short:"o" help:"The output dir of the check."`
}

type CheckCmd struct {
	CheckGlobal
}

// [Interface Func]
func (c *CheckCmd) Run() error {
	c.fillDefault()
	if err := c.validate(); err != nil {
		return err
	}
	return nil
}
