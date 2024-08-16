package checkcontroller

import (
	"yhc/defs/confdef"
)

type CheckCmd struct {
	CheckGlobal
}

// [Interface Func]
func (c *CheckCmd) Run() error {
	if err := c.initConfig(); err != nil {
		return err
	}
	return c.Check()
}

func (c *CheckCmd) initConfig() error {
	yhcConf := confdef.GetYHCConf()
	if err := confdef.InitMetricConf(yhcConf.MetricPaths); err != nil {
		return err
	}
	if err := confdef.InitModuleConf(yhcConf.DefaultModulePath); err != nil {
		return err
	}
	return nil
}
