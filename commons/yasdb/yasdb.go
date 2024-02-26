package yasdb

import (
	"os"

	constdef "yhc/defs/constants"
	"yhc/defs/errdef"
	"yhc/utils/stringutil"

	"git.yasdb.com/go/yaslog"
	"git.yasdb.com/pandora/yasqlgo"
)

type NodeInfo struct {
	DatabaseName string
	NodeID       string
	ListenAddr   string
	Role         string
	Connected    bool
	User         string
	Password     string
	Check        bool
}

type YashanDB struct {
	YasdbHome     string
	YasdbData     string
	YasdbUser     string
	YasdbPassword string
	ListenAddr    string
	DatabaseName  string
}

func (y *YashanDB) ValidHome() error {
	if stringutil.IsEmpty(y.YasdbHome) {
		return errdef.NewItemEmpty(constdef.YASDB_HOME)
	}
	if err := y.validatePath(y.YasdbHome); err != nil {
		return err
	}
	return nil
}

func (y *YashanDB) ValidData() error {
	if !stringutil.IsEmpty(y.YasdbData) {
		if err := y.validatePath(y.YasdbData); err != nil {
			return err
		}
	}
	return nil
}

func (y *YashanDB) ValidUser() error {
	if stringutil.IsEmpty(y.YasdbUser) {
		return errdef.NewItemEmpty(constdef.YASDB_USER)
	}
	return nil
}

func (y *YashanDB) ValidPassword() error {
	if stringutil.IsEmpty(y.YasdbPassword) {
		return errdef.NewItemEmpty(constdef.YASDB_PASSWORD)
	}
	return nil
}

func (y *YashanDB) ValidUserAndPwd(log yaslog.YasLog) error {
	if err := y.ValidHome(); err != nil {
		return err
	}
	if err := y.ValidData(); err != nil {
		return err
	}
	if err := y.ValidUser(); err != nil {
		return err
	}
	if err := y.ValidPassword(); err != nil {
		return err
	}
	tx := yasqlgo.NewLocalInstance(y.YasdbUser, y.YasdbPassword, y.YasdbHome, y.YasdbData, log)
	if err := tx.CheckPassword(); err != nil {
		return err
	}
	return nil
}

func (y *YashanDB) validatePath(path string) error {
	_, err := os.Stat(path)
	return err
}
