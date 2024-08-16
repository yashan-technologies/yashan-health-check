package main

import (
	"yhc/commons/flags"
	checkcontroller "yhc/internal/api/controller/yhcctlcontroller/check"
)

type App struct {
	flags.Globals
	Check        checkcontroller.CheckCmd        `cmd:"check" name:"check" help:"The check command is used to yashan health check."`
	AfterInstall checkcontroller.AfterInstallCmd `cmd:"after-install" name:"after-install" help:"The after-install command is used to verify the installation of Yashandb after it has been installed."`
}
