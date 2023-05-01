package pids

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	// given
	pt := NewPidTable()

	pt.Alloc()
	pt.Alloc()
	pt.Alloc()

	// then
	assert.Equal(t, len(pt.List()), 3)

	// when
	pt.Free(Pid(2))

	// then
	assert.Equal(t, len(pt.List()), 2)
}

func TestContains(t *testing.T) {
	// given
	pt := NewPidTable()

	pids := []Pid{
		pt.Alloc(),
		pt.Alloc(),
		pt.Alloc(),
	}

	// then
	assert.True(t, pt.Contains(pids[0]))
	assert.True(t, pt.Contains(pids[1]))
	assert.True(t, pt.Contains(pids[2]))

	// and when
	pt.Free(pids[1])

	// then
	assert.True(t, pt.Contains(pids[0]))
	assert.False(t, pt.Contains(pids[1]))
	assert.True(t, pt.Contains(pids[2]))
}

func TestPidAssignment(t *testing.T) {
	// given
	pt := NewPidTable()

	pids := []Pid{
		pt.Alloc(),
		pt.Alloc(),
		pt.Alloc(),
	}

	// then pids should grow
	assert.Equal(t, pids[0], Pid(1))
	assert.Equal(t, pids[1], Pid(2))
	assert.Equal(t, pids[2], Pid(3))

	// and when
	pt.Free(Pid(2))
	pid := pt.Alloc()

	// then unused Pid(2) should be reassigned
	assert.Equal(t, pid, Pid(2))
}
