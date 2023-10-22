package yhccommons

import (
	"os"
	"strconv"

	"yhc/commons/yasdb"
	"yhc/defs/runtimedef"
	"yhc/utils/userutil"
	"yhc/utils/yasdbutil"

	"git.yasdb.com/go/yaslog"
)

func QueryYasdb(log yaslog.YasLog, db *yasdb.YashanDB, sql string, timeout int) ([]map[string]string, error) {
	yasdb := yasdbutil.NewYashanDB(log, db)
	return yasdb.QueryMultiRows(sql, timeout)
}

func ChownToExecuter(path string) error {
	if !userutil.IsCurrentUserRoot() {
		return nil
	}
	user := runtimedef.GetExecuter()
	uid, _ := strconv.ParseInt(user.Uid, 10, 64)
	gid, _ := strconv.ParseInt(user.Gid, 10, 64)
	if uid == userutil.ROOT_USER_UID {
		return nil
	}
	return os.Chown(path, int(uid), int(gid))
}
