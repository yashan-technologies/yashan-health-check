// The userutil package encapsulates functions related to users and user groups.
package userutil

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strconv"
)

const (
	NotRunSudo         = "not run sudo"
	PasswordIsRequired = "password is required"
)

const (
	ENV_SUDO_USER = "SUDO_USER"
	ROOT_USER_UID = 0
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
	u, err := user.LookupId(strconv.FormatInt(int64(id), 10))
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
