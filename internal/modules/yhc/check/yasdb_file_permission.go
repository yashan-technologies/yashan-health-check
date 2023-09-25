package check

import (
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/fileutil"
)

const (
	KEY_FILE_PATH  = "filePath"
	KEY_PERMISSION = "permission"
)

func (c *YHCChecker) GetYasdbFilePermission() (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_YASDB_FILE_PERMISSION,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_YASDB_FILE_PERMISSION))
	permissionMap, errs := fileutil.GetFilesAccess(c.base.DBInfo.YasdbData)
	if err != nil {
		data.Error = err.Error()
		log.Errorf("failed to check files permission in %s, err: %v", c.base.DBInfo.YasdbData, err)
		return
	}

	res := []map[string]string{}
	for filePath, fileMode := range permissionMap {
		res = append(res, map[string]string{
			KEY_FILE_PATH:  filePath,
			KEY_PERMISSION: fileMode.String(),
		})
	}
	for filePath, err := range errs {
		res = append(res, map[string]string{
			KEY_FILE_PATH:  filePath,
			KEY_PERMISSION: err.Error(),
		})
	}
	data.Details = res
	return
}
