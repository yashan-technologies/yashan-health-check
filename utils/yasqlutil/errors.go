package yasqlutil

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrInvalidSelect        = NewYasqlError("cannot select nil")
	ErrNoColumns            = NewYasqlError("cannot select no columns")
	ErrTableNotSet          = NewYasqlError("table not set")
	ErrInvalidDest          = NewYasqlError("type of dest should be pointer")
	ErrConnect              = NewYasqlError("connect failed")
	ErrDBNotOpen            = NewYasqlError("the database is not open")
	ErrYasqlHomeNotExist    = NewYasqlError("path: yasql-home not exist")
	ErrInvalidUserPwd       = NewYasqlError("invalid username/password, login denied")
	ErrLackLoginAuth        = NewYasqlError("user lacks CREATE SESSION privilege; logon denied")
	ErrNoDBUser             = NewYasqlError("user does not exist")
	ErrDBUserLackPrivileges = NewYasqlError("user lacks privileges")
	ErrInvalidUserPwdValue  = NewYasqlError("invalid username/password")
	ErrRecordNotFound       = NewYasqlError("record not found")
)

var (
	EnterRegexp   = regexp.MustCompile(`\n`)
	DivLineRegexp = regexp.MustCompile(`_+`)
	YasRegexp     = regexp.MustCompile(`YAS-\d{5}`)
)

type YasErr struct {
	Prefix   string
	Msg      string
	yasqlMsg string
}

type yasqlError struct {
	msg       string
	originMsg string
}

func (e *yasqlError) Error() string {
	return e.msg
}

func NewYasqlError(message string) *yasqlError {
	return &yasqlError{msg: message}
}

func IsYasqlError(err error) bool {
	var e *yasqlError
	return errors.As(err, &e)
}

func (e *yasqlError) SetOriginErr(s string) {
	e.originMsg = s
}

func NewYasErr(stdout string) error {
	resolve := EnterRegexp.ReplaceAllString(stdout, "_")
	resolve = DivLineRegexp.ReplaceAllString(resolve, "_")
	fields := strings.Split(resolve, "_")
	e := &YasErr{}
	if len(fields) < 2 {
		e.Msg = stdout
		return e
	}
	var (
		yasErr   = strings.TrimSpace(fields[0])
		yasqlErr = strings.TrimSpace(fields[1])
	)
	match := YasRegexp.FindStringSubmatch(yasErr)
	if len(match) < 1 {
		e.Msg = stdout
		return e
	}
	yasPrefix := match[0]
	yasErr = strings.TrimSpace(strings.ReplaceAll(yasErr, yasPrefix, ""))
	e.Prefix = yasPrefix
	e.Msg = yasErr
	e.yasqlMsg = yasqlErr
	return e
}

func (y *YasErr) Error() string {
	if len(y.Prefix) == 0 {
		return y.Msg
	}
	return fmt.Sprintf("%s:%s", y.Prefix, y.Msg)
}
