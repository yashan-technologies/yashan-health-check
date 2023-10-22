package runtimedef

import (
	"os"
	"path"
	"path/filepath"

	"yhc/utils/stringutil"
)

const (
	_ENV_YHC_HOME       = "YHC_HOME"
	_ENV_YHC_DEBUG_MODE = "YHC_DEBUG_MODE"
)

const (
	_DIR_NAME_LOG     = "log"
	_DIR_NAME_STATIC  = "static"
	_DIR_NAME_SCRIPTS = "scripts"
)

var _yhcHome string

func GetYHCHome() string {
	return _yhcHome
}

func GetLogPath() string {
	return path.Join(_yhcHome, _DIR_NAME_LOG)
}

func GetStaticPath() string {
	return path.Join(_yhcHome, _DIR_NAME_STATIC)
}

func GetScriptsPath() string {
	return path.Join(_yhcHome, _DIR_NAME_SCRIPTS)
}

func setYHCHome(v string) {
	_yhcHome = v
}

func isDebugMode() bool {
	return !stringutil.IsEmpty(os.Getenv(_ENV_YHC_DEBUG_MODE))
}

func getYHCHomeEnv() string {
	return os.Getenv(_ENV_YHC_HOME)
}

// genYHCHomeFromEnv generates ${YHC_HOME} from env, using YHC_HOME env as YHCHome in debug mode.
func genYHCHomeFromEnv() (yhcHome string, err error) {
	yhcHomeEnv := getYHCHomeEnv()
	if isDebugMode() && !stringutil.IsEmpty(yhcHomeEnv) {
		yhcHomeEnv, err = filepath.Abs(yhcHomeEnv)
		if err != nil {
			return
		}
		yhcHome = yhcHomeEnv
		return
	}
	return
}

// genYHCHomeFromRelativePath generates ${YHC_HOME} from relative path to the executable bin.
// executable bin locates at ${YHC_HOME}/bin/${executable}
func genYHCHomeFromRelativePath() (yhcHome string, err error) {
	executeable, err := getExecutable()
	if err != nil {
		return
	}
	yhcHome, err = filepath.Abs(path.Dir(path.Dir(executeable)))
	return
}

func initYHCHome() (err error) {
	yhcHome, err := genYHCHomeFromEnv()
	if err != nil {
		return
	}
	if !stringutil.IsEmpty(yhcHome) {
		setYHCHome(yhcHome)
		return
	}
	yhcHome, err = genYHCHomeFromRelativePath()
	if err != nil {
		return
	}
	setYHCHome(yhcHome)
	return
}
