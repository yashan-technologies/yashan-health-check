package main

import (
	"fmt"

	"yhc/commons/flags"
	"yhc/defs/compiledef"

	"git.yasdb.com/go/yaserr"
	"github.com/alecthomas/kong"
)

const (
	_APP_NAME        = "yasdb-go"
	_APP_DESCRIPTION = "Yasdb-go is used to query or exec sql in YashanDB."
)

func main() {
	var app App
	options := flags.NewAppOptions(_APP_NAME, _APP_DESCRIPTION, compiledef.GetAPPVersion())
	ctx := kong.Parse(&app, options...)
	if err := ctx.Run(); err != nil {
		fmt.Println(yaserr.Unwrap(err))
	}
}
