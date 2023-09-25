package check

import (
	"fmt"
	"path"
	"strings"
	"time"

	"yhc/defs/timedef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/stringutil"

	"git.yasdb.com/go/yaslog"
)

const (
	PARAMETER_RUN_LOG_FILE_PATH = "RUN_LOG_FILE_PATH"

	KEY_YASDB_RUN_LOG     = "run"
	KEY_YASDB_RUN_LOG_ERR = "errno"

	NAME_YASDB_RUN_LOG = "run.log"
)

func (c *YHCChecker) GetYasdbRedoLog() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_REDO_LOG)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbRedoLogCount() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_REDO_LOG_COUNT)
	log.Module.Error(data.Details)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbRunLogError() (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_YASDB_RUN_LOG_ERROR,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_YASDB_RUN_LOG_ERROR))
	var res []string
	runLogPath, err := c.getRunLogPath(log)
	if err != nil {
		log.Errorf("failed to get run log path, err: %v", err)
		data.Error = err.Error()
		return
	}
	runLogFiles, err := c.getLogFiles(log, runLogPath, KEY_YASDB_RUN_LOG)
	if err != nil {
		log.Error(err)
		data.Error = err.Error()
		return
	}
	if res, err = c.getYasdbRunLogError(log, runLogFiles); err != nil {
		log.Error(err)
		data.Error = err.Error()
		return
	}
	data.Details = res
	return
}

func (c *YHCChecker) getYasdbRunLogError(log yaslog.YasLog, srcs []string) (res []string, err error) {
	timeParseFunc := func(date time.Time, line string) (t time.Time, err error) {
		fields := strings.Split(line, stringutil.STR_BLANK_SPACE)
		if len(fields) < 2 {
			err = fmt.Errorf("invalid line: %s, skip", line)
			return
		}
		timeStr := fmt.Sprintf("%s %s", fields[0], fields[1])
		return time.ParseInLocation(timedef.TIME_FORMAT_WITH_MICROSECOND, timeStr, time.Local)
	}
	for _, f := range srcs {
		logEndTime := time.Now()
		if path.Base(f) != NAME_YASDB_RUN_LOG {
			fileds := strings.Split(strings.TrimSuffix(path.Base(f), ".log"), stringutil.STR_HYPHEN)
			if len(fileds) < 2 {
				log.Errorf("failed to get log end time from %s, skip", f)
				continue
			}
			if logEndTime, err = time.ParseInLocation(timedef.TIME_FORMAT_IN_FILE, fileds[1], time.Local); err != nil {
				log.Errorf("failed to parse log end time from %s", fileds[1])
				continue
			}
		}
		if logEndTime.Before(c.base.Start) {
			// no need to write into dest
			log.Debugf("skip run log file: %s", f)
			continue
		}
		if res, err = c.collectLogError(log, f, time.Now(), KEY_YASDB_RUN_LOG_ERR, timeParseFunc); err != nil {
			return
		}
	}
	return
}

func (c *YHCChecker) getRunLogPath(log yaslog.YasLog) (path string, err error) {
	path, err = c.querySingleParameter(log, PARAMETER_RUN_LOG_FILE_PATH)
	if err != nil {
		return
	}
	return strings.ReplaceAll(path, stringutil.STR_QUESTION_MARK, c.base.DBInfo.YasdbData), nil
}
