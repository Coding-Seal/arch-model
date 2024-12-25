package dequeue

type Dequeue[T any] struct {
	front, rear, capacity, size int
	data                        []T
}

func New[T any](n int) *Dequeue[T] {
	return &Dequeue[T]{
		data:     make([]T, n),
		front:    1,
		rear:     0,
		size:     0,
		capacity: n,
	}
}

func (q *Dequeue[T]) Full() bool {
	return q.size == q.capacity
}

func (q *Dequeue[T]) Empty() bool {
	return q.size == 0
}

func (q *Dequeue[T]) Back() (T, bool) {
	var res T
	if q.Empty() {
		return res, false
	}

	res = q.data[q.rear]

	return res, true
}

func (q *Dequeue[T]) Front() (T, bool) {
	var res T
	if q.Empty() {
		return res, false
	}

	res = q.data[q.front]

	return res, true
}

func (q *Dequeue[T]) PushBack(elem T) bool {
	if q.Full() {
		return false
	}

	q.rear = mod((q.rear + 1), q.capacity)
	q.data[q.rear] = elem
	q.size++

	return true
}

func (q *Dequeue[T]) PushFront(elem T) bool {
	if q.Full() {
		return false
	}

	q.front = mod((q.front - 1), q.capacity)
	q.data[q.front] = elem
	q.size++

	return true
}

func (q *Dequeue[T]) PopBack() bool {
	var null T

	if q.Empty() {
		return false
	}

	q.data[q.rear] = null
	q.rear = mod((q.rear - 1), q.capacity)
	q.size--

	return true
}

func (q *Dequeue[T]) PopFront() bool {
	var null T

	if q.Empty() {
		return false
	}

	q.data[q.front] = null
	q.front = mod((q.front + 1), q.capacity)
	q.size--

	return true
}

func mod(a, b int) int {
	return (a%b + b) % b
}
