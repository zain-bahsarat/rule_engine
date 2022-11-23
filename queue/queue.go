package queue

import (
	"sync"
)

// queue struct definition
type queue struct {
	mtx   *sync.Mutex
	items []interface{}
}

func New() *queue {
	return &queue{
		mtx:   &sync.Mutex{},
		items: make([]interface{}, 0),
	}
}

// Enqueue adds the item into the queue
func (q *queue) Enqueue(item interface{}) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	q.items = append(q.items, item)
}

// Dequeue removes the item from the queue and RETURNS item
func (q *queue) Dequeue() interface{} {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	item := q.items[0]
	q.items = q.items[1:]
	return item
}

// Items returns the queue items
func (q *queue) Items() []interface{} {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	return q.items
}

// Length returns the queue length
func (q *queue) Length() int {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	return len(q.items)
}
