package main

import (
	"context"
	"fmt"
	"os/user"
	"time"

	"yhc/db"
	"yhc/utils/userutil"
)

const (
	_GROUP_YASDBA = "YASDBA"
)

const (
	INSTANCE_STATUS_SQL = "SELECT STATUS,VERSION FROM V$INSTANCE"
)

type YashanDB struct {
	YasdbUser     string
	YasdbPassword string
	ListenAddr    string
	DataPath      string
}

func NewYashanDB(yasdbUser, yasdbPassword, listenAddr, dataPath string) *YashanDB {
	return &YashanDB{
		YasdbUser:     yasdbUser,
		YasdbPassword: yasdbPassword,
		ListenAddr:    listenAddr,
		DataPath:      dataPath,
	}
}

func (y *YashanDB) Driver() (*db.YasDBDriver, error) {
	if y.isUdsOpen() {
		driver, err := y.udsDriver()
		if err != nil {
			return y.tcpDriver()
		}
		return driver, nil
	}
	return y.tcpDriver()
}

func (y *YashanDB) isUdsOpen() bool {
	u, err := user.Current()
	if err != nil {
		return false
	}
	gs := userutil.GetUserGroups(u)
	if len(gs) == 0 {
		return false
	}
	for _, g := range gs {
		if g == _GROUP_YASDBA {
			return true
		}
	}
	return false
}

func (y *YashanDB) udsDriver() (*db.YasDBDriver, error) {
	return db.NewYasDBDriver(db.WithDataPath(y.DataPath))
}

func (y *YashanDB) tcpDriver() (*db.YasDBDriver, error) {
	return db.NewYasDBDriver(
		db.WithUser(y.YasdbUser),
		db.WithPassword(y.YasdbPassword),
		db.WithAddr(y.ListenAddr),
	)
}

func (y *YashanDB) ExecSQL(query string, timeout int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	db, err := y.Driver()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.DB.ExecContext(ctx, query)
	return err
}

func (y *YashanDB) Query(query string, timeout int) ([]map[string]string, error) {
	var result []map[string]string
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	db, err := y.Driver()
	if err != nil {
		return result, err
	}
	defer db.Close()

	rows, err := db.DB.QueryContext(ctx, query)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		return result, err
	}
	cols, err := rows.Columns()
	if err != nil {
		return result, err
	}
	for rows.Next() {
		columnData := make([]interface{}, len(cols))
		scanArgs := make([]interface{}, len(cols))

		for i := range columnData {
			scanArgs[i] = &columnData[i]
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return result, err
		}
		mapValue := make(map[string]string)
		for i, colName := range cols {
			if columnData[i] == nil {
				mapValue[colName] = ""
				continue
			}
			mapValue[colName] = fmt.Sprint(columnData[i])
		}
		result = append(result, mapValue)
	}
	return result, nil
}
