package flags

import (
	"fmt"

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
	return table.String()
}
