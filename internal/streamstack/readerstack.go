package streamstack

import (
	"bufio"
	"io"
	"sync"
)

type ChildReader struct {
	closed bool
	stack  *ReaderStack
}

type ReaderStack struct {
	rd       bufio.Reader
	rdLocker sync.Locker // only one child is allowed to use rd at a time

	children []*ChildReader
	cond     *sync.Cond
}

func NewReaderStack(rd io.Reader) *ReaderStack {
	locker := &sync.Mutex{}
	return &ReaderStack{
		rd:       *bufio.NewReader(rd),
		children: make([]*ChildReader, 0),
		cond:     sync.NewCond(locker),
		rdLocker: &sync.Mutex{},
	}
}

func (s *ReaderStack) AddChild() *ChildReader {
	s.cond.L.Lock()

	res := &ChildReader{closed: false, stack: s}
	s.children = append(s.children, res)

	s.cond.Broadcast()
	s.cond.L.Unlock()
	return res
}

func (s *ReaderStack) CloseWithChildren(firstChild *ChildReader) {
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

func (c *ChildReader) Read(output []byte) (n int, err error) {
	return c.stack.Read(c, output)
}

func (c *ChildReader) Close() error {
	c.stack.CloseWithChildren(c)
	return nil
}

func (s *ReaderStack) Read(c *ChildReader, output []byte) (n int, err error) {
	s.cond.L.Lock()
	for {
		// we hold s.cond.L.Lock() in the body of this loop and
		// temporarily release it while waiting in the if statements

		if c.closed {
			s.cond.L.Unlock()
			return 0, io.EOF
		}

		if c != s.children[len(s.children)-1] {
			s.cond.Wait()
			continue
		}

		// s.cond.L is locked first, then s.rdMx
		s.rdLocker.Lock()

		if s.rd.Buffered() == 0 {
			// wait for more bytes / error in the input stream
			s.cond.L.Unlock()
			_, err = s.rd.ReadByte()
			if err == nil {
				s.rd.UnreadByte()
			}
			s.rdLocker.Unlock()
			s.cond.L.Lock()
			continue
		} else {
			lim := s.rd.Buffered()
			if lim > len(output) {
				lim = len(output)
			}
			n, err := s.rd.Read(output[n:lim])
			s.rdLocker.Unlock()
			s.cond.L.Unlock()
			return n, err
		}
	}
}
