package runtimedef

import "yhc/utils/osutil"

var (
	_osRelease osutil.OSRelease
)

func GetOSRelease() osutil.OSRelease {
	return _osRelease
}

func initOSRelease() (err error) {
	osRelease, err := osutil.GetOSRelease()
	if err != nil {
		return
	}
	setOSRelease(*osRelease)
	return
}

func setOSRelease(osRelease osutil.OSRelease) {
	_osRelease = osRelease
}
