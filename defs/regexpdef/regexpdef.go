package regexpdef

import "regexp"

const (
	range_format               = `^[1-9][0-9]*[Mdhms]$`
	time_format                = `^[1-9]\d{3}-(0\d|1[0-2])-([012]\d|3[01])(-(0\d|1\d|2[0-3])(-([0-5]\d))?)?$`
	path_format                = `^(\/?(\.{1,2}|[^\s]+)\/?([^\s]+\/)*[^\s]*)$`
	space_format               = `\s+`
	yasdb_process_format       = `^[/\w._-]*yasdb[\0](?i)(open|nomount|mount)[\0]-D[\0][/\w._-]+[\0]`
	key_value_format           = `^([^=]+)=(.*)$`
	mutiple_line_format        = `\n\s*\n`
	lsblk_ignore_device_format = `^(ram|loop|fd|(h|s|v|xv)d[a-z]|nvme\d+n\d+p)\d+$`
	lsblk_output_format        = `([A-Z]+)=(?:\"(.*?)\")`
)

var (
	TimeRegexp              = regexp.MustCompile(time_format)
	RangeRegexp             = regexp.MustCompile(range_format)
	PathRegexp              = regexp.MustCompile(path_format)
	SpaceRegexp             = regexp.MustCompile(space_format)
	YasdbProcessRegexp      = regexp.MustCompile(yasdb_process_format)
	KeyValueRegexp          = regexp.MustCompile(key_value_format)
	MultiLineRegexp         = regexp.MustCompile(mutiple_line_format)
	LsblkIgnoreDeviceRegexp = regexp.MustCompile(lsblk_ignore_device_format)
	LsblkOutputRegexp       = regexp.MustCompile(lsblk_output_format)
)
