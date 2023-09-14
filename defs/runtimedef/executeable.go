package runtimedef

import (
	"os"

	"yhc/utils/fileutil"
)

var (
	_owner fileutil.Owner
)

func GetExecuteableOwner() fileutil.Owner {
	return _owner
}

func getExecutable() (executeable string, err error) {
	executeable, err = os.Executable()
	if err != nil {
		return
	}
	return fileutil.GetRealPath(executeable)
}

func setOwner(owner fileutil.Owner) {
	_owner = fileutil.Owner{
		Uid:       owner.Uid,
		Gid:       owner.Gid,
		Username:  owner.Username,
		GroupName: owner.GroupName,
	}
}

func initExecuteable() (err error) {
	executeable, err := getExecutable()
	if err != nil {
		return
	}
	var owner fileutil.Owner
	owner, e := fileutil.GetOwner(executeable)
	if e != nil {
		// if failed to get user name, just fill user id and group id
		userId, groupId, _ := fileutil.GetOwnerID(executeable)
		owner.Uid = int(userId)
		owner.Gid = int(groupId)
	}
	setOwner(owner)
	return
}
