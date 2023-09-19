package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "git.yasdb.com/go/yasdb-go"
)

type DriverOpts func(*YasDBDriver)

func WithUser(user string) DriverOpts { return func(d *YasDBDriver) { d.User = user } }

func WithPassword(pwd string) DriverOpts { return func(d *YasDBDriver) { d.Password = pwd } }

func WithAddr(addr string) DriverOpts { return func(d *YasDBDriver) { d.Addr = addr } }

func WithDataPath(datapath string) DriverOpts { return func(d *YasDBDriver) { d.DataPath = datapath } }

type YasDBDriver struct {
	YasdbHome string
	User      string
	Password  string
	Addr      string
	DataPath  string
	DB        *sql.DB
}

func NewYasDBDriver(opts ...DriverOpts) (*YasDBDriver, error) {
	d := &YasDBDriver{}
	for _, opt := range opts {
		opt(d)
	}
	if err := d.Connect(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *YasDBDriver) Connect() error {
	var db *sql.DB
	var err error
	if d.DataPath == "" {
		d.formatPassword()
		db, err = sql.Open("yasdb", fmt.Sprintf("%s/%s@%s", d.User, d.Password, d.Addr))
	} else {
		db, err = sql.Open("yasdb", d.DataPath)
	}
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	d.DB = db
	return nil
}

func (d *YasDBDriver) Close() {
	if d.DB != nil {
		d.DB.Close()
		d.DB = nil
	}
}

// formatPassword add '\' before '/', '@' or '\'.
func (d *YasDBDriver) formatPassword() {
	var newPwd strings.Builder

	for _, r := range d.Password {
		if r == '\\' || r == '@' || r == '/' {
			newPwd.WriteRune('\\')
		}
		newPwd.WriteRune(r)
	}
	d.Password = newPwd.String()
}
