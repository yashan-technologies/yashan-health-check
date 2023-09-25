package reporter

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"time"

	"yhc/defs/bashdef"
	"yhc/defs/timedef"
	yhccommons "yhc/internal/modules/yhc/check/commons"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/execerutil"
	"yhc/utils/fileutil"

	"git.yasdb.com/go/yasutil/fs"
)

const (
	_PACKAGE_NAME_FORMATTER = "ytc-%s"
	_DATA_NAME_FORMATTER    = "data-%s.json"
)

type YHCReport struct {
	BeginTime time.Time                             `json:"beginTime"`
	EndTime   time.Time                             `json:"endTime"`
	CheckBase *define.CheckerBase                   `json:"checkBase"`
	Items     map[define.MetricName]*define.YHCItem `json:"items"`
}

func NewYHCReport(checkBase *define.CheckerBase) *YHCReport {
	return &YHCReport{
		CheckBase: checkBase,
		Items:     map[define.MetricName]*define.YHCItem{},
	}
}

func (r *YHCReport) GenResult() (string, error) {
	log := log.Module.M("gen-result")
	if err := r.mkdir(); err != nil {
		log.Errorf("mkdir err: %s", err.Error())
		return "", err
	}
	if err := r.genDataJson(); err != nil {
		log.Errorf("gen data err: %s", err.Error())
		return "", err
	}
	// TODO write report
	if err := r.tarResult(); err != nil {
		log.Errorf("tar result failed: %s", err)
		return "", err
	}
	if err := r.chownResult(); err != nil {
		log.Errorf("chown result failed: %s", err)
	}
	return r.genPackageTarPath(), nil
}

func (r *YHCReport) genDataJson() error {
	dataJson := path.Join(r.genDataPath(), fmt.Sprintf(_DATA_NAME_FORMATTER, r.BeginTime.Format(timedef.TIME_FORMAT_IN_FILE)))
	bytes, err := json.MarshalIndent(r.Items, "", "    ")
	if err != nil {
		return err
	}
	if err := fileutil.WriteFile(dataJson, bytes); err != nil {
		return err
	}
	return nil
}

func (r *YHCReport) genPackageTarPath() string {
	return path.Join(r.CheckBase.Output, r.genPackageTarName())
}

func (r *YHCReport) genPackageName() string {
	return fmt.Sprintf(_PACKAGE_NAME_FORMATTER, r.BeginTime.Format(timedef.TIME_FORMAT_IN_FILE))
}

func (r *YHCReport) genPackageDir() string {
	return path.Join(r.CheckBase.Output, r.genPackageName())
}

func (r *YHCReport) genPackageTarName() string {
	return fmt.Sprintf("%s.tar.gz", r.genPackageName())
}

func (r *YHCReport) genDataPath() string {
	return path.Join(r.genPackageDir(), "data")
}

func (r *YHCReport) mkdir() error {
	if !fs.IsDirExist(r.CheckBase.Output) {
		if err := fs.Mkdir(r.CheckBase.Output); err != nil {
			return err
		}
		if err := yhccommons.ChownToExecuter(r.CheckBase.Output); err != nil {
			log.Module.Warnf("chown %s failed: %s", r.CheckBase.Output, err)
		}
	}
	if err := fs.Mkdir(r.genDataPath()); err != nil {
		return err
	}
	return nil
}

func (r *YHCReport) tarResult() error {
	command := fmt.Sprintf("cd %s;%s czvf %s %s;rm -rf %s", r.CheckBase.Output, bashdef.CMD_TAR, r.genPackageTarName(), r.genPackageName(), r.genPackageName())
	executer := execerutil.NewExecer(log.Logger)
	ret, _, stderr := executer.Exec(bashdef.CMD_BASH, "-c", command)
	if ret != 0 {
		return errors.New(stderr)
	}
	return nil
}

func (r *YHCReport) chownResult() error {
	return yhccommons.ChownToExecuter(r.genPackageTarPath())
}
