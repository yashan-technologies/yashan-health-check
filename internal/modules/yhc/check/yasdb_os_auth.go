package check

import (
	"path"
	"strings"

	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/fileutil"
	"yhc/utils/stringutil"
	"yhc/utils/userutil"

	ini "gopkg.in/ini.v1"
)

const (
	KEY_YASDB_OS_AUTH     = "ENABLE_LOCAL_OSAUTH"
	KEY_YASDBA_GROUP_USER = "YASDBA_GROUP_USER"

	VALUE_ON  = "on"
	VALUE_OFF = "off"

	GROUP_YASDBA       = "YASDBA"
	DIR_CONFIG         = "config"
	FILE_YASDB_NET_INI = "yasdb_net.ini"
)

func (c *YHCChecker) GetYasdbOSAuth() (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_YASDB_OS_AUTH,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_YASDB_OS_AUTH))
	yasdbNetIniPath := path.Join(c.base.DBInfo.YasdbData, DIR_CONFIG, FILE_YASDB_NET_INI)
	res := map[string]any{}
	res[KEY_YASDB_OS_AUTH] = VALUE_ON
	if fileutil.IsExist(yasdbNetIniPath) {
		iniConf, e := ini.Load(yasdbNetIniPath)
		if e != nil {
			err = e
			data.Error = err.Error()
			log.Errorf("failed to load yasdb_net.ini, err: %v", err)
			return
		}
		for _, section := range iniConf.Sections() {
			key, e := section.GetKey(KEY_YASDB_OS_AUTH)
			if e != nil || key.String() != VALUE_OFF {
				break
			}
			if key.String() == VALUE_OFF {
				res[KEY_YASDB_OS_AUTH] = VALUE_OFF
				data.Details = res
				return
			}
		}
	}
	users, err := userutil.GetUserOfGroup(log, GROUP_YASDBA)
	if err != nil {
		data.Error = err.Error()
		log.Errorf("failed to get user of group %s, err: %v", GROUP_YASDBA, err)
		return
	}
	res[KEY_YASDBA_GROUP_USER] = strings.Join(users, stringutil.STR_COMMA)
	data.Details = res
	return
}
