package streamstack

import (
	"io"
	"sync"
)

type ChildWriter struct {
	closed bool
	stack  *WriterStack
}

type WriterStack struct {
	wr       io.Writer
	children []*ChildWriter
	cond     *sync.Cond
}

func NewWriterStack(wr io.Writer) *WriterStack {
	locker := &sync.Mutex{}
	return &WriterStack{
		wr:       wr,
		children: make([]*ChildWriter, 0),
		cond:     sync.NewCond(locker),
	}
}

func (s *WriterStack) AddChild() *ChildWriter {
	s.cond.L.Lock()

	res := &ChildWriter{closed: false, stack: s}
	s.children = append(s.children, res)

	s.cond.Broadcast()
	s.cond.L.Unlock()
	return res
}

func (s *WriterStack) CloseWithChildren(firstChild *ChildWriter) {
	s.cond.L.Lock()
	idx := -1
	for i, child := range s.children {
		if firstChild == child {
			idx = i
		}

		if idx != -1 {
			child.closed = true
		}
	}
	if idx != -1 {
		s.children = s.children[:idx]
		s.cond.Broadcast()
	}
	s.cond.L.Unlock()
}

func (c *ChildWriter) Close() error {
	c.stack.CloseWithChildren(c)
	return nil
}

func (c *ChildWriter) Write(p []byte) (n int, err error) {
	c.stack.cond.L.Lock()
	defer c.stack.cond.L.Unlock()
	if c.closed {
		return 0, io.ErrClosedPipe
	}
	return c.stack.wr.Write(p)
}
