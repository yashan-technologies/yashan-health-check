package check

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"yhc/defs/confdef"
	"yhc/defs/timedef"
	"yhc/internal/modules/yhc/check/alertgenner"
	"yhc/internal/modules/yhc/check/define"
	"yhc/internal/modules/yhc/check/gopsutil"
	"yhc/internal/modules/yhc/check/jsonparser"
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
	define.METRIC_HOST_HISTORY_CPU_USAGE: define.WT_CPU,
	define.METRIC_HOST_CURRENT_CPU_USAGE: define.WT_CPU,
}

var SQLMap = map[define.MetricName]string{
	define.METRIC_YASDB_CONTROLFILE:        define.SQL_QUERY_CONTROLFILE,
	define.METRIC_YASDB_CONTROLFILE_COUNT:  define.SQL_QUERY_CONTROLFILE_COUNT,
	define.METRIC_YASDB_DATABASE:           define.SQL_QUERY_DATABASE,
	define.METRIC_YASDB_INSTANCE:           define.SQL_QUERY_INSTANCE,
	define.METRIC_YASDB_INDEX_BLEVEL:       define.SQL_QUERY_INDEX_BLEVEL,
	define.METRIC_YASDB_INDEX_COLUMN:       define.SQL_QUERY_INDEX_COLUMN,
	define.METRIC_YASDB_INDEX_INVISIBLE:    define.SQL_QUERY_INDEX_INVISIBLE,
	define.METRIC_YASDB_LISTEN_ADDR:        define.SQL_QUERY_LISTEN_ADDR,
	define.METRIC_YASDB_TABLESPACE:         define.SQL_QUERY_TABLESPACE,
	define.METRIC_YASDB_WAIT_EVENT:         define.SQL_QUERY_WAIT_EVENT,
	define.METRIC_YASDB_REPLICATION_STATUS: define.SQL_QUERY_REPLICATION_STATUS,
	define.METRIC_YASDB_PARAMETER:          define.SQL_QUERY_PARAMETER,
	define.METRIC_YASDB_OBJECT_COUNT:       define.SQL_QUERY_TOTAL_OBJECT,
	define.METRIC_YASDB_OBJECT_OWNER:       define.SQL_QUERY_OWNER_OBJECT,
	define.METRIC_YASDB_OBJECT_TABLESPACE:  define.SQL_QUERY_TABLESPACE_OBJECT,
	define.METRIC_YASDB_REDO_LOG:           define.SQL_QUERY_LOGFILE,
	define.METRIC_YASDB_REDO_LOG_COUNT:     define.SQL_QUERY_LOGFILE_COUNT,
}

type logTimeParseFunc func(date time.Time, line string) (time.Time, error)

type Checker interface {
	CheckFuncs(metrics []*confdef.YHCMetric) map[string]func() error
	GetResult(startCheck, endCheck time.Time) (map[define.MetricName]*define.YHCItem, *define.PandoraReport)
}

type YHCChecker struct {
	mtx        sync.RWMutex
	base       *define.CheckerBase
	metrics    []*confdef.YHCMetric
	Result     map[define.MetricName]*define.YHCItem
	FailedItem map[define.MetricName]error
}

func NewYHCChecker(base *define.CheckerBase, metrics []*confdef.YHCMetric) *YHCChecker {
	return &YHCChecker{
		base:       base,
		metrics:    metrics,
		Result:     map[define.MetricName]*define.YHCItem{},
		FailedItem: map[define.MetricName]error{},
	}
}

// [Interface Func]
func (c *YHCChecker) GetResult(startCheck, endCheck time.Time) (map[define.MetricName]*define.YHCItem, *define.PandoraReport) {
	c.fillterFailed()
	c.genAlerts()
	return c.Result, c.genReportJson(startCheck, endCheck)
}

func (c *YHCChecker) genReportJson(startCheck, endCheck time.Time) *define.PandoraReport {
	log := log.Module.M("gen-report-json")
	parser := jsonparser.NewJsonParser(log, *c.base, startCheck, endCheck, c.metrics, c.Result)
	return parser.Parse()
}

func (c *YHCChecker) genAlerts() {
	log := log.Module.M("gen-alert")
	alertGenner := alertgenner.NewAlterGenner(log, c.metrics, c.Result)
	c.Result = alertGenner.GenAlerts()
}

func (c *YHCChecker) fillterFailed() {
	for _, metric := range c.metrics {
		name := define.MetricName(metric.Name)
		item, ok := c.Result[name]
		if !ok {
			c.FailedItem[name] = fmt.Errorf("could not find result of metric %s", name)
			continue
		}
		if !stringutil.IsEmpty(item.Error) {
			delete(c.Result, name)
			c.FailedItem[name] = errors.New(item.Error)
		}
	}
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
		define.METRIC_HOST_INFO:                c.GetHostInfo,
		define.METRIC_HOST_CPU_INFO:            c.GetHostCPUInfo,
		define.METRIC_HOST_HISTORY_CPU_USAGE:   c.GetHostHistoryCPUUsage,
		define.METRIC_HOST_CURRENT_CPU_USAGE:   c.GetHostCurrentCPUUsage,
		define.METRIC_YASDB_CONTROLFILE:        c.GetYasdbControlFile,
		define.METRIC_YASDB_CONTROLFILE_COUNT:  c.GetYasdbControlFileCount,
		define.METRIC_YASDB_DATABASE:           c.GetYasdbDatabase,
		define.METRIC_YASDB_FILE_PERMISSION:    c.GetYasdbFilePermission,
		define.METRIC_YASDB_INDEX_BLEVEL:       c.GetYasdbIndexBlevel,
		define.METRIC_YASDB_INDEX_COLUMN:       c.GetYasdbIndexColumn,
		define.METRIC_YASDB_INDEX_INVISIBLE:    c.GetYasdbIndexInvisible,
		define.METRIC_YASDB_INSTANCE:           c.GetYasdbInstance,
		define.METRIC_YASDB_LISTEN_ADDR:        c.GetYasdbListenAddr,
		define.METRIC_YASDB_RUN_LOG_ERROR:      c.GetYasdbRunLogError,
		define.METRIC_YASDB_REDO_LOG:           c.GetYasdbRedoLog,
		define.METRIC_YASDB_REDO_LOG_COUNT:     c.GetYasdbRedoLogCount,
		define.METRIC_YASDB_OBJECT_COUNT:       c.GetYasdbObjectCount,
		define.METRIC_YASDB_OBJECT_OWNER:       c.GetYasdbOwnerObject,
		define.METRIC_YASDB_OBJECT_TABLESPACE:  c.GetYasdbTablespaceObject,
		define.METRIC_YASDB_REPLICATION_STATUS: c.GetYasdbReplicationStatus,
		define.METRIC_YASDB_PARAMETER:          c.GetYasdbParameter,
		define.METRIC_YASDB_SESSION:            c.GetYasdbSession,
		define.METRIC_YASDB_TABLESPACE:         c.GetYasdbTablespace,
		define.METRIC_YASDB_WAIT_EVENT:         c.GetYasdbWaitEvent,
	}
	return
}

func (c *YHCChecker) fillResult(data *define.YHCItem) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.Result[data.Name] = data
}

func (c *YHCChecker) querySingleRow(name define.MetricName) (*define.YHCItem, error) {
	data := &define.YHCItem{
		Name: name,
	}
	log := log.Module.M(string(name))
	sql, err := c.getSQL(name)
	if err != nil {
		log.Errorf("failed to get  sql of %s, err: %v", name, err)
		data.Error = err.Error()
		return data, err
	}
	metric, err := c.getMetric(name)
	if err != nil {
		log.Errorf("failed to get metric by name %s, err: %v", name, err)
		data.Error = err.Error()
		return data, err
	}
	yasdb := yasdbutil.NewYashanDB(log, c.base.DBInfo)
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
	data.Details = c.convertSqlData(metric, res[0])
	return data, nil
}

func (c *YHCChecker) queryMultiRows(name define.MetricName) (*define.YHCItem, error) {
	data := &define.YHCItem{
		Name: name,
	}
	log := log.Module.M(string(name))
	sql, err := c.getSQL(name)
	if err != nil {
		log.Errorf("failed to get  sql of %s, err: %v", name, err)
		data.Error = err.Error()
		return data, err
	}
	metric, err := c.getMetric(name)
	if err != nil {
		log.Errorf("failed to get metric by name %s, err: %v", name, err)
		data.Error = err.Error()
		return data, err
	}
	yasdb := yasdbutil.NewYashanDB(log, c.base.DBInfo)
	res, err := yasdb.QueryMultiRows(sql, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		log.Errorf("failed to get data with sql '%s', err: %v", sql, err)
		data.Error = err.Error()
		return data, err
	}
	data.Details = c.convertMultiSqlData(metric, res)
	return data, nil
}

func (c *YHCChecker) querySingleParameter(log yaslog.YasLog, name string) (string, error) {
	sql := fmt.Sprintf(SQL_QUERY_SINGLE_PARAMETER_FORMATER, name)
	yasdb := yasdbutil.NewYashanDB(log, c.base.DBInfo)
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
		if t.Before(c.base.Start) {
			continue
		}
		if t.After(c.base.End) {
			break
		}
		res = append(res, txt)
	}
	return res, nil
}

func (c *YHCChecker) hostHistoryWorkload(log yaslog.YasLog, name define.MetricName) (resp define.WorkloadOutput, err error) {
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
	args := c.genHistoryWorkloadArgs(c.base.Start, c.base.End, sarDir)
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

func (c *YHCChecker) convertObjectData(object interface{}) (res map[string]any, err error) {
	res = make(map[string]any)
	bytes, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &res)
	return
}

func (c *YHCChecker) convertMultiSqlData(metric *confdef.YHCMetric, datas []map[string]string) []map[string]interface{} {
	res := []map[string]interface{}{}
	for _, data := range datas {
		res = append(res, c.convertSqlData(metric, data))
	}
	return res
}

func (c *YHCChecker) convertSqlData(metric *confdef.YHCMetric, data map[string]string) map[string]interface{} {
	log := log.Module.M("convert-sql-data")
	res := make(map[string]interface{})
	for _, col := range metric.NumberColumns {
		value, ok := data[col]
		if !ok {
			log.Debugf("column %s not found, skip", col)
			continue
		}
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Errorf("failed to parse column %s to float64, value: %s, metric: %s, err: %v", col, value, metric.Name, err)
			continue
		}
		res[col] = f
	}
	for k, v := range data {
		if _, ok := res[k]; !ok {
			res[k] = v
		}
	}
	return res
}

func (c *YHCChecker) getMetric(name define.MetricName) (*confdef.YHCMetric, error) {
	for _, metric := range c.metrics {
		if metric.Name == string(name) {
			return metric, nil
		}
	}
	return nil, fmt.Errorf("failed to found metric by name %s", name)
}

func (c *YHCChecker) getSQL(name define.MetricName) (string, error) {
	metric, err := c.getMetric(name)
	if err != nil {
		return "", err
	}
	if stringutil.IsEmpty(metric.SQL) {
		return SQLMap[name], nil
	}
	return metric.SQL, nil
}
