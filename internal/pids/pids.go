package pids

type Pid int

type PidTable struct {
	pids []Pid
}

func NewPidTable() *PidTable {
	return &PidTable{
		pids: make([]Pid, 0),
	}
}

func (pt *PidTable) Alloc() Pid {
	for pid := Pid(1); ; pid = Pid(pid + 1) {
		if !pt.Contains(pid) {
			pt.pids = append(pt.pids, pid)
			return pid
		}
	}
}

func (pt *PidTable) Free(pid Pid) bool {
	for i, pidT := range pt.pids {
		if pid == pidT {
			pt.pids = append(pt.pids[:i], pt.pids[i+1:]...)
			return true
		}
	}
	return false
}

func (pt *PidTable) Contains(pid Pid) bool {
	for _, pidT := range pt.pids {
		if pidT == pid {
			return true
		}
	}
	return false
}

func (pt *PidTable) List() []Pid {
	res := make([]Pid, len(pt.pids))
	copy(res, pt.pids)
	return res
}
