package processutil

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"yhc/commons/constants"
	"yhc/defs/timedef"
	"yhc/log"

	"github.com/shirou/gopsutil/process"
)

const (
	_lineSep = "\n"
	_itemSep = "\t"
	_kvSep   = ":"

	YasdbProcessReg = ".*yasdb (?i:(nomount|mount|open)) -D %s"
)

type Process struct {
	Pid             int     `json:"pid"`
	Name            string  `json:"name"`
	FullCommand     string  `json:"fullCommand"`
	Cmdline         string  `json:"cmdline"`
	ReadableCmdline string  `json:"readableCmdline"`
	User            string  `json:"user"`
	CPUPercent      float64 `json:"cpuPercent"`
	MemoryPercent   float64 `json:"memoryPercent"`
	CreateTime      string  `json:"createTime"`
	Status          string  `json:"status"`
}

func NewProcess(pid int) *Process {
	return &Process{
		Pid: pid,
	}
}

func (p *Process) IsRunning() (*os.Process, bool) {
	process, err := os.FindProcess(p.Pid)
	if err != nil {
		return nil, false
	}
	err = process.Signal(syscall.Signal(0))
	if err == nil || err == syscall.EPERM {
		return process, true
	}
	return nil, false
}

func (p *Process) getCmdline() ([]byte, error) {
	_, running := p.IsRunning()
	if !running {
		return nil, fmt.Errorf("process of pid(%d) not found", p.Pid)
	}
	filename := fmt.Sprintf("/proc/%d/cmdline", p.Pid)
	return os.ReadFile(filename)
}

func (p *Process) getCPUPercent() (float64, error) {
	_, running := p.IsRunning()
	if !running {
		return 0, fmt.Errorf("process of pid(%d) not found", p.Pid)
	}
	proc, err := process.NewProcess(int32(p.Pid))
	if err != nil {
		return 0, err
	}
	return proc.CPUPercent()
}

func (p *Process) getMemoryPercent() (float64, error) {
	_, running := p.IsRunning()
	if !running {
		return 0, fmt.Errorf("process of pid(%d) not found", p.Pid)
	}
	proc, err := process.NewProcess(int32(p.Pid))
	if err != nil {
		return 0, err
	}
	memPercent, err := proc.MemoryPercent()
	return float64(memPercent), err
}

func (p *Process) getCreateTime() (int64, error) {
	_, running := p.IsRunning()
	if !running {
		return 0, fmt.Errorf("process of pid(%d) not found", p.Pid)
	}
	proc, err := process.NewProcess(int32(p.Pid))
	if err != nil {
		return 0, err
	}
	return proc.CreateTime()
}

func (p *Process) getStatus() (string, error) {
	_, running := p.IsRunning()
	if !running {
		return "0", fmt.Errorf("process of pid(%d) not found", p.Pid)
	}
	proc, err := process.NewProcess(int32(p.Pid))
	if err != nil {
		return "", err
	}
	return proc.Status()
}

func (p *Process) GetCmdline() (string, error) {
	log := log.Util.M("process")
	data, err := p.getCmdline()
	if err != nil {
		return "", err
	}
	var bytes []byte
	for _, b := range data {
		// remove NUT(null)
		if b == 0 {
			log.Infof("get cmdline remove null byte")
			continue
		}
		bytes = append(bytes, b)
	}
	return string(bytes), nil
}

func (p *Process) GetCmdlineItems() ([]string, error) {
	data, err := p.getCmdline()
	if err != nil {
		return nil, err
	}
	var items []string
	lastIndex := -1
	for i, b := range data {
		if b == 0 {
			items = append(items, string(data[lastIndex+1:i]))
			lastIndex = i
			continue
		}
	}
	return items, nil
}

func (p *Process) parseInt(key, value string) (int, error) {
	key = strings.TrimSpace(key)
	if len(value) == 0 {
		return 0, nil
	}
	num, err := strconv.ParseInt(value, constants.BASE_DECIMAL, constants.BIT_SIZE_64)
	if err != nil {
		return 0, fmt.Errorf("%s of %s is not a number", value, key)
	}
	return int(num), nil
}

func (p *Process) GetUids() ([]int, error) {
	filename := fmt.Sprintf("/proc/%d/status", p.Pid)
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	content := string(data)
	lines := strings.Split(content, _lineSep)
	var uids []int
	for _, line := range lines {
		items := strings.Split(line, _itemSep)
		if len(items) < 2 {
			continue
		}
		key := strings.TrimSuffix(items[0], _kvSep)
		key = strings.ToLower(key)
		switch key {
		case "uid":
			for _, i := range items[1:] {
				num, err := p.parseInt(key, i)
				if err != nil {
					uids = append(uids, 0)
					continue
				}
				uids = append(uids, num)
			}
		}
	}
	if len(uids) == 0 {
		return nil, errors.New("get empty uid")
	}
	return uids, nil
}

func GetUsernameById(id int) (string, error) {
	u, err := user.LookupId(strconv.FormatInt(int64(id), constants.BASE_DECIMAL))
	if err != nil {
		return "", err
	}
	return u.Username, nil
}

func (p *Process) FindBaseInfo() error {
	_, running := p.IsRunning()
	if !running {
		return fmt.Errorf("process of pid(%d) not running", p.Pid)
	}
	items, err := p.GetCmdlineItems()
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return errors.New("empty cmdline")
	}
	uids, err := p.GetUids()
	if err != nil {
		return err
	}
	user, err := GetUsernameById(uids[0])
	if err != nil {
		return err
	}
	cpuPercent, err := p.getCPUPercent()
	if err != nil {
		return err
	}
	memoryPercent, err := p.getMemoryPercent()
	if err != nil {
		return err
	}
	createTimestamp, err := p.getCreateTime()
	if err != nil {
		return err
	}
	createTime := time.Unix(createTimestamp/1000, 0).Format(timedef.TIME_FORMAT)
	status, err := p.getStatus()
	if err != nil {
		return err
	}
	p.FullCommand = items[0]
	p.Cmdline = strings.Join(items, "")
	p.ReadableCmdline = strings.Join(items, " ")
	p.Name = path.Base(p.FullCommand)
	p.User = user
	p.CPUPercent = cpuPercent
	p.MemoryPercent = memoryPercent
	p.CreateTime = createTime
	p.Status = status
	return nil
}

func ListPids() ([]int, error) {
	dir, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	dirnames, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	var pids []int
	for _, fname := range dirnames {
		pid, err := strconv.ParseInt(fname, constants.BASE_DECIMAL, constants.BIT_SIZE_32)
		if err != nil {
			continue
		}
		pids = append(pids, int(pid))
	}
	sort.Slice(pids, func(i, j int) bool { return pids[i] < pids[j] })
	return pids, nil
}

// ListProcess lists all running processes with base information.
func ListProcess() ([]Process, error) {
	pids, err := ListPids()
	if err != nil {
		return nil, err
	}
	var processes []Process
	for _, pid := range pids {
		p := NewProcess(pid)
		if err := p.FindBaseInfo(); err != nil {
			continue
		}
		processes = append(processes, *p)
	}
	return processes, nil
}

// ListProcessByCmdline list all running processes by cmdline with base information.
func ListProcessByCmdline(user, cmdline string, isRegexp bool) ([]Process, error) {
	log := log.Util.M("process")
	// remove spaces of cmdline
	cmdline = strings.Join(strings.Fields(strings.TrimSpace(cmdline)), "")
	if len(cmdline) == 0 {
		log.Warnf("cmdline is empty, skip list process")
		return nil, nil
	}
	list, err := ListProcess()
	if err != nil {
		log.Warnf("list process err: %s", err.Error())
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	// cmdline and user are both equal
	matchFunc := func(p Process, user, cmdline string) bool {
		if isRegexp {
			match, _ := regexp.MatchString(cmdline, p.Cmdline)
			if match {
				if len(user) == 0 {
					return true
				}
				return p.User == user
			}
			return false
		}
		if p.Cmdline == cmdline {
			if len(user) == 0 {
				return true
			}
			return p.User == user
		}
		return false
	}
	var procs []Process
	for _, p := range list {
		if matchFunc(p, user, cmdline) {
			procs = append(procs, p)
		}
	}
	return procs, nil
}

func ListAnyUserProcessByCmdline(cmdline string, isRegexp bool) ([]Process, error) {
	return ListProcessByCmdline("", cmdline, isRegexp)
}

func GetYasdbProcess(dataPath string) ([]Process, error) {
	return ListAnyUserProcessByCmdline(fmt.Sprintf(YasdbProcessReg, dataPath), true)
}
