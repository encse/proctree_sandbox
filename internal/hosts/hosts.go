package hosts

import (
	"io"
	"sync"
	"time"

	"github.com/encse/proctree_sandbox/internal/pids"
	"github.com/encse/proctree_sandbox/internal/streamstack"
)

type Host struct {
	Name   string
	pids   *pids.PidTable
	vprocs map[pids.Pid]Vproc
	mx     sync.RWMutex
}

type Conn struct {
	stdinStack  *streamstack.ReaderStack
	stdoutStack *streamstack.WriterStack
}

func NewConn(stdin io.Reader, stdout io.Writer) Conn {
	return Conn{
		stdinStack:  streamstack.NewReaderStack(stdin),
		stdoutStack: streamstack.NewWriterStack(stdout),
	}
}

type Executable interface {
	Name() string
	Run(env Env) error
}

type Env struct {
	Host   *Host
	Pid    pids.Pid
	Args   []string
	Stdin  io.Reader
	Stdout io.Writer
	Conn   Conn
}

type Vproc struct {
	Pid        pids.Pid
	Executable Executable
	Stdin      io.ReadCloser
	Stdout     io.WriteCloser
	Started    time.Time
	Conn       Conn
}

func New(name string) *Host {
	return &Host{
		Name:   name,
		pids:   pids.NewPidTable(),
		vprocs: map[pids.Pid]Vproc{},
	}
}

func (env Env) Exec(host *Host, executable Executable) error {
	return Exec(env.Conn, host, executable)
}

func Exec(conn Conn, host *Host, executable Executable) error {
	host.mx.Lock()
	pid := host.pids.Alloc()

	stdin := conn.stdinStack.AddChild()
	stdout := conn.stdoutStack.AddChild()

	env := Env{
		Pid:    pid,
		Host:   host,
		Stdin:  stdin,
		Stdout: stdout,
		Conn:   conn,
	}

	vproc := Vproc{
		Pid:        pid,
		Executable: executable,
		Stdin:      stdin,
		Stdout:     stdout,
		Conn:       conn,
		Started:    time.Now(),
	}

	host.vprocs[pid] = vproc
	host.mx.Unlock()

	res := executable.Run(env)
	host.Terminate(vproc)

	return res
}

func (host *Host) Terminate(vproc Vproc) {
	host.mx.Lock()
	vproc.Stdin.Close()
	vproc.Stdout.Close()
	delete(host.vprocs, vproc.Pid)
	host.pids.Free(vproc.Pid)
	host.mx.Unlock()
}

func (host *Host) GetVprocs() []Vproc {
	host.mx.RLock()
	defer host.mx.RUnlock()

	res := make([]Vproc, 0, len(host.vprocs))
	for _, vproc := range host.vprocs {
		res = append(res, vproc)
	}
	return res
}
