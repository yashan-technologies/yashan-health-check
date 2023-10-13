package check

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"yhc/commons/yasdb"
	"yhc/defs/bashdef"
	"yhc/defs/confdef"
	"yhc/defs/runtimedef"
	yhccommons "yhc/internal/modules/yhc/check/commons"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/execerutil"
	"yhc/utils/fileutil"
	"yhc/utils/osutil"
	"yhc/utils/stringutil"
	"yhc/utils/userutil"

	"git.yasdb.com/go/yaslog"
)

const (
	YAS_USER_LACK_AUTH               = "YAS-02213"
	YAS_TABLE_OR_VIEW_DOES_NOT_EXIST = "YAS-02012"
)

type checkFunc func(log yaslog.YasLog, db *yasdb.YashanDB, metric *confdef.YHCMetric) *define.NoNeedCheckMetric

type getPathFunc func(log yaslog.YasLog, db *yasdb.YashanDB) (string, error)

var (
	NeedCheckMetricMap = map[define.MetricName]struct{}{
		define.METRIC_YASDB_OBJECT_COUNT:                                     {},
		define.METRIC_YASDB_OBJECT_OWNER:                                     {},
		define.METRIC_YASDB_OBJECT_TABLESPACE:                                {},
		define.METRIC_YASDB_INDEX_BLEVEL:                                     {},
		define.METRIC_YASDB_INDEX_COLUMN:                                     {},
		define.METRIC_YASDB_INDEX_INVISIBLE:                                  {},
		define.METRIC_YASDB_INDEX_OVERSIZED:                                  {},
		define.METRIC_YASDB_INDEX_TABLE_INDEX_NOT_TOGETHER:                   {},
		define.METRIC_YASDB_SEQUENCE_NO_AVAILABLE:                            {},
		define.METRIC_YASDB_TASK_RUNNING:                                     {},
		define.METRIC_YASDB_PACKAGE_NO_PACKAGE_PACKAGE_BODY:                  {},
		define.METRIC_YASDB_SECURITY_LOGIN_MAXIMUM_LOGIN_ATTEMPTS:            {},
		define.METRIC_YASDB_SECURITY_USER_NO_OPEN:                            {},
		define.METRIC_YASDB_SECURITY_USER_WITH_SYSTEM_TABLE_PRIVILEGES:       {},
		define.METRIC_YASDB_SECURITY_USER_WITH_DBA_ROLE:                      {},
		define.METRIC_YASDB_SECURITY_USER_ALL_PRIVILEGE_OR_SYSTEM_PRIVILEGES: {},
		define.METRIC_YASDB_SECURITY_USER_USE_SYSTEM_TABLESPACE:              {},
		define.METRIC_YASDB_SECURITY_AUDIT_CLEANUP_TASK:                      {},
		define.METRIC_YASDB_SECURITY_AUDIT_FILE_SIZE:                         {},
		define.METRIC_YASDB_TABLESPACE:                                       {},

		define.METRIC_YASDB_RUN_LOG_DATABASE_CHANGES: {},
		define.METRIC_YASDB_RUN_LOG_ERROR:            {},
		define.METRIC_YASDB_ALERT_LOG_ERROR:          {},
		define.METRIC_HOST_SYSTEM_LOG_ERROR:          {},
	}

	NeedCheckMetricFuncMap = map[define.MetricName]checkFunc{
		define.METRIC_YASDB_OBJECT_COUNT:                                     checkDBAPrivileges,
		define.METRIC_YASDB_OBJECT_OWNER:                                     checkDBAPrivileges,
		define.METRIC_YASDB_OBJECT_TABLESPACE:                                checkDBAPrivileges,
		define.METRIC_YASDB_INDEX_BLEVEL:                                     checkDBAPrivileges,
		define.METRIC_YASDB_INDEX_COLUMN:                                     checkDBAPrivileges,
		define.METRIC_YASDB_INDEX_INVISIBLE:                                  checkDBAPrivileges,
		define.METRIC_YASDB_INDEX_OVERSIZED:                                  checkDBAPrivileges,
		define.METRIC_YASDB_INDEX_TABLE_INDEX_NOT_TOGETHER:                   checkDBAPrivileges,
		define.METRIC_YASDB_SEQUENCE_NO_AVAILABLE:                            checkDBAPrivileges,
		define.METRIC_YASDB_TASK_RUNNING:                                     checkDBAPrivileges,
		define.METRIC_YASDB_PACKAGE_NO_PACKAGE_PACKAGE_BODY:                  checkDBAPrivileges,
		define.METRIC_YASDB_SECURITY_LOGIN_MAXIMUM_LOGIN_ATTEMPTS:            checkDBAPrivileges,
		define.METRIC_YASDB_SECURITY_USER_NO_OPEN:                            checkDBAPrivileges,
		define.METRIC_YASDB_SECURITY_USER_WITH_SYSTEM_TABLE_PRIVILEGES:       checkDBAPrivileges,
		define.METRIC_YASDB_SECURITY_USER_WITH_DBA_ROLE:                      checkDBAPrivileges,
		define.METRIC_YASDB_SECURITY_USER_ALL_PRIVILEGE_OR_SYSTEM_PRIVILEGES: checkDBAPrivileges,
		define.METRIC_YASDB_SECURITY_USER_USE_SYSTEM_TABLESPACE:              checkDBAPrivileges,
		define.METRIC_YASDB_SECURITY_AUDIT_CLEANUP_TASK:                      checkAuditEnableAndDBA,
		define.METRIC_YASDB_SECURITY_AUDIT_FILE_SIZE:                         checkAuditEnableAndDBA,
		define.METRIC_YASDB_TABLESPACE:                                       checkDBAPrivileges,
		define.METRIC_YASDB_RUN_LOG_DATABASE_CHANGES:                         checkPermission,
		define.METRIC_YASDB_RUN_LOG_ERROR:                                    checkPermission,
		define.METRIC_YASDB_ALERT_LOG_ERROR:                                  checkPermission,
		define.METRIC_HOST_SYSTEM_LOG_ERROR:                                  checkPermission,
		define.METRIC_HOST_DMESG_LOG_ERROR:                                   checkDmesg,
	}

	metricPermissionPathMap = map[string]getPathFunc{
		string(define.METRIC_YASDB_RUN_LOG_DATABASE_CHANGES): getRunLogPath,
		string(define.METRIC_YASDB_RUN_LOG_ERROR):            getRunLogPath,
		string(define.METRIC_YASDB_ALERT_LOG_ERROR):          getAlertLogPath,
		string(define.METRIC_HOST_SYSTEM_LOG_ERROR):          systemLogPath,
	}
)

func (c *YHCChecker) CheckSarAccess() error {
	cmd := []string{
		"-c",
		bashdef.CMD_SAR,
		"-V",
	}
	exe := execerutil.NewExecer(log.Module)
	ret, _, stderr := exe.Exec(bashdef.CMD_BASH, cmd...)
	if ret != 0 {
		return fmt.Errorf("failed to check sar command, err: %s", stderr)
	}
	return nil
}

func checkDBAPrivileges(log yaslog.YasLog, db *yasdb.YashanDB, metric *confdef.YHCMetric) *define.NoNeedCheckMetric {
	sql := metric.SQL
	if stringutil.IsEmpty(sql) {
		sql = SQLMap[define.MetricName(metric.Name)]
	}
	if _, err := yhccommons.QueryYasdb(log, db, sql, confdef.GetYHCConf().SqlTimeout); err != nil {
		if strings.Contains(err.Error(), YAS_USER_LACK_AUTH) || strings.Contains(err.Error(), YAS_TABLE_OR_VIEW_DOES_NOT_EXIST) {
			return &define.NoNeedCheckMetric{
				Name:        metric.NameAlias,
				Error:       err,
				Description: "需要DBA权限",
			}
		}
		log.Warnf("pre check %s err: %s", metric.NameAlias, err.Error())
		return nil
	}
	return nil
}

func checkAuditEnable(log yaslog.YasLog, db *yasdb.YashanDB, metric *confdef.YHCMetric) *define.NoNeedCheckMetric {
	sql := "select value from v$parameter where name = 'UNIFIED_AUDITING'"
	res, err := yhccommons.QueryYasdb(log, db, sql, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		log.Errorf("query UNIFIED_AUDITING err: %s", err.Error())
		return nil
	}
	if len(res) == 0 {
		log.Warnf("parameter UNIFIED_AUDITING not exist")
		return nil
	}
	if res[0]["VALUE"] == "TRUE" {
		return nil
	}
	return &define.NoNeedCheckMetric{
		Name:        metric.NameAlias,
		Error:       errors.New("系统审计未打开"),
		Description: "需要打开系统审计",
	}
}

func checkAuditEnableAndDBA(log yaslog.YasLog, db *yasdb.YashanDB, metric *confdef.YHCMetric) *define.NoNeedCheckMetric {
	if res := checkAuditEnable(log, db, metric); res != nil {
		return res
	}
	return checkDBAPrivileges(log, db, metric)
}

func checkPermission(log yaslog.YasLog, db *yasdb.YashanDB, metric *confdef.YHCMetric) *define.NoNeedCheckMetric {
	pathFunc, ok := metricPermissionPathMap[metric.Name]
	if !ok {
		log.Errorf("get metric %s path func failed", metric.Name)
		return nil
	}
	path, err := pathFunc(log, db)
	if err != nil {
		log.Errorf("get metric %s path failed", metric.Name)
		return nil
	}
	err = fileutil.CheckAccess(path)
	if err == nil {
		return nil
	}
	res := &define.NoNeedCheckMetric{
		Name:  metric.NameAlias,
		Error: err,
	}
	if os.IsPermission(err) {
		res.Description = fmt.Sprintf("用户%s没有%s访问权限", userutil.CurrentUser, path)
		return res
	}
	if os.IsNotExist(err) {
		res.Description = fmt.Sprintf("文件:%s不存在", path)
		return res
	}
	res.Description = err.Error()
	return res
}

func checkDmesg(log yaslog.YasLog, db *yasdb.YashanDB, metric *confdef.YHCMetric) *define.NoNeedCheckMetric {
	if userutil.IsCurrentUserRoot() {
		return nil
	}
	if runtimedef.GetOSRelease().Id != osutil.KYLIN_ID {
		return nil
	}
	return &define.NoNeedCheckMetric{
		Name:        metric.NameAlias,
		Description: "麒麟系统执行dmesg命令需要root权限",
		Error:       errors.New("dmesg need root permission"),
	}
}

func getRunLogPath(log yaslog.YasLog, db *yasdb.YashanDB) (p string, err error) {
	sql := fmt.Sprintf("select * from v$parameter where name = '%s'", PARAMETER_RUN_LOG_FILE_PATH)
	res, err := yhccommons.QueryYasdb(log, db, sql, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		return
	}
	p = path.Join(strings.ReplaceAll(res[0]["VALUE"], "?", db.YasdbData), NAME_YASDB_RUN_LOG)
	return
}

func getAlertLogPath(log yaslog.YasLog, db *yasdb.YashanDB) (p string, err error) {
	p = path.Join(db.YasdbData, KEY_YASDB_LOG, KEY_YASDB_ALERT_LOG, NAME_YASDB_ALERT_LOG)
	return
}

func systemLogPath(log yaslog.YasLog, db *yasdb.YashanDB) (p string, err error) {
	p, err = getSystemLogName()
	return
}
