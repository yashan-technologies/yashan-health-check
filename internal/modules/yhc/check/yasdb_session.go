package check

import (
	"fmt"
	"strconv"

	"yhc/commons/constants"
	"yhc/defs/confdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/yasdbutil"
)

const (
	SQL_QUERY_SESSION = `select type from v$session`

	PARAMETER_MAX_SESSIONS = "MAX_SESSIONS"

	KEY_USER_SESSIONS       = "USER_SESSIONS"
	KYE_BACKGROUND_SESSIONS = "BACKGROUND_SESSIONS"
	KEY_MAX_SESSIONS        = "MAX_SESSIONS"
	KEY_SESSION_USAGE       = "SESSION_USAGE"
)

func (c *YHCChecker) GetYasdbSession() (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_YASDB_SESSION,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_YASDB_SESSION))
	yasdb := yasdbutil.NewYashanDB(log, c.Yasdb)
	var userSessions, backgroundSessions int
	sessions, err := yasdb.QueryMultiRows(SQL_QUERY_SESSION, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		log.Errorf("failed to get data with sql %s, err: %v", SQL_QUERY_SESSION, err)
		data.Error = err.Error()
		return
	}
	for _, session := range sessions {
		if session["TYPE"] == "BACKGROUND" {
			backgroundSessions++
			continue
		}
		userSessions++
	}
	maxSessionsStr, err := c.querySingleParameter(log, PARAMETER_MAX_SESSIONS)
	if err != nil {
		data.Error = err.Error()
		log.Error("query %s failed, err: %v", PARAMETER_MAX_SESSIONS, err)
		return
	}
	maxSessions, err := strconv.ParseInt(maxSessionsStr, constants.BASE_DECIMAL, constants.BIT_SIZE_64)
	if err != nil {
		data.Error = err.Error()
		log.Error("faid to parse string %s to int, err: %v", maxSessionsStr, err)
		return
	}
	sessionUsageStr := fmt.Sprintf("%.2f", float64(userSessions+backgroundSessions)/float64(maxSessions)*100)
	sessionUsage, err := strconv.ParseFloat(sessionUsageStr, 64)
	if err != nil {
		data.Error = err.Error()
		log.Error("faid to parse string %s to float64, err: %v", sessionUsageStr, err)
		return
	}
	res := map[string]interface{}{
		KEY_USER_SESSIONS:       userSessions,
		KYE_BACKGROUND_SESSIONS: backgroundSessions,
		KEY_MAX_SESSIONS:        maxSessions,
		KEY_SESSION_USAGE:       sessionUsage,
	}
	data.Details = res
	return
}
