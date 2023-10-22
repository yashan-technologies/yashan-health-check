// This is the main package for yhcctl.
// Yhcctl is used to manage the yashan health check.
package main

import (
	"fmt"
	"strings"

	"yhc/commons/flags"
	"yhc/commons/std"
	"yhc/defs/compiledef"
	"yhc/defs/confdef"
	"yhc/defs/runtimedef"
	"yhc/log"

	"git.yasdb.com/go/yaserr"
	"github.com/alecthomas/kong"
)

const (
	_APP_NAME        = "yhcctl"
	_APP_DESCRIPTION = "Yhcctl is used to manage the yashan health check."
)

func main() {
	var app App
	options := flags.NewAppOptions(_APP_NAME, _APP_DESCRIPTION, compiledef.GetAPPVersion())
	ctx := kong.Parse(&app, options...)
	if err := initApp(app); err != nil {
		ctx.FatalIfErrorf(err)
	}
	finalize := std.GetRedirecter().RedirectStd()
	defer finalize()
	std.WriteToFile(fmt.Sprintf("execute: %s %s\n", _APP_NAME, strings.Join(ctx.Args, " ")))
	if err := ctx.Run(); err != nil {
		fmt.Println(yaserr.Unwrap(err))
	}
}

func initLogger(logPath, level string) error {
	optFuncs := []log.OptFunc{
		log.SetLogPath(logPath),
		log.SetLevel(level),
	}
	return log.InitLogger(_APP_NAME, log.NewLogOption(optFuncs...))
}

func initApp(app App) error {
	if err := runtimedef.InitRuntime(); err != nil {
		return err
	}
	if err := confdef.InitYHCConf(app.Config); err != nil {
		return err
	}
	if err := initLogger(runtimedef.GetLogPath(), confdef.GetYHCConf().LogLevel); err != nil {
		return err
	}
	if err := std.InitRedirecter(); err != nil {
		return err
	}
	return nil
}
