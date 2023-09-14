package regexpdef

import "regexp"

const (
	range_format         = `^[1-9][0-9]*[Mdhms]$`
	time_format          = `^[1-9]\d{3}-(0\d|1[0-2])-([012]\d|3[01])(-(0\d|1\d|2[0-3])(-([0-5]\d))?)?$`
	path_format          = `^(\/?(\.{1,2}|[^\s]+)\/?([^\s]+\/)*[^\s]*)$`
	space_format         = `\s+`
	yasdb_process_format = `^[/\w._-]*yasdb[\0](?i)(open|nomount|mount)[\0]-D[\0][/\w._-]+[\0]`
	key_value_format     = `^([^=]+)=(.*)$`
)

var (
	TimeRegex         = regexp.MustCompile(time_format)
	RangeRegex        = regexp.MustCompile(range_format)
	PathRegex         = regexp.MustCompile(path_format)
	SpaceRegex        = regexp.MustCompile(space_format)
	YasdbProcessRegex = regexp.MustCompile(yasdb_process_format)
	KeyValueRegex     = regexp.MustCompile(key_value_format)
)
