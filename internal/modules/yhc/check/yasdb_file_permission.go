package check

import (
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/fileutil"
)

const (
	KEY_FILE_PATH  = "filePath"
	KEY_PERMISSION = "permission"
	KEY_OWNER      = "owner"
	KEY_GROUP      = "group"
)

func (c *YHCChecker) GetYasdbFilePermission(name string) (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_YASDB_FILE_PERMISSION,
	}
	defer c.fillResults(data)

	log := log.Module.M(string(define.METRIC_YASDB_FILE_PERMISSION))
	permissionMap, errs := fileutil.GetFilesAccess(c.base.DBInfo.YasdbData)
	if err != nil {
		data.Error = err.Error()
		log.Errorf("failed to check files permission in %s, err: %v", c.base.DBInfo.YasdbData, err)
		return
	}

	res := []map[string]string{}
	for filePath, err := range errs {
		res = append(res, map[string]string{
			KEY_FILE_PATH:  filePath,
			KEY_PERMISSION: err.Error(),
		})
	}
	for filePath, fileMode := range permissionMap {
		var keyOwner, keyGroup string
		owner, err := fileutil.GetOwner(filePath)
		if err != nil {
			log.Errorf("failed to get owner of %s, err: %v", filePath, err)
		} else {
			keyOwner = owner.Username
			keyGroup = owner.GroupName
		}
		res = append(res, map[string]string{
			KEY_FILE_PATH:  filePath,
			KEY_OWNER:      keyOwner,
			KEY_GROUP:      keyGroup,
			KEY_PERMISSION: fileMode.String(),
		})
	}
	data.Details = res
	return
}
