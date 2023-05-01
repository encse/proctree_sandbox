package vproc

import "sync"

type Pid int

type Vproc struct {
	pid Pid
}

type VprocTable struct {
	vprocs []Vproc
	mx     sync.RWMutex
}

func NewProcTable() *VprocTable {
	return &VprocTable{
		vprocs: make([]Vproc, 0),
	}
}

func (pt *VprocTable) Start() Vproc {
	pt.mx.Lock()
	defer pt.mx.Unlock()

	for pid := Pid(1); ; pid = Pid(pid + 1) {
		if pt.getVprocByPidLocked(pid) == nil {
			vproc := Vproc{pid: pid}
			pt.vprocs = append(pt.vprocs, vproc)
			return vproc
		}
	}
}

func (pt *VprocTable) Remove(pid Pid) bool {
	pt.mx.Lock()
	defer pt.mx.Unlock()

	for i, vproc := range pt.vprocs {
		if vproc.pid == pid {
			pt.vprocs = append(pt.vprocs[:i], pt.vprocs[i+1:]...)
			return true
		}
	}
	return false
}

func (pt *VprocTable) GetVprocByPid(pid Pid) *Vproc {
	pt.mx.RLock()
	defer pt.mx.RUnlock()
	return pt.getVprocByPidLocked(pid)
}

func (pt *VprocTable) GetVprocList() []Vproc {
	pt.mx.RLock()
	defer pt.mx.RUnlock()
	res := make([]Vproc, len(pt.vprocs))
	copy(res, pt.vprocs)
	return res
}

func (pt *VprocTable) getVprocByPidLocked(pid Pid) *Vproc {

	for _, vproc := range pt.vprocs {
		if vproc.pid == pid {
			return &vproc
		}
	}
	return nil
}
