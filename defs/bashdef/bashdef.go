// The bashdef package defines the shell commands that will be used in the project.
// The shell commands used in the project must be defined here.
package bashdef

import "fmt"

const (
	CMD_BASH          = "bash"
	CMD_SUDO          = "/usr/bin/sudo"
	CMD_CAT_BY_REGEXP = "cat %s | grep -v grep | grep -E '%s'"
	CMD_TAR           = "tar"
	CMD_YASDB         = "yasdb"
	CMD_YASQL         = "yasql"
	CMD_CAT           = "cat"
	CMD_SAR           = "sar"
	CMD_SYSTEMCTL     = "systemctl"
	CMD_UFW           = "ufw"
	CMD_DMESG         = "dmesg"
	CMD_COMMAND       = "command"
	CMD_CP            = "cp"
)

const (
	COLOR_GREEN  = "\033[32m"
	COLOR_RED    = "\033[31m"
	COLOR_YELLOW = "\033[33m"
	COLOR_BLUE   = "\033[34m"
	COLOR_RESET  = "\033[0m"
)

func WithColor(s string, color string) string {
	return fmt.Sprintf("%s%s%s", color, s, COLOR_RESET)
}
