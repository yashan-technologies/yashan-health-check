package daemon

type DaemonCmd struct {
	Start   startCmd   `cmd:"start"   name:"start"   help:"Start yashan health check daemon."`
	Stop    stopCmd    `cmd:"stop"    name:"stop"    help:"Stop yashan health check daemon."`
	Restart restartCmd `cmd:"restart" name:"restart" help:"Restart yashan health check daemon."`
	Status  statusCmd  `cmd:"status"  name:"status"  help:"Show yashan health check daemon status."`
	Reload  reloadCmd  `cmd:"reload"  name:"reload"  help:"Reload yashan health check daemon."`
}

// [Interface Func]
func (c DaemonCmd) Run() error {
	return nil
}
