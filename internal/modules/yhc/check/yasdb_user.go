package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbUserLoginPasswordStrength() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_SECURITY_LOGIN_PASSWORD_STRENGTH)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbUserLoginMaximumAttempts() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_SECURITY_LOGIN_MAXIMUM_LOGIN_ATTEMPTS)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbUserNoOpen() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_SECURITY_USER_NO_OPEN)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbUserWithSystemTablePrivileges() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_SECURITY_USER_WITH_SYSTEM_TABLE_PRIVILEGES)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbUserWithDBARole() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_SECURITY_USER_WITH_DBA_ROLE)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbUserAllPrivilegeOrSystemPrivileges() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_SECURITY_USER_ALL_PRIVILEGE_OR_SYSTEM_PRIVILEGES)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbUserUseSystemTablespace() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_SECURITY_USER_USE_SYSTEM_TABLESPACE)
	defer c.fillResult(data)
	return
}
