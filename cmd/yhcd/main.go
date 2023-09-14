package main

import (
	"yhc/commons/flags"
	"yhc/defs/compiledef"
	"yhc/defs/confdef"
	"yhc/defs/runtimedef"
	"yhc/log"

	"github.com/alecthomas/kong"
)

const (
	_APP_NAME        = "yhcd"
	_APP_DESCRIPTION = "yhcd is the daemon process of yashan health check"
)

func initApp(a App) error {
	if err := runtimedef.InitRuntime(); err != nil {
		return err
	}
	if err := confdef.InitYHCConf(a.Config); err != nil {
		return err
	}
	if err := log.InitLogger(_APP_NAME, log.NewLogOption()); err != nil {
		return err
	}
	return nil
}

func main() {
	var app App
	options := flags.NewAppOptions(_APP_NAME, _APP_DESCRIPTION, compiledef.GetAPPVersion())
	ctx := kong.Parse(&app, options...)
	if err := initApp(app); err != nil {
		ctx.FatalIfErrorf(err)
	}
	if err := ctx.Run(); err != nil {
		ctx.FatalIfErrorf(err)
	}
}
