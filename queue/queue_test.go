package queue

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type action struct {
	Name  string
	Value interface{}
}

func TestQueue(t *testing.T) {

	testcases := []struct {
		actions  []action
		expected []interface{}
	}{
		{
			actions:  []action{{Name: "E", Value: 1}, {Name: "E", Value: 2}},
			expected: []interface{}{1, 2},
		},
		{
			actions:  []action{{Name: "E", Value: 1}, {Name: "E", Value: 2}, {Name: "D", Value: 2}},
			expected: []interface{}{2},
		},
		{
			actions:  []action{{Name: "E", Value: "OK"}, {Name: "E", Value: 2}, {Name: "D", Value: 2}},
			expected: []interface{}{2},
		},
	}

	for i, tt := range testcases {
		queue := New()
		for _, action := range tt.actions {
			if action.Name == "E" {
				queue.Enqueue(action.Value)
			} else {
				queue.Dequeue()
			}
		}

		assert.Equal(t, tt.expected, queue.Items(), fmt.Sprintf("tests[%d] - queue items not equal", i))
		assert.Equal(t, len(tt.expected), queue.Length())
	}
}
