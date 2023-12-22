package reporter

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"yhc/defs/bashdef"
	"yhc/defs/confdef"
	"yhc/defs/timedef"
	yhccommons "yhc/internal/modules/yhc/check/commons"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/execerutil"
	"yhc/utils/fileutil"

	"git.yasdb.com/go/yasutil/execer"
	"git.yasdb.com/go/yasutil/fs"
)

const (
	_PACKAGE_NAME_FORMATTER          = "yhc-%s"
	_DATA_NAME_FORMATTER             = "data-%s.json"
	_REPORT_JSON_NAME_FORMATTER      = "report-%s.json"
	_FAILED_ITEM_JSON_NAME_FORMATTER = "failed-%s.json"
	_REPORT_NAME_FORMATTER           = "report-%s.html"
	_WORD_REPORT_NAME_FORMATTER      = "report-%s.docx"

	_DIR_HTML_TEMPLATE  = "html-template"
	_FILE_HTML_TEMPLATE = "template.html"

	_SCRIPTS          = "scripts"
	_WORD_GENNER_PATH = "wordgenner"

	_TEMPLATE_KEY               = "$GLOBAL={}"
	_TEMPLATE_REPLACE_FORMATTER = "$GLOBAL=%s"
)

type YHCReport struct {
	YHCHome    string                                `json:"YHCHome"`
	BeginTime  time.Time                             `json:"beginTime"`
	EndTime    time.Time                             `json:"endTime"`
	CheckBase  *define.CheckerBase                   `json:"checkBase"`
	Items      map[define.MetricName]*define.YHCItem `json:"items"`
	Report     *define.PandoraReport
	FailedItem map[define.MetricName]string
}

func NewYHCReport(yhcHome string, checkBase *define.CheckerBase) *YHCReport {
	return &YHCReport{
		YHCHome:   yhcHome,
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
	if err := r.genReportJson(); err != nil {
		log.Errorf("gen data err: %s", err.Error())
		return "", err
	}
	if err := r.genFailedItemJson(); err != nil {
		log.Errorf("gen data err: %s", err.Error())
		return "", err
	}
	if err := r.genReport(); err != nil {
		log.Errorf("gen report failed: %s", err)
		return "", err
	}
	if err := r.tarResult(); err != nil {
		log.Errorf("tar result failed: %s", err)
		return "", err
	}
	if err := r.chownResult(); err != nil {
		log.Errorf("chown result failed: %s", err)
	}
	return r.genPackageTarPath(), nil
}

func (r *YHCReport) genReport() error {
	if err := r.genHtmlReport(); err != nil {
		return err
	}
	if err := r.genWordReport(); err != nil {
		return err
	}
	return nil
}

func (r *YHCReport) genHtmlReport() error {
	log := log.Module.M("gen-html")
	if confdef.GetYHCConf().SkipGenHtmlReport {
		log.Debug("skip to gen html report")
		return nil
	}
	templateFile := r.getHtmlTemplateFile()
	f, err := os.Open(templateFile)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	content, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(r.Report)
	if err != nil {
		return err
	}
	replacement := fmt.Sprintf(_TEMPLATE_REPLACE_FORMATTER, string(jsonData))
	contentStr := string(content)
	newContentStr := strings.Replace(contentStr, _TEMPLATE_KEY, replacement, 1)
	return fileutil.WriteFile(r.genReportFilePath(), []byte(newContentStr))
}

func (r *YHCReport) genWordReport() error {
	log := log.Module.M("gen-word")
	if confdef.GetYHCConf().SkipGenWordReport {
		log.Debug("skip to gen word report")
		return nil
	}
	wordGenner := r.getWordGennerFile()
	exec := execer.NewExecer(log)
	cmd := []string{
		wordGenner,
		"-i",
		r.getReportJsonFile(),
		"-o",
		r.getWordReportFile(),
	}
	ret, _, stderr := exec.Exec(bashdef.CMD_BASH, "-c", strings.Join(cmd, " "))
	if ret != 0 {
		err := fmt.Errorf("gen word report err: %s", stderr)
		log.Error(err)
		return err
	}
	return nil
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

func (r *YHCReport) genReportJson() error {
	dataJson := r.getReportJsonFile()
	bytes, err := json.MarshalIndent(r.Report, "", "    ")
	if err != nil {
		return err
	}
	if err := fileutil.WriteFile(dataJson, bytes); err != nil {
		return err
	}
	return nil
}

func (r *YHCReport) genFailedItemJson() error {
	dataJson := r.getFailedItemJsonFile()
	bytes, err := json.MarshalIndent(r.FailedItem, "", "    ")
	if err != nil {
		return err
	}
	if err := fileutil.WriteFile(dataJson, bytes); err != nil {
		return err
	}
	return nil
}

func (r *YHCReport) getReportJsonFile() string {
	return path.Join(r.genDataPath(), fmt.Sprintf(_REPORT_JSON_NAME_FORMATTER, r.BeginTime.Format(timedef.TIME_FORMAT_IN_FILE)))
}

func (r *YHCReport) getFailedItemJsonFile() string {
	return path.Join(r.genDataPath(), fmt.Sprintf(_FAILED_ITEM_JSON_NAME_FORMATTER, r.BeginTime.Format(timedef.TIME_FORMAT_IN_FILE)))
}

func (r *YHCReport) getWordReportFile() string {
	return path.Join(r.genPackageDir(), fmt.Sprintf(_WORD_REPORT_NAME_FORMATTER, r.BeginTime.Format(timedef.TIME_FORMAT_IN_FILE)))
}

func (r *YHCReport) genReportFilePath() string {
	return path.Join(r.genPackageDir(), fmt.Sprintf(_REPORT_NAME_FORMATTER, r.BeginTime.Format(timedef.TIME_FORMAT_IN_FILE)))
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

func (r *YHCReport) getHtmlTemplateFile() string {
	return path.Join(r.YHCHome, _DIR_HTML_TEMPLATE, _FILE_HTML_TEMPLATE)
}

func (r *YHCReport) getWordGennerFile() string {
	return path.Join(r.YHCHome, _SCRIPTS, _WORD_GENNER_PATH, "wordgenner")
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
