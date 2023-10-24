package check

import (
	"yhc/defs/confdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/fileutil"
	"yhc/utils/yasdbutil"

	"git.yasdb.com/go/yaserr"
)

const (
	KEY_BACKUP_SET_PATH       = "PATH"
	KEY_BACKUP_SET_PATH_EXIST = "EXIST"

	STR_FALSE = "FALSE"
	STR_TRUE  = "TRUE"
)

func (c *YHCChecker) GetYasdbBackupSetPath(name string) (err error) {
	data := &define.YHCItem{Name: define.METRIC_YASDB_BACKUP_SET_PATH}
	defer c.fillResult(data)

	logger := log.Module.M(string(define.METRIC_YASDB_BACKUP_SET_PATH))
	yasdb := yasdbutil.NewYashanDB(logger, c.base.DBInfo)
	paths, err := yasdb.QueryMultiRows(define.SQL_QUERY_BACKUP_SET_PATH, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		err = yaserr.Wrap(err)
		logger.Error(err)
		data.Error = err.Error()
		return
	}

	var content []map[string]string
	if len(paths) == 0 {
		content = append(content, map[string]string{
			KEY_BACKUP_SET_PATH:       "NO BACKUP SET PATH FOUND",
			KEY_BACKUP_SET_PATH_EXIST: STR_FALSE,
		})
		data.Details = content
		return
	}

	for _, p := range paths {
		path := p[KEY_BACKUP_SET_PATH]
		exist := STR_FALSE
		if fileutil.IsExist(path) {
			exist = STR_TRUE
		}
		content = append(content, map[string]string{
			KEY_BACKUP_SET_PATH:       path,
			KEY_BACKUP_SET_PATH_EXIST: exist,
		})
	}
	data.Details = content
	return
}
