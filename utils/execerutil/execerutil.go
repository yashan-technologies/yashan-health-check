package execerutil

import (
	"git.yasdb.com/go/yaslog"
	"git.yasdb.com/go/yasutil/execer"
)

type Execer struct {
	execer.Execer
}

func NewExecer(log yaslog.YasLog, opts ...execer.ExecerOpt) *Execer {
	opts = append(opts, execer.WithPrintResult()) // default print result in debug mode
	return &Execer{
		Execer: *execer.NewExecer(log, opts...),
	}
}

func (e *Execer) Exec(bin string, arg ...string) (int, string, string) {
	return e.Execer.Exec(bin, arg...)
}

func (e *Execer) EnvExec(env []string, bin string, arg ...string) (int, string, string) {
	return e.Execer.EnvExec(env, bin, arg...)
}

func (e *Execer) Daemonize(bin string, arg ...string) error {
	return e.Execer.Daemonize(bin, arg...)
}

func (e *Execer) NohupProcess(env []string, logPath string, args ...string) error {
	return e.Execer.NohupProcess(env, logPath, args...)
}
