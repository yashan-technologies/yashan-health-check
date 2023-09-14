package flags

import (
	"fmt"

	"yhc/defs/compiledef"

	"git.yasdb.com/go/yasutil/tabler"
	"github.com/alecthomas/kong"
)

type showFlag bool

// [Interface Func]
// BeforeReset shows software compilation information and terminates with a 0 exit status.
func (s showFlag) BeforeReset(app *kong.Kong, vars kong.Vars) error {
	fmt.Fprint(app.Stdout, s.genContent())
	app.Exit(0)
	return nil
}

// genContent generates data of software compilation information.
func (s showFlag) genContent() string {
	titles := []*tabler.RowTitle{
		{Name: "KEY"},
		{Name: "VALUE"},
	}
	table := tabler.NewTable("App Information", titles...)
	_ = table.AddColumn("App Version", compiledef.GetAPPVersion())
	_ = table.AddColumn("Go Version", compiledef.GetGoVersion())
	_ = table.AddColumn("Git Commit", compiledef.GetGitCommitID())
	_ = table.AddColumn("Git Describe", compiledef.GetGitDescribe())
	return table.String()
}
