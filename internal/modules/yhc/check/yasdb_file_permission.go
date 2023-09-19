package check

import (
	"fmt"

	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/fileutil"
)

type FilePermissionResp struct {
	PermissionMap map[string]string `json:"permissionMap"`
	WarningMap    map[string]string `json:"warningMap"`
}

func (c *YHCChecker) GetYasdbFilePermission() (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_YASDB_FILE_PERMISSION,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_YASDB_FILE_PERMISSION))
	permissionMap, errs := fileutil.GetFilesAccess(c.Yasdb.YasdbData)
	if err != nil {
		data.Error = err.Error()
		log.Errorf("failed to check files permission in %s, err: %v", c.Yasdb.YasdbData, err)
		return
	}
	res := make(map[string]string)
	otherWriteMap := map[string]string{}
	for filepath, filemode := range permissionMap {
		if fileutil.CheckOtherWrite(filemode) {
			otherWriteMap[filepath] = filemode.String()
		}
		res[filepath] = filemode.String()
	}
	for filePath, err := range errs {
		res[filePath] = fmt.Sprintf("check %s err: %v", filePath, err)
	}
	data.Details = FilePermissionResp{
		PermissionMap: res,
		WarningMap:    otherWriteMap,
	}
	return
}
