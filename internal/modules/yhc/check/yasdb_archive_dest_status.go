package check

import (
	"fmt"
	"strings"

	"yhc/defs/confdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/yasdbutil"

	"git.yasdb.com/go/yaserr"
	"git.yasdb.com/go/yaslog"
)

const (
	DATABASE_ROLE_PRIMARY = "PRIMARY"
	DATABASE_ROLE_STANDBY = "STANDBY"

	KEY_DATABASE_ROLE = "DATABASE_ROLE"
)

func (c *YHCChecker) GetYasdbArchiveDestStatus(name string) (err error) {
	log := log.Module.M(string(define.METRIC_YASDB_ARCHIVE_DEST_STATUS))
	role, err := c.getNodeRole(log)
	if err != nil {
		c.fillResults(&define.YHCItem{
			Name:  define.METRIC_YASDB_ARCHIVE_DEST_STATUS,
			Error: err.Error(),
		})
		return
	}
	if role == DATABASE_ROLE_PRIMARY {
		return c.GetPrimaryMultiRowData(string(define.METRIC_YASDB_ARCHIVE_DEST_STATUS))
	}
	return c.getYasdbReplicationStatus(log)
}

// 获取节点主备角色信息
func (c *YHCChecker) getNodeRole(log yaslog.YasLog) (string, error) {
	yasdb := yasdbutil.NewYashanDB(log, c.base.DBInfo)
	res, err := yasdb.QueryMultiRows(define.SQL_QUERY_DATABASE, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		return "", err
	}
	if len(res) == 0 {
		return "", fmt.Errorf("failed to get info by sql '%s'", define.SQL_QUERY_DATABASE)
	}
	return strings.ToUpper(res[0][KEY_DATABASE_ROLE]), nil
}

func (c *YHCChecker) getYasdbReplicationStatus(log yaslog.YasLog) (err error) {
	data := &define.YHCItem{Name: define.METRIC_YASDB_ARCHIVE_DEST_STATUS}
	defer c.fillResults(data)

	yasdb := yasdbutil.NewYashanDB(log, c.base.DBInfo)
	res, err := yasdb.QueryMultiRows(define.SQL_QUERY_REPLICATION_STATUS, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		data.Error = err.Error()
		return
	}
	metric, err := c.getMetric(define.METRIC_YASDB_ARCHIVE_DEST_STATUS)
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		data.Error = err.Error()
		return
	}
	data.Details = c.convertMultiSqlData(metric, res)
	return
}
