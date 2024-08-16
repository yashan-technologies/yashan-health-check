package checkcontroller

import "yhc/defs/confdef"

type AfterInstallCmd struct {
	CheckGlobal
}

func (c *AfterInstallCmd) Run() error {
	if err := c.InitConfig(); err != nil {
		return err
	}
	return c.Check()
}

func (c *AfterInstallCmd) InitConfig() error {
	yhcConf := confdef.GetYHCConf()
	if err := confdef.InitMetricConf(yhcConf.AfterInstallMetricPath); err != nil {
		return err
	}
	if err := confdef.InitModuleConf(yhcConf.AfterInstallModulePath); err != nil {
		return err
	}
	return nil
}
