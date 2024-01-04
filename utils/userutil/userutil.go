// The userutil package encapsulates functions related to users and user groups.
package userutil

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	"yhc/commons/constants"
	"yhc/defs/bashdef"
	"yhc/utils/execerutil"
	"yhc/utils/stringutil"

	"git.yasdb.com/go/yaslog"
)

const (
	NotRunSudo         = "not run sudo"
	PasswordIsRequired = "password is required"
)

const (
	ENV_SUDO_USER = "SUDO_USER"
	ROOT_USER_UID = 0
	ETC_GROUP     = "/etc/group"
)

var (
	ErrSudoNeedPwd = errors.New("a password is required")
	ErrNotRunSudo  = errors.New("user may not run sudo")
)

var (
	CurrentUser string
)

func init() {
	user, err := GetCurrentUser()
	if err != nil {
		panic(err)
	}
	CurrentUser = user
}

// GetUsernameById returns username by user ID.
func GetUsernameById(id int) (username string, err error) {
	u, err := user.LookupId(strconv.FormatInt(int64(id), constants.BASE_DECIMAL))
	if err != nil {
		return
	}
	username = u.Username
	return
}

// GetCurrentUser returns the current username.
func GetCurrentUser() (string, error) {
	return GetUsernameById(os.Getuid())
}

// GetUserGroups return groups of the user
func GetUserGroups(u *user.User) []string {
	groupids, err := u.GroupIds()
	if err != nil {
		return nil
	}
	groups := []string{}
	for _, gid := range groupids {
		g, err := user.LookupGroupId(gid)
		if err != nil {
			return nil
		}
		groups = append(groups, g.Name)
	}
	return groups
}

// IsCurrentUserRoot checks whether the current user is root.
func IsCurrentUserRoot() bool {
	return os.Getuid() == ROOT_USER_UID
}

// IsSysUserExists checks if the OS user exists.
func IsSysUserExists(username string) bool {
	_, err := user.Lookup(username)
	return err == nil
}

// IsSysGroupExists checks if the OS user group exists.
func IsSysGroupExists(groupname string) bool {
	_, err := user.LookupGroup(groupname)
	return err == nil
}

func GetRealUser() (*user.User, error) {
	if IsCurrentUserRoot() {
		username := os.Getenv(ENV_SUDO_USER)
		if len(username) == 0 {
			return user.LookupId(fmt.Sprint(ROOT_USER_UID))
		}
		return user.Lookup(username)
	}
	return user.LookupId(fmt.Sprint(os.Getuid()))
}

func GetUserOfGroup(log yaslog.YasLog, groupName string) ([]string, error) {
	execer := execerutil.NewExecer(log)
	cmd := fmt.Sprintf("%s %s", bashdef.CMD_CAT, ETC_GROUP)
	ret, stdout, stderr := execer.Exec(bashdef.CMD_BASH, "-c", cmd)
	if ret != 0 {
		return nil, errors.New(stderr)
	}
	var users []string
	groups := strings.Split(strings.TrimSpace(stdout), stringutil.STR_NEWLINE)
	for _, group := range groups {
		arr := strings.Split(group, stringutil.STR_COLON)
		// just like:'YASDBA:x:1021:yashan,mongodb,oracle,db,ycm,ny,golang' or 'db:x:1026:'
		if arr[0] != groupName || len(arr) < 4 || len(arr[3]) <= 0 {
			continue
		}
		users = append(users, strings.Split(arr[3], stringutil.STR_COMMA)...)
	}
	return users, nil
}
