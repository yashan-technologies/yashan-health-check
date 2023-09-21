package checkcontroller

import (
	"errors"
	"os"
	"path"
	"strings"
	"time"

	"yhc/commons/yasdb"
	"yhc/defs/confdef"
	constdef "yhc/defs/constants"
	checkhandler "yhc/internal/api/handler/yhcctlhandler/check"
	"yhc/internal/modules/yhc/check/define"
	yhcyasdb "yhc/internal/modules/yhc/yasdb"
	"yhc/log"
	"yhc/utils/processutil"
	"yhc/utils/stringutil"
	"yhc/utils/timeutil"
	"yhc/utils/yasqlutil"
)

type CheckGlobal struct {
	Range  string `name:"range"  short:"r" help:"The time range of the check, such as '1M', '1d', '1h', '1m'. If <range> is given, <start> and <end> will be discard."`
	Start  string `name:"start"  short:"s" help:"The start datetime of the check, such as 'yyyy-MM-dd', 'yyyy-MM-dd-hh', 'yyyy-MM-dd-hh-mm'"`
	End    string `name:"end"    short:"e" help:"The end timestamp of the check, such as 'yyyy-MM-dd', 'yyyy-MM-dd-hh', 'yyyy-MM-dd-hh-mm', default value is current datetime."`
	Output string `name:"output" short:"o" help:"The output dir of the check."`
}

type CheckCmd struct {
	CheckGlobal
}

// [Interface Func]
func (c *CheckCmd) Run() error {
	c.fillDefault()
	if err := c.validate(); err != nil {
		return err
	}

	metricConf := confdef.GetMetricConf()
	modules := c.transferToModuleMetric(metricConf)
	yasdb := c.newYasdb()
	StartTerminalView(modules, yasdb)
	if globalExitCode != EXIT_CONTINUE {
		return errors.New(exitCodeMap[globalExitCode])
	}
	if err := c.fillListenAddr(yasdb); err != nil {
		log.Controller.Errorf("fill listen addr err: %s", err.Error())
		return err
	}
	checkerBase := c.genCheckBase(yasdb)
	handler := checkhandler.NewCheckHandler(globalFilterModule, checkerBase)
	if err := handler.Check(); err != nil {
		return err
	}
	return nil
}

func (c *CheckCmd) transferToModuleMetric(config *confdef.YHCMetricConfig) (modules []*constdef.ModuleMetrics) {
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

func (c *CheckCmd) yasdbPath() (yasdbHome, yasdbData string) {
	yasdbData = os.Getenv(constdef.YASDB_DATA)
	yasdbHome = os.Getenv(constdef.YASDB_HOME)
	processYasdbHome, processYasdbData := c.yasdbPathFromProcess()
	if stringutil.IsEmpty(yasdbHome) {
		yasdbHome = processYasdbHome
	}
	if stringutil.IsEmpty(yasdbData) {
		yasdbData = processYasdbData
	}
	return
}

func (c *CheckCmd) yasdbPathFromProcess() (yasdbHome, yasdbData string) {
	log := log.Controller.M("get yasdb process from cmdline")
	processes, err := processutil.ListAnyUserProcessByCmdline(_base_yasdb_process_format, true)
	if err != nil {
		log.Errorf("get process err: %s", err.Error())
		return
	}
	if len(processes) == 0 {
		log.Infof("process result is empty")
		return
	}
	for _, p := range processes {
		fields := strings.Split(p.ReadableCmdline, "-D")
		if len(fields) < 2 {
			log.Infof("process cmdline: %s format err, skip", p.ReadableCmdline)
			continue
		}
		yasdbData = trimSpace(fields[1])
		full := trimSpace(p.FullCommand)
		if !path.IsAbs(full) {
			return
		}
		yasdbHome = path.Dir(path.Dir(full))
		return
	}
	return
}

func (c *CheckCmd) newYasdb() *yasdb.YashanDB {
	home, data := c.yasdbPath()
	yasdb := &yasdb.YashanDB{
		YasdbData: data,
		YasdbHome: home,
	}
	return yasdb
}

func (c *CheckCmd) fillListenAddr(db *yasdb.YashanDB) error {
	tx := yasqlutil.GetLocalInstance(db.YasdbUser, db.YasdbPassword, db.YasdbHome, db.YasdbData)
	listenAddr, err := yhcyasdb.QueryParameter(tx, yhcyasdb.LISTEN_ADDR)
	if err != nil {
		return err
	}
	db.ListenAddr = trimSpace(listenAddr)
	return nil
}

func (c *CheckCmd) genCheckBase(db *yasdb.YashanDB) *define.CheckerBase {
	start, end, _ := c.getStartAndEnd()
	return &define.CheckerBase{
		DBInfo: db,
		Start:  start,
		End:    end,
		Output: c.Output,
	}
}

func trimSpace(s string) string {
	return strings.TrimSpace(s)
}

func (c *CheckCmd) getStartAndEnd() (start time.Time, end time.Time, err error) {
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

func (c *CheckCmd) getRangeFlagTime() (start, end time.Time, err error) {
	var r time.Duration
	r, err = timeutil.GetDuration(c.Range)
	if err != nil {
		return
	}
	end = time.Now()
	start = end.Add(-r)
	return
}

func (c *CheckCmd) getStartEndFlagTime(defRange time.Duration) (start, end time.Time, err error) {
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
