package check

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"yhc/commons/yasdb"
	"yhc/defs/bashdef"
	"yhc/defs/confdef"
	"yhc/defs/timedef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/internal/modules/yhc/check/gopsutil"
	"yhc/internal/modules/yhc/check/sar"
	"yhc/log"
	"yhc/utils/stringutil"
	"yhc/utils/yasdbutil"

	"git.yasdb.com/go/yaslog"
)

const (
	SQL_QUERY_SINGLE_PARAMETER_FORMATER = "select value from v$parameter where name='%s'"
)

const (
	KEY_CURRENT = "current"
	KEY_HISTORY = "history"
)

var MetricNameToWorkloadTypeMap = map[define.MetricName]define.WorkloadType{
	define.METRIC_HOST_CPU_USAGE: define.WT_CPU,
}

type logTimeParseFunc func(date time.Time, line string) (time.Time, error)

type Checker interface {
	CheckFuncs(metrics []*confdef.YHCMetric) map[string]func() error
	GetResult() *define.YHCModule
}

type YHCChecker struct {
	Result     *define.YHCModule
	Yasdb      *yasdb.YashanDB
	FailedItem map[define.MetricName]error
	StartTime  time.Time
	EndTime    time.Time
}

func NewYHCChecker(base *define.CheckerBase) *YHCChecker {
	return &YHCChecker{
		Yasdb:     base.DBInfo,
		Result:    new(define.YHCModule),
		StartTime: base.Start,
		EndTime:   base.End,
	}
}

// [Interface Func]
func (c *YHCChecker) GetResult() *define.YHCModule {
	return c.Result
}

// [Interface Func]
func (c *YHCChecker) CheckFuncs(metrics []*confdef.YHCMetric) (res map[string]func() error) {
	res = make(map[string]func() error)
	defaultFuncMap := c.funcMap()
	for _, metric := range metrics {
		if metric.Default {
			fn, ok := defaultFuncMap[define.MetricName(metric.Name)]
			if !ok {
				log.Module.Errorf("failed to find function of default metric %s", metric.Name)
				continue
			}
			res[metric.Name] = fn
			continue
		}
		fn, err := c.GenCustomCheckFunc(metric)
		if err != nil {
			log.Module.Errorf("failed to gen function of custom metric %s", metric.Name)
			continue
		}
		res[metric.Name] = fn
	}
	return
}

func (c *YHCChecker) funcMap() (res map[define.MetricName]func() error) {
	res = map[define.MetricName]func() error{
		define.METRIC_HOST_INFO:               c.GetHostInfo,
		define.METRIC_HOST_CPU_INFO:           c.GetHostCPUInfo,
		define.METRIC_HOST_CPU_USAGE:          c.GetHostCPUUsage,
		define.METRIC_YASDB_CONTROLFILE:       c.GetYasdbControlFile,
		define.METRIC_YASDB_DATABASE:          c.GetYasdbDatabase,
		define.METRIC_YASDB_FILE_PERMISSION:   c.GetYasdbFilePermission,
		define.METRIC_YASDB_INDEX_BLEVEL:      c.GetYasdbIndexBlevel,
		define.METRIC_YASDB_INDEX_COLUMN:      c.GetYasdbIndexColumn,
		define.METRIC_YASDB_INDEX_INVISIBLE:   c.GetYasdbIndexInvisible,
		define.METRIC_YASDB_INSTANCE:          c.GetYasdbInstance,
		define.METRIC_YASDB_LISTEN_ADDR:       c.GetYasdbListenAddr,
		define.METRIC_YASDB_RUN_LOG_ERROR:     c.GetYasdbRunLogError,
		define.METRIC_YASDB_REDO_LOG:          c.GetYasdbRedoLog,
		define.METRIC_YASDB_OBJECT_COUNT:      c.GetYasdbObjectCount,
		define.METRIC_YASDB_OBJECT_OWNER:      c.GetYasdbOwnerObject,
		define.METRIC_YASDB_OBJECT_TABLESPACE: c.GetYasdbTablespaceObject,
		define.METRIC_YASDB_PARAMETER:         c.GetYasdbParameter,
		define.METRIC_YASDB_SESSION:           c.GetYasdbSession,
		define.METRIC_YASDB_TABLESPACE:        c.GetYasdbTablespace,
		define.METRIC_YASDB_WAIT_EVENT:        c.GetYasdbWaitEvent,
	}
	return
}

func (c *YHCChecker) fillResult(data *define.YHCItem) {
	c.Result.Set(data)
}

func (c *YHCChecker) querySingleRow(name define.MetricName, sql string) (*define.YHCItem, error) {
	data := &define.YHCItem{
		Name: name,
	}
	log := log.Module.M(string(name))
	yasdb := yasdbutil.NewYashanDB(log, c.Yasdb)
	res, err := yasdb.QueryMultiRows(sql, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		log.Errorf("failed to get data with sql %s, err: %v", sql, err)
		data.Error = err.Error()
		return data, err
	}
	if len(res) == 0 {
		err = fmt.Errorf("failed to get data with sql %s", sql)
		log.Error(err)
		data.Error = err.Error()
		return data, err
	}
	data.Details = res[0]
	return data, nil
}

func (c *YHCChecker) queryMultiRows(name define.MetricName, sql string) (*define.YHCItem, error) {
	data := &define.YHCItem{
		Name: name,
	}
	log := log.Module.M(string(name))
	yasdb := yasdbutil.NewYashanDB(log, c.Yasdb)
	res, err := yasdb.QueryMultiRows(sql, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		log.Errorf("failed to get data with sql '%s', err: %v", sql, err)
		data.Error = err.Error()
		return data, err
	}
	data.Details = res
	return data, nil
}

func (c *YHCChecker) querySingleParameter(log yaslog.YasLog, name string) (string, error) {
	sql := fmt.Sprintf(SQL_QUERY_SINGLE_PARAMETER_FORMATER, name)
	yasdb := yasdbutil.NewYashanDB(log, c.Yasdb)
	res, err := yasdb.QueryMultiRows(sql, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		log.Errorf("failed to query value of parameter %s, err: %v", name, err)
		return "", err
	}
	if len(res) <= 0 {
		err = fmt.Errorf("failed to query value of parameter %s", name)
		log.Error(err)
		return "", err
	}
	return res[0]["VALUE"], nil
}

func (c *YHCChecker) getLogFiles(log yaslog.YasLog, logPath string, prefix string) (logFiles []string, err error) {
	entrys, err := os.ReadDir(logPath)
	if err != nil {
		log.Error(err)
		return
	}
	for _, entry := range entrys {
		if !entry.Type().IsRegular() || !strings.HasPrefix(entry.Name(), prefix) {
			continue
		}
		logFiles = append(logFiles, path.Join(logPath, entry.Name()))
	}
	// sort with file name
	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i] < logFiles[j]
	})
	return
}

func (c *YHCChecker) collectLogError(log yaslog.YasLog, src string, date time.Time, errKey string, timeParseFunc logTimeParseFunc) ([]string, error) {
	res := []string{}
	srcFile, err := os.Open(src)
	if err != nil {
		return res, err
	}
	defer srcFile.Close()

	var t time.Time
	scanner := bufio.NewScanner(srcFile)
	for scanner.Scan() {
		txt := scanner.Text()
		line := stringutil.RemoveExtraSpaces(strings.TrimSpace(txt))
		if stringutil.IsEmpty(line) && strings.Contains(line, errKey) {
			continue
		}
		if t, err = timeParseFunc(date, line); err != nil {
			log.Error("skip line: %s, err: %s", txt, err.Error())
			continue
		}
		if t.Before(c.StartTime) {
			continue
		}
		if t.After(c.EndTime) {
			break
		}
		res = append(res, txt)
	}
	return res, nil
}

func (c *YHCChecker) hostHistoryWorkload(log yaslog.YasLog, name define.MetricName, start, end time.Time) (resp define.WorkloadOutput, err error) {
	// get sar args
	workloadType, ok := MetricNameToWorkloadTypeMap[name]
	if !ok {
		err = fmt.Errorf("failed to get workload type from metric name: %s", name)
		log.Error(err)
		return
	}
	sarArg, ok := sar.WorkloadTypeToSarArgMap[workloadType]
	if !ok {
		err = fmt.Errorf("failed to get SAR arg from workload type: %s", workloadType)
		log.Error(err)
		return
	}
	// collect
	sarCollector := sar.NewSar(log)
	sarDir := confdef.GetYHCConf().GetSarDir()
	if stringutil.IsEmpty(sarDir) {
		sarDir = sarCollector.GetSarDir()
	}
	sarOutput := make(define.WorkloadOutput)
	args := c.genHistoryWorkloadArgs(start, end, sarDir)
	for _, arg := range args {
		output, e := sarCollector.Collect(workloadType, sarArg, arg)
		if e != nil {
			log.Error(e)
			continue
		}
		for timestamp, output := range output {
			sarOutput[timestamp] = output
		}
	}
	resp = sarOutput
	return
}

func (c *YHCChecker) genHistoryWorkloadArgs(start, end time.Time, sarDir string) (args []string) {
	// get data between start and end
	var dates []time.Time
	begin := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	for date := begin; !date.After(end); date = date.AddDate(0, 0, 1) {
		dates = append(dates, date)
	}
	for i, date := range dates {
		var startArg, endArg, fileArg string
		// the frist
		if i == 0 && !date.Equal(start) {
			startArg = fmt.Sprintf("-s %s", start.Format(timedef.TIME_FORMAT_TIME))
		}
		// the last one
		if i == len(dates)-1 {
			if date.Equal(end) {
				// skip
				continue
			}
			endArg = fmt.Sprintf("-e %s", end.Format(timedef.TIME_FORMAT_TIME))
		}
		fileArg = fmt.Sprintf("-f %s", path.Join(sarDir, fmt.Sprintf("sa%s", date.Format(timedef.TIME_FORMAT_DAY))))
		args = append(args, fmt.Sprintf("%s %s %s", fileArg, startArg, endArg))
	}
	return
}

func (c *YHCChecker) hostCurrentWorkload(log yaslog.YasLog, name define.MetricName, hasSar bool) (resp define.WorkloadOutput, err error) {
	// global conf
	scrapeInterval, scrapeTimes := confdef.GetYHCConf().GetScrapeInterval(), confdef.GetYHCConf().GetScrapeTimes()
	// get sar args
	workloadType, ok := MetricNameToWorkloadTypeMap[name]
	if !ok {
		err = fmt.Errorf("failed to get workload type from metric name: %s", name)
		log.Error(err)
		return
	}
	if !hasSar {
		// use gopsutil to calculate by ourself
		return gopsutil.Collect(workloadType, scrapeInterval, scrapeTimes)
	}
	sarArg, ok := sar.WorkloadTypeToSarArgMap[workloadType]
	if !ok {
		err = fmt.Errorf("failed to get SAR arg from workload type: %s", workloadType)
		log.Error(err)
		return
	}
	sarCollector := sar.NewSar(log)
	return sarCollector.Collect(workloadType, sarArg, strconv.Itoa(scrapeInterval), strconv.Itoa(scrapeTimes))

}

func (c *YHCChecker) hostWorkload(log yaslog.YasLog, name define.MetricName) (resp define.HostWorkResponse, err error) {
	details := map[string]interface{}{}
	hasSar := c.CheckSarAccess() == nil
	resp.DataType = define.DATATYPE_GOPSUTIL
	resp.Errors = make(map[string]string)

	// collect historyworkload
	if hasSar {
		resp.DataType = define.DATATYPE_SAR
		if historyNetworkWorkload, e := c.hostHistoryWorkload(log, name, c.StartTime, c.EndTime); e != nil {
			err = fmt.Errorf("failed to collect history %s, err: %s", name, e.Error())
			resp.Errors[KEY_HISTORY] = err.Error()
			log.Error(err)
		} else {
			details[KEY_HISTORY] = historyNetworkWorkload
		}
	} else {
		e := fmt.Errorf("cannot find command '%s'", bashdef.CMD_SAR)
		resp.Errors[KEY_HISTORY] = e.Error()
		log.Error(e)
	}
	// collect current workload
	if currentNetworkWorkload, e := c.hostCurrentWorkload(log, name, hasSar); e != nil {
		err = fmt.Errorf("failed to collect current %s, err: %s", name, e.Error())
		resp.Errors[KEY_CURRENT] = err.Error()
		log.Error(err)
	} else {
		details[KEY_CURRENT] = currentNetworkWorkload
	}
	resp.Data = details
	return
}
