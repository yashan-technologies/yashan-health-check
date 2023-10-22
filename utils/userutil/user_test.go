package userutil_test

import (
	"fmt"
	"testing"

	"yhc/log"
	"yhc/utils/userutil"
)

func TestG(t *testing.T) {
	initLogger("/tmp", "debug")
	user, e := userutil.GetUserOfGroup(log.Controller, "YASDBA")
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(user)
}

func initLogger(logPath, level string) error {
	optFuncs := []log.OptFunc{
		log.SetLogPath(logPath),
		log.SetLevel(level),
	}
	return log.InitLogger("_APP_NAME", log.NewLogOption(optFuncs...))
}
