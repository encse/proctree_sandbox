package streamstack

import (
	"bytes"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestStreamStack(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	s := NewReaderStack(buf)

	mx := sync.Mutex{}
	children := make([]*ChildReader, 0)
	wg := sync.WaitGroup{}
	for i := 0; i < 500; i++ {
		r := rand.Int() % 100
		if r < 10 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mx.Lock()
				childStream := s.AddChild()
				children = append(children, childStream)
				mx.Unlock()
				time.Sleep(2 * time.Second)
				s.CloseWithChildren(childStream)
			}()
		} else if r < 20 {
			mx.Lock()
			if len(children) > 0 {
				child := children[rand.Int()%len(children)]
				s.CloseWithChildren(child)
			}
			mx.Unlock()
		} else if r < 60 {
			mx.Lock()
			if len(children) > 0 {
				child := children[rand.Int()%len(children)]
				child.Read(make([]byte, 5))
			}
			mx.Unlock()
		} else {
			buf.Write([]byte("test"))
		}
	}

	wg.Wait()

}
