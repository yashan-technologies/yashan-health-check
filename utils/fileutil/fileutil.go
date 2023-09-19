// The stringutil package encapsulates functions related to files.
package fileutil

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"yhc/defs/errdef"
	"yhc/defs/regexpdef"
	"yhc/utils/userutil"
)

const (
	ROOT_DIR = "/"
)

var (
	ErrSyscallStatNotSupported = errors.New("syscall stat not supported")
)

type Owner struct {
	Uid       int
	Gid       int
	Username  string
	GroupName string
}

// IsPathSymlink checks whether a path is a real path or a link, and returns the real path when it is a link.
func IsPathSymlink(path string) (isSymlink bool, realPath string, err error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return
	}
	isSymlink = fi.Mode()&os.ModeSymlink != 0
	if isSymlink {
		realPath, err = os.Readlink(path)
	}
	return
}

// GetRealPath returns the real path directly when path is a real path,
// and returns the real path pointed to by a link when path is a link.
func GetRealPath(path string) (realPath string, err error) {
	return filepath.EvalSymlinks(path)
}

// GetOwnerID gets the user ID and user group ID to which path belongs.
func GetOwnerID(path string) (uid uint32, gid uint32, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return
	}
	state, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		err = ErrSyscallStatNotSupported
		return
	}
	uid, gid = state.Uid, state.Gid
	return
}

// GetOwnerID gets the username and user group name to which path belongs.
func GetOwner(path string) (owner Owner, err error) {
	uid, gid, err := GetOwnerID(path)
	if err != nil {
		return
	}
	u, err := user.LookupId(fmt.Sprint(uid))
	if err != nil {
		return
	}
	g, err := user.LookupGroupId(fmt.Sprint(gid))
	if err != nil {
		return
	}
	owner = Owner{
		Uid:       int(uid),
		Gid:       int(gid),
		Username:  u.Username,
		GroupName: g.Name,
	}
	return
}

// judge file or path is exist
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

func ReadFile(file string) ([]byte, error) {
	_, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(file)
}

// 获取系统配置,仅限配置文件是key = value类型的
func GetConfByKey(configPath string, key string) (string, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	res := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = regexpdef.SpaceRegex.ReplaceAllString(line, "")
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, key) {
			res = strings.TrimPrefix(line, fmt.Sprintf("%s=", key))
			break
		}
	}
	return res, nil
}

func CheckAccess(p string) error {
	_, err := os.Stat(p)
	if err != nil {
		return err
	}
	file, err := os.Open(p)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func GetPidByPidFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	buffer := make([]byte, fileinfo.Size())
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}
	pidUint := binary.LittleEndian.Uint32(buffer)
	return strconv.FormatUint(uint64(pidUint), 10), nil
}

func IsAncestorDir(ancestorDir, dir string) bool {
	if !path.IsAbs(ancestorDir) || !path.IsAbs(dir) {
		return false
	}
	if ancestorDir == dir || ancestorDir == ROOT_DIR {
		return true
	}
	for i := dir; i != ROOT_DIR; i = path.Dir(i) {
		if i == ancestorDir {
			return true
		}
	}
	return false
}

func ComparePathDepth(path1, path2 string) int {
	sep := string(filepath.Separator)
	depth1 := strings.Count(path1, sep)
	depth2 := strings.Count(path2, sep)
	return depth1 - depth2

}

func CheckUserWrite(path string) error {
	file, err := os.Stat(path)
	if err != nil {
		return err
	}
	if file.Mode().Perm()&syscall.S_IWRITE == 0 {
		return errdef.NewErrPermissionDenied(userutil.CurrentUser, path)
	}
	return nil
}

func CheckUserRead(path string) error {
	file, err := os.Stat(path)
	if err != nil {
		return err
	}
	if file.Mode().Perm()&syscall.S_IREAD == 0 {
		return errdef.NewErrPermissionDenied(userutil.CurrentUser, path)
	}
	return nil
}

func CheckUserExec(path string) error {
	file, err := os.Stat(path)
	if err != nil {
		return err
	}
	if file.Mode().Perm()&syscall.S_IEXEC == 0 {
		return errdef.NewErrPermissionDenied(userutil.CurrentUser, path)
	}
	return nil
}

func CheckDirAccess(dir string, excludeMap map[string]struct{}) (res map[string]error, err error) {
	res = make(map[string]error)
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if _, ok := excludeMap[path]; ok {
			return nil
		}
		if err != nil {
			res[path] = err
			return nil
		}
		if err := CheckAccess(path); err != nil {
			res[path] = err
			return nil
		}
		return nil
	})
	return
}

func GetFilesAccess(dir string) (map[string]os.FileMode, map[string]error) {
	permissionsMap := make(map[string]os.FileMode)
	errs := make(map[string]error)
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errs[dir] = err
			return nil
		}
		if !info.IsDir() {
			permissionsMap[path] = info.Mode().Perm()
		}
		return nil
	})
	return permissionsMap, errs
}

func CheckOtherWrite(fileMode os.FileMode) bool {
	return (fileMode.Perm() & syscall.S_IWOTH) != 0
}
