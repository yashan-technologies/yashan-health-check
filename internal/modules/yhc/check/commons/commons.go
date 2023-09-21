package yhccommons

import (
	"yhc/commons/yasdb"
	"yhc/utils/yasdbutil"

	"git.yasdb.com/go/yaslog"
)

func QueryYasdb(log yaslog.YasLog, db *yasdb.YashanDB, sql string, timeout int) ([]map[string]string, error) {
	yasdb := yasdbutil.NewYashanDB(log, db)
	return yasdb.QueryMultiRows(sql, timeout)
}
