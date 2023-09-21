package check

import (
	"fmt"

	"yhc/commons/yasdb"
	"yhc/defs/bashdef"
	"yhc/defs/confdef"
	yhccommons "yhc/internal/modules/yhc/check/commons"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/execerutil"

	"git.yasdb.com/go/yaslog"
)

type checkFunc func(log yaslog.YasLog, db *yasdb.YashanDB, metric *confdef.YHCMetric) *define.NoNeedCheckMetric

var (
	NeedCheckMetricMap = map[define.MetricName]struct{}{
		define.METRIC_YASDB_OBJECT_COUNT:      {},
		define.METRIC_YASDB_OBJECT_OWNER:      {},
		define.METRIC_YASDB_OBJECT_TABLESPACE: {},
		define.METRIC_YASDB_INDEX_BLEVEL:      {},
		define.METRIC_YASDB_INDEX_COLUMN:      {},
		define.METRIC_YASDB_INDEX_INVISIBLE:   {},
		define.METRIC_YASDB_TABLESPACE:        {},
	}

	NeedCheckMetricFuncMap = map[define.MetricName]checkFunc{
		define.METRIC_YASDB_OBJECT_COUNT:      checkYasdbObjectCount,
		define.METRIC_YASDB_OBJECT_OWNER:      checkYasdbObjectOwner,
		define.METRIC_YASDB_OBJECT_TABLESPACE: checkYasdbObjectTablespace,
		define.METRIC_YASDB_INDEX_BLEVEL:      checkYasdbIndexBlevel,
		define.METRIC_YASDB_INDEX_COLUMN:      checkYasdbIndexColumn,
		define.METRIC_YASDB_INDEX_INVISIBLE:   checkYasdbIndexInvisible,
		define.METRIC_YASDB_TABLESPACE:        checkYasdbTableSpace,
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

func checkYasdbObjectCount(log yaslog.YasLog, db *yasdb.YashanDB, metric *confdef.YHCMetric) *define.NoNeedCheckMetric {
	if _, err := yhccommons.QueryYasdb(log, db, SQL_QUERY_TOTAL_OBJECT, confdef.GetYHCConf().SqlTimeout); err != nil {
		return &define.NoNeedCheckMetric{
			Name:        metric.NameAlias,
			Error:       err,
			Description: fmt.Sprintf("%s need dba privileges", metric.NameAlias),
		}
	}
	return nil
}

func checkYasdbObjectTablespace(log yaslog.YasLog, db *yasdb.YashanDB, metric *confdef.YHCMetric) *define.NoNeedCheckMetric {
	if _, err := yhccommons.QueryYasdb(log, db, SQL_QUERY_TABLESPACE_OBJECT, confdef.GetYHCConf().SqlTimeout); err != nil {
		return &define.NoNeedCheckMetric{
			Name:        metric.NameAlias,
			Error:       err,
			Description: fmt.Sprintf("%s need dba privileges", metric.NameAlias),
		}
	}
	return nil
}

func checkYasdbObjectOwner(log yaslog.YasLog, db *yasdb.YashanDB, metric *confdef.YHCMetric) *define.NoNeedCheckMetric {
	if _, err := yhccommons.QueryYasdb(log, db, SQL_QUERY_OWNER_OBJECT, confdef.GetYHCConf().SqlTimeout); err != nil {
		return &define.NoNeedCheckMetric{
			Name:        metric.NameAlias,
			Error:       err,
			Description: fmt.Sprintf("%s need dba privileges", metric.NameAlias),
		}
	}
	return nil
}

func checkYasdbIndexBlevel(log yaslog.YasLog, db *yasdb.YashanDB, metric *confdef.YHCMetric) *define.NoNeedCheckMetric {
	if _, err := yhccommons.QueryYasdb(log, db, SQL_QUERY_INDEX_BLEVEL, confdef.GetYHCConf().SqlTimeout); err != nil {
		return &define.NoNeedCheckMetric{
			Name:        metric.NameAlias,
			Error:       err,
			Description: fmt.Sprintf("%s need dba privileges", metric.NameAlias),
		}
	}
	return nil
}

func checkYasdbIndexColumn(log yaslog.YasLog, db *yasdb.YashanDB, metric *confdef.YHCMetric) *define.NoNeedCheckMetric {
	if _, err := yhccommons.QueryYasdb(log, db, SQL_QUERY_INDEX_COLUMN, confdef.GetYHCConf().SqlTimeout); err != nil {
		return &define.NoNeedCheckMetric{
			Name:        metric.NameAlias,
			Error:       err,
			Description: fmt.Sprintf("%s need dba privileges", metric.NameAlias),
		}
	}
	return nil
}

func checkYasdbIndexInvisible(log yaslog.YasLog, db *yasdb.YashanDB, metric *confdef.YHCMetric) *define.NoNeedCheckMetric {
	if _, err := yhccommons.QueryYasdb(log, db, SQL_QUERY_INDEX_INVISIBLE, confdef.GetYHCConf().SqlTimeout); err != nil {
		return &define.NoNeedCheckMetric{
			Name:        metric.NameAlias,
			Error:       err,
			Description: fmt.Sprintf("%s need dba privileges", metric.NameAlias),
		}
	}
	return nil
}

func checkYasdbTableSpace(log yaslog.YasLog, db *yasdb.YashanDB, metric *confdef.YHCMetric) *define.NoNeedCheckMetric {
	if _, err := yhccommons.QueryYasdb(log, db, SQL_QUERY_TABLESPACE, confdef.GetYHCConf().SqlTimeout); err != nil {
		return &define.NoNeedCheckMetric{
			Name:        metric.NameAlias,
			Error:       err,
			Description: fmt.Sprintf("%s need dba privileges", metric.NameAlias),
		}
	}
	return nil
}
