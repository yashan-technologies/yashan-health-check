package log

import (
	"fmt"
	"os"
	"path"
	"strings"

	"yhc/defs/compiledef"
	"yhc/defs/runtimedef"
	"yhc/utils/fileutil"
	"yhc/utils/userutil"

	"git.yasdb.com/go/yaslog"
	"git.yasdb.com/go/yasutil/fs"
)

const (
	LEVEL_DEBUG   = "DEBUG"
	LEVEL_NOTICE  = "NOTICE"
	LEVEL_INFO    = "INFO"
	LEVEL_WARNING = "WARNING"
	LEVEL_ERROR   = "ERROR"
	LEVEL_FATAL   = "FATAL"
)

const (
	_LOG_MODULE_CONTROLLER = "controller"
	_LOG_MODULE_HANDLER    = "handler"
	_LOG_MODULE_MODULE     = "module"
	_LOG_MODULE_YASQL      = "yasql"
	_LOG_MODULE_EXECER     = "execer"
)

// default option values
const (
	DEFAULT_LOG_LEVEL        = LEVEL_INFO
	DEFAULT_MAX_SIZE_BYTES   = 10 * (1 << 20) // 10 MB
	DEFAULT_INTERVAL_SECONDS = 30
)

// NODE: supported loggers
var (
	Logger     yaslog.YasLog
	Controller yaslog.YasLog
	Handler    yaslog.YasLog
	Module     yaslog.YasLog
	Yasql      yaslog.YasLog
	// TODO: more loggers
)

var _levelMap = map[string]yaslog.LogLevel{
	LEVEL_DEBUG:   yaslog.DEBUG,
	LEVEL_NOTICE:  yaslog.NOTICE,
	LEVEL_INFO:    yaslog.INFO,
	LEVEL_WARNING: yaslog.WARNING,
	LEVEL_ERROR:   yaslog.ERROR,
	LEVEL_FATAL:   yaslog.FATAL,
}

type Option struct {
	dir             string
	level           string
	maxSizeBytes    int64 // when the current log file size reaches <maxSizeBytes>, switch to another log file
	intervalSeconds int   // check the current log file size every <intervalSeconds>
	console         bool  // whether to print logs in the terminal
}

func (opt *Option) set(funcs []OptFunc) {
	for _, f := range funcs {
		f(opt)
	}
}

type OptFunc func(*Option)

func SetLevel(level string) OptFunc {
	return func(opt *Option) { opt.level = level }
}

func SetMaxSize(sizeBytes int64) OptFunc {
	return func(opt *Option) { opt.maxSizeBytes = sizeBytes }
}

func SetInterval(intervalSeconds int) OptFunc {
	return func(opt *Option) { opt.intervalSeconds = intervalSeconds }
}

func SetLogPath(logPath string) OptFunc {
	return func(opt *Option) { opt.dir = logPath }
}

func SetConsole(v bool) OptFunc {
	return func(opt *Option) { opt.console = v }
}

func DefaultLogOption() *Option {
	return &Option{
		level:           "DEBUG",
		dir:             runtimedef.GetLogPath(),
		maxSizeBytes:    DEFAULT_MAX_SIZE_BYTES,
		intervalSeconds: DEFAULT_INTERVAL_SECONDS,
	}
}

func NewLogOption(funcs ...OptFunc) *Option {
	opt := DefaultLogOption()
	opt.set(funcs)
	return opt
}

func InitLogger(servername string, opt *Option) error {
	fname := path.Join(opt.dir, fmt.Sprintf("%s.log", servername))
	if err := prepare(opt.dir, fname); err != nil {
		return nil
	}
	logger := yaslog.NewRotateLogger(fname, servername, opt.toYaslogOptions())
	Logger = logger
	Controller = logger.NewModule(_LOG_MODULE_CONTROLLER)
	Handler = logger.NewModule(_LOG_MODULE_HANDLER)
	Module = logger.NewModule(_LOG_MODULE_MODULE)
	Yasql = logger.NewModule(_LOG_MODULE_YASQL)
	logger.Infof("VERSION: %s", compiledef.GetAPPVersion())
	return nil
}

// create log path and log file
func prepare(logPath, fname string) error {
	owner := runtimedef.GetExecuteableOwner()
	if !fs.IsDirExist(logPath) {
		if err := fs.Mkdir(logPath); err != nil {
			return err
		}
	}
	if !fs.IsFileExist(fname) {
		if err := fileutil.WriteFile(fname, nil); err != nil {
			return err
		}
	}
	if userutil.IsCurrentUserRoot() && (owner.Uid != 0 || owner.Gid != 0) {
		_ = os.Chown(logPath, owner.Uid, owner.Gid)
		_ = os.Chown(fname, owner.Uid, owner.Gid)
	}
	return nil
}

func (opt *Option) toYaslogOptions() *yaslog.LogOption {
	var optFuncs []yaslog.LogOptFunc
	level, ok := _levelMap[strings.ToUpper(opt.level)]
	if !ok {
		level = yaslog.INFO
	}
	optFuncs = append(optFuncs, yaslog.SetFlag(yaslog.Ldate|yaslog.Ltime|yaslog.Lmidfile))
	optFuncs = append(optFuncs, yaslog.SetLevel(level))
	optFuncs = append(optFuncs, yaslog.SetFileSize(opt.maxSizeBytes))
	optFuncs = append(optFuncs, yaslog.SetInterval(opt.intervalSeconds))
	optFuncs = append(optFuncs, yaslog.SetConsole(opt.console))
	return yaslog.NewLogOption(optFuncs...)
}
