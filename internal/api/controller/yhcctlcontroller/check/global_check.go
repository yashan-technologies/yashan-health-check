package checkcontroller

import (
	"errors"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"yhc/commons/std"
	"yhc/commons/yasdb"
	"yhc/defs/confdef"
	constdef "yhc/defs/constants"
	"yhc/defs/errdef"
	"yhc/defs/regexpdef"
	"yhc/defs/runtimedef"
	checkhandler "yhc/internal/api/handler/yhcctlhandler/check"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/fileutil"
	"yhc/utils/jsonutil"
	"yhc/utils/stringutil"
	"yhc/utils/timeutil"
	"yhc/utils/userutil"

	"git.yasdb.com/go/yasutil/fs"
	"git.yasdb.com/go/yasutil/tabler"
)

const (
	f_type   = "type"
	f_range  = "range"
	f_start  = "start"
	f_end    = "end"
	f_output = "output"

	range_help = "you must ensure that the number before (M|d|h|m) is greater than 0"
)

var (
	_examplesTime = []string{
		"yyyy-MM-dd",
		"yyyy-MM-dd-hh",
		"yyyy-MM-dd-hh-mm",
	}

	_examplesRange = []string{
		"1M",
		"1d",
		"1h",
		"1m",
	}
)

type CheckGlobal struct {
	Range              string `name:"range"               short:"r"          help:"The time range of the check, such as '1M', '1d', '1h', '1m'. If <range> is given, <start> and <end> will be discard."`
	Start              string `name:"start"               short:"s"          help:"The start datetime of the check, such as 'yyyy-MM-dd', 'yyyy-MM-dd-hh', 'yyyy-MM-dd-hh-mm'"`
	End                string `name:"end"                 short:"e"          help:"The end timestamp of the check, such as 'yyyy-MM-dd', 'yyyy-MM-dd-hh', 'yyyy-MM-dd-hh-mm', default value is current datetime."`
	Output             string `name:"output"              short:"o"          help:"The output dir of the check."`
	DisableInteraction bool   `name:"disable-interaction" short:"d"          help:"Disable interaction."`
	MultipleNodes      bool   `name:"multiple-nodes"      short:"m"          help:"Check multiple nodes."`
	YasdbHome          string `name:"yasdb-home"          help:"Home path of YashanDB(env: YASDB_HOME)."`
	YasdbData          string `name:"yasdb-data"          help:"Data path of YashanDB(env: YASDB_DATA)."`
	YasdbUser          string `name:"user"          short:"u"          help:"YashanDB user for checking."`
	YasdbPassword      string `name:"password"      short:"p"          help:"YashanDB user password for checking."`
}

func (c *CheckGlobal) Check() error {
	c.fillDefault()
	if err := c.validate(); err != nil {
		return err
	}
	log.Controller.Debugf("module report: %s", jsonutil.ToJSONString(confdef.GetModuleConf()))
	var modules []*constdef.ModuleMetrics
	yasdb, modules := c.getViewModels()
	globalYasdb = &YashanDB{
		YashanDB:    yasdb,
		Mutex:       sync.Mutex{},
		checkStatus: STATUS_NOT_CHECK,
	}
	c.fillYasdbFromFlags(globalYasdb.YashanDB)
	if c.DisableInteraction {
		// use yasql query LISTEN_ADDR and fill yasdb
		if err := fillListenAddrAndDBName(globalYasdb.YashanDB); err != nil {
			log.Controller.Errorf("fill listen addr err: %s", err.Error())
			return err
		}
		if c.MultipleNodes {
			fillNodeInfos(globalYasdb)
			if err := validateCheckedNodes(); err != nil {
				return errors.New("no node can be checked")
			}
		}
		validateMetrics(globalYasdb.YashanDB, modules)
		if len(moduleNoNeedCheckMetrics) != 0 {
			std.WriteToFile("the following metric will not be checked \n")
			noNeedStr := genNoNeedCheckMetricsStr()
			std.WriteToFile(noNeedStr)
		}
	} else {
		StartTerminalView(modules, globalYasdb, c.MultipleNodes)
	}
	// globalExitCode will be fill after terminal view exit
	if globalExitCode != EXIT_CONTINUE {
		return errors.New(exitCodeMap[globalExitCode])
	}
	checkerBase := c.genCheckBase(globalYasdb, c.MultipleNodes)
	// write user choose yashan health check to console.log
	c.writeUserChoose()
	// globalFilterModule will be fill after user choose metrics
	handler := checkhandler.NewCheckHandler(globalFilterModule, checkerBase)
	if err := handler.Check(); err != nil {
		return err
	}
	return nil
}

func trimSpace(s string) string {
	return strings.TrimSpace(s)
}

func (c *CheckGlobal) genUserChooseMetricsStr(modules []*constdef.ModuleMetrics) string {
	t := tabler.NewTable("",
		tabler.NewRowTitle("module", 25),
		tabler.NewRowTitle("module checked", 10),
		tabler.NewRowTitle("metric", 25),
		tabler.NewRowTitle("metric check", 10),
	)
	for _, module := range modules {
		moduleAlias := confdef.GetModuleAlias(module.Name)
		moduleChecked := strconv.FormatBool(module.Enabled)
		for i, metric := range module.Metrics {
			if i != 0 {
				moduleAlias = ""
				moduleChecked = ""
			}
			if err := t.AddColumn(moduleAlias, moduleChecked, metric.NameAlias, metric.Enabled); err != nil {
				log.Module.Errorf("add column err: %s", err.Error())
			}
		}
	}
	return t.String()
}

func (c *CheckGlobal) getStartAndEnd() (start time.Time, end time.Time, err error) {
	defer func() {
		end = end.Add(time.Minute)
	}()
	conf := confdef.GetYHCConf()
	defRange := conf.GetRange()
	// range
	if !stringutil.IsEmpty(c.Range) {
		start, end, err = c.getRangeFlagTime()
		if err != nil {
			return
		}
		return
	}
	// start or end
	if !stringutil.IsEmpty(c.Start) || !stringutil.IsEmpty(c.End) {
		start, end, err = c.getStartEndFlagTime(defRange)
		if err != nil {
			return
		}
		return
	}
	// no flag input with default
	end = time.Now()
	start = end.Add(-defRange)
	return
}

func (c *CheckGlobal) writeUserChoose() {
	std.WriteToFile("user choose module metric result: \n")
	userChooseStr := c.genUserChooseMetricsStr(globalFilterModule)
	std.WriteToFile(userChooseStr)
}

func (c *CheckGlobal) getRangeFlagTime() (start, end time.Time, err error) {
	var r time.Duration
	r, err = timeutil.GetDuration(c.Range)
	if err != nil {
		return
	}
	end = time.Now()
	start = end.Add(-r)
	return
}

func (c *CheckGlobal) getStartEndFlagTime(defRange time.Duration) (start, end time.Time, err error) {
	if !stringutil.IsEmpty(c.Start) {
		start, err = timeutil.GetTimeDivBySepa(c.Start, stringutil.STR_HYPHEN)
		if err != nil {
			return
		}
		// only start
		if stringutil.IsEmpty(c.End) {
			end = start.Add(defRange)
			return
		}
		// both start end
		end, err = timeutil.GetTimeDivBySepa(c.End, stringutil.STR_HYPHEN)
		if err != nil {
			return
		}
		return
	}
	// only end
	end, err = timeutil.GetTimeDivBySepa(c.End, stringutil.STR_HYPHEN)
	if err != nil {
		return
	}
	start = end.Add(-defRange)
	return
}

func (c *CheckGlobal) validateOutput() error {
	output := c.Output
	if !regexpdef.PathRegexp.Match([]byte(output)) {
		return errdef.ErrPathFormat
	}
	if !path.IsAbs(output) {
		output = path.Join(runtimedef.GetYHCHome(), output)
	}
	_, err := os.Stat(output)
	if err != nil {
		if os.IsPermission(err) {
			return errdef.NewErrPermissionDenied(userutil.CurrentUser, output)
		}
		if !os.IsNotExist(err) {
			return err
		}
		if err := fs.Mkdir(output); err != nil {
			log.Controller.Errorf("create output err: %s", err.Error())
			if os.IsPermission(err) {
				return errdef.NewErrPermissionDenied(userutil.CurrentUser, output)
			}
			return err
		}
	}
	return fileutil.CheckUserWrite(output)
}

func (c *CheckGlobal) fillDefault() {
	if stringutil.IsEmpty(c.Output) {
		c.Output = confdef.GetYHCConf().Output
	}
	if !path.IsAbs(c.Output) {
		c.Output = path.Join(runtimedef.GetYHCHome(), c.Output)
	}
	c.Output = path.Clean(c.Output)
}

func (c *CheckGlobal) fillYasdbFromFlags(yasdb *yasdb.YashanDB) {
	if len(c.YasdbHome) > 0 {
		yasdb.YasdbHome = c.YasdbHome
	}
	if len(c.YasdbData) > 0 {
		yasdb.YasdbData = c.YasdbData
	}
	yasdb.YasdbUser = c.YasdbUser
	yasdb.YasdbPassword = c.YasdbPassword
}

func (c *CheckGlobal) getViewModels() (*yasdb.YashanDB, []*constdef.ModuleMetrics) {
	metricConf := confdef.GetMetricConf()
	modules := c.transferToModuleMetric(metricConf)
	yasdb := newYasdb()
	return yasdb, modules
}

func (c *CheckGlobal) transferToModuleMetric(config *confdef.YHCMetricConfig) (modules []*constdef.ModuleMetrics) {
	log := log.Controller.M("transfer metric conf")
	modules = make([]*constdef.ModuleMetrics, 0)
	m := make(map[string]*constdef.ModuleMetrics)
	for _, metric := range config.Metrics {
		if !metric.Enabled {
			log.Debugf("metric %s disable, skip to transfer", metric.Name)
			continue
		}
		if _, ok := m[metric.ModuleName]; !ok {
			m[metric.ModuleName] = &constdef.ModuleMetrics{
				Name:    metric.ModuleName,
				Enabled: true,
				Metrics: make([]*confdef.YHCMetric, 0),
			}
		}
		m[metric.ModuleName].Metrics = append(m[metric.ModuleName].Metrics, metric)
	}
	for _, module := range define.Level1ModuleOrder {
		if _, ok := m[string(module)]; ok {
			modules = append(modules, m[string(module)])
		}
	}
	return modules
}

func (c *CheckGlobal) validate() error {
	if err := c.validateRange(); err != nil {
		return err
	}
	if err := c.validateStartAndEnd(); err != nil {
		return err
	}
	if err := c.validateOutput(); err != nil {
		return err
	}
	return nil
}

func (c *CheckGlobal) validateRange() error {
	conf := confdef.GetYHCConf()
	log.Controller.Debugf("conf: %s\n", jsonutil.ToJSONString(conf))
	log.Controller.Debugf("cmd: %s", jsonutil.ToJSONString(c))
	if stringutil.IsEmpty(c.Range) {
		return nil
	}
	if !regexpdef.RangeRegexp.MatchString(c.Range) {
		return errdef.NewErrYHCFlag(f_range, c.Range, _examplesRange, range_help)
	}
	minDuration, maxDuration, err := conf.GetMinAndMaxDuration()
	if err != nil {
		log.Controller.Errorf("get duration err: %s", err.Error())
		return err
	}
	log.Controller.Debugf("get min %s max %s", minDuration.String(), maxDuration.String())
	r, err := timeutil.GetDuration(c.Range)
	if err != nil {
		return err
	}
	if r > maxDuration {
		return errdef.NewGreaterMaxDur(conf.MaxDuration)
	}
	if r < minDuration {
		return errdef.NewLessMinDur(conf.MinDuration)
	}
	return nil
}

func (c *CheckGlobal) validateStartAndEnd() error {
	conf := confdef.GetYHCConf()
	var (
		startNotEmpty, endNotEmpty bool
		start, end                 time.Time
		err                        error
	)
	if !stringutil.IsEmpty(c.Start) {
		if !regexpdef.TimeRegexp.MatchString(c.Start) {
			return errdef.NewErrYHCFlag(f_start, c.Start, _examplesTime, "")
		}
		start, err = timeutil.GetTimeDivBySepa(c.Start, stringutil.STR_HYPHEN)
		if err != nil {
			return err
		}
		now := time.Now()
		if start.After(now) {
			return errdef.ErrStartShouldLessCurr
		}
		startNotEmpty = true
	}
	if !stringutil.IsEmpty(c.End) {
		if !regexpdef.TimeRegexp.MatchString(c.End) {
			return errdef.NewErrYHCFlag(f_end, c.End, _examplesTime, "")
		}
		end, err = timeutil.GetTimeDivBySepa(c.End, stringutil.STR_HYPHEN)
		if err != nil {
			return err
		}
		endNotEmpty = true
	}
	if startNotEmpty && endNotEmpty {
		minDuration, maxDuration, err := conf.GetMinAndMaxDuration()
		if err != nil {
			log.Controller.Errorf("get duration err: %s", err.Error())
			return err
		}
		if end.Before(start) {
			return errdef.ErrEndLessStart
		}
		r := end.Sub(start)
		if r > maxDuration {
			return errdef.NewGreaterMaxDur(conf.MaxDuration)
		}
		if r < minDuration {
			return errdef.NewLessMinDur(conf.MaxDuration)
		}
	}
	return nil
}

func (c *CheckGlobal) genCheckBase(db *YashanDB, multipleNodes bool) *define.CheckerBase {
	start, end, _ := c.getStartAndEnd()
	var nodes []*yasdb.NodeInfo
	for _, node := range db.Nodes {
		if node.Check && node.Connected {
			nodes = append(nodes, node)
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].NodeID < nodes[j].NodeID
	})
	return &define.CheckerBase{
		DBInfo:        db.YashanDB,
		Start:         start,
		End:           end,
		Output:        c.Output,
		NodeInfos:     nodes,
		MultipleNodes: multipleNodes,
	}
}
