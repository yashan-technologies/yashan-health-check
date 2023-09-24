package reporter

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"time"

	"yhc/defs/bashdef"
	"yhc/defs/timedef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/execerutil"
	"yhc/utils/fileutil"

	"git.yasdb.com/go/yaslog"
	"git.yasdb.com/go/yasutil/fs"
)

type YHCReport struct {
	BeginTime time.Time                    `json:"beginTime"`
	EndTime   time.Time                    `json:"endTime"`
	CheckBase *define.CheckerBase          `json:"checkBase"`
	Modules   map[string]*define.YHCModule `json:"modules"`
}

func NewYHCReport(checkBase *define.CheckerBase) *YHCReport {
	return &YHCReport{
		CheckBase: checkBase,
		Modules:   make(map[string]*define.YHCModule),
	}
}

func (r *YHCReport) GenReport() (string, error) {
	log := log.Module.M("gen result")
	if err := r.mkdir(); err != nil {
		log.Errorf("mkdir err: %s", err.Error())
		return "", err
	}
	if err := r.genData(log); err != nil {
		log.Errorf("gen data err: %s", err.Error())
		return "", err
	}
	path, err := r.tarResult()
	if err != nil {
		log.Errorf("tar check result package err: %s", path)
	}
	return path, nil
}

func (r *YHCReport) genData(log yaslog.YasLog) error {
	dataJson := path.Join(r.getDataDir(), "data.json")
	result := r.getModuleResult()
	bytes, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Errorf("module data result marshal err: %s", err.Error())
		return err
	}
	if err := fileutil.WriteFile(dataJson, bytes); err != nil {
		log.Errorf("write data.json err: %s", err.Error())
		return err
	}
	return nil
}

func (r *YHCReport) getModuleResult() (res map[string]map[define.MetricName]*define.YHCItem) {
	res = make(map[string]map[define.MetricName]*define.YHCItem)
	for moduleName, moduleData := range r.Modules {
		res[moduleName] = moduleData.Items()
	}
	return
}

func (r *YHCReport) getPackageName() string {
	return fmt.Sprintf("yhc-%s", r.getTimeStr(r.BeginTime))
}

func (r *YHCReport) getPackageDir() string {
	return path.Join(r.CheckBase.Output, r.getPackageName())
}

func (r *YHCReport) getPackageTarName() string {
	return fmt.Sprintf("%s.tar.gz", r.getPackageName())
}

func (r *YHCReport) getDataDir() string {
	return path.Join(r.getPackageDir(), "data")
}

func (r *YHCReport) getTimeStr(t time.Time) string {
	return t.Format(timedef.TIME_FORMAT_IN_FILE)
}

func (r *YHCReport) mkdir() error {
	if err := fs.Mkdir(r.getPackageDir()); err != nil {
		return err
	}
	if err := fs.Mkdir(r.getDataDir()); err != nil {
		return err
	}
	return nil
}

func (r *YHCReport) tarResult() (path string, err error) {
	tarName := r.getPackageTarName()
	packageDir := r.getPackageName()
	command := fmt.Sprintf("cd %s;%s czvf %s %s;rm -rf %s", r.CheckBase.Output, bashdef.CMD_TAR, tarName, packageDir, packageDir)
	executer := execerutil.NewExecer(log.Logger)
	ret, _, stderr := executer.Exec(bashdef.CMD_BASH, "-c", command)
	if ret != 0 {
		err = errors.New(stderr)
		return
	}
	path = tarName
	return
}
