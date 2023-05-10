package algorithms

import (
	"github.com/gammazero/deque"
	"github.com/stretchr/testify/assert"
	"testing"
)

func CheckConnectivity(t *testing.T, data []map[int]bool) {
	head := 0
	found := make(map[int]bool)
	queue := deque.New[int]()
	queue.PushBack(head)

	for queue.Len() > 0 {
		curr := queue.PopFront()
		if found[curr] {
			continue
		}
		found[curr] = true
		for i, v := range data[curr] {
			if !v {
				continue
			}
			if !found[i] {
				queue.PushBack(i)
			}
		}
	}
	for k := 0; k < len(data); k++ {
		v, ok := found[k]
		assert.True(t, ok && v)
	}
}

func CheckGraph(t *testing.T, data []map[int]bool) {
	for i := range data {
		for j := range data {
			if i == j {
				assert.Equal(t, false, data[i][j])
			}
			assert.Equal(t, data[i][j], data[j][i])
		}
	}
}
