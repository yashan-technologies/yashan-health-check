package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

const (
	CMD_QUERY = "query"
	CMD_EXEC  = "exec"
)

type App struct {
	Version kong.VersionFlag `help:"Show version."                                   short:"v"`
	Show    bool             `help:"[Hidden] Show software compilation information." hidden:"true" name:"show"`

	CmdType       string `help:"The type of command, one of (query,exec)"  name:"type"     short:"t"`
	YasdbUser     string `help:"The user of YashanDB"                      name:"user"     short:"u"`
	YasdbPassword string `help:"The password of YashanDB user"             name:"password" short:"p"`
	ListenAddr    string `help:"The listen address of YashanDB"            name:"addr"     short:"a"`
	DataPath      string `help:"the data path of YashanDB"                 name:"datapath" short:"d"`
	Sql           string `help:"the sql of yashandb"                       name:"sql"      short:"s"`
	Timeout       int    `help:"The timeout of connection, unit is second" name:"timeout"`
}

// [Interface Func]
func (c *App) Run() error {
	if err := c.Valid(); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	yasdb := NewYashanDB(c.YasdbUser, c.YasdbPassword, c.ListenAddr, c.DataPath)
	if c.CmdType == CMD_QUERY {
		res, err := yasdb.Query(c.Sql, c.Timeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		data, err := json.Marshal(res)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(data))
		return nil
	}
	if err := yasdb.ExecSQL(c.Sql, c.Timeout); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Println("succeed")
	return nil
}

func (c *App) Valid() error {
	c.TrimSpace()
	if len(c.Sql) == 0 {
		return fmt.Errorf("empty sql")
	}
	if c.CmdType != CMD_QUERY && c.CmdType != CMD_EXEC {
		return fmt.Errorf("invalid command type %s", c.CmdType)
	}
	if len(c.DataPath) == 0 && (len(c.YasdbUser) == 0 || len(c.YasdbPassword) == 0 || len(c.ListenAddr) == 0) {
		return fmt.Errorf("insufficient parameters entered")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("invalid timeout %d", c.Timeout)
	}
	return nil
}

func (c *App) TrimSpace() {
	c.Sql = strings.TrimSpace(c.Sql)
	c.CmdType = strings.TrimSpace(c.CmdType)
	c.DataPath = strings.TrimSpace(c.DataPath)
	c.YasdbUser = strings.TrimSpace(c.YasdbUser)
	c.YasdbPassword = strings.TrimSpace(c.YasdbPassword)
	c.ListenAddr = strings.TrimSpace(c.ListenAddr)
}
