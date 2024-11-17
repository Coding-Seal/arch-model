package dequeue

type Dequeue[T any] struct {
	front, rear, capacity, size int
	data                        []T
}

func New[T any](n int) *Dequeue[T] {
	return &Dequeue[T]{
		data:     make([]T, n),
		front:    0,
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

	if !q.Empty() {
		q.rear = (q.rear + 1) % q.capacity
	}

	q.data[q.rear] = elem
	q.size++

	return true
}

func (q *Dequeue[T]) PushFront(elem T) bool {
	if q.Full() {
		return false
	}

	if !q.Empty() {
		if q.front == 0 {
			q.front = q.capacity - 1
		} else {
			q.front--
		}
	}

	q.data[q.front] = elem
	q.size++

	return true
}

func (q *Dequeue[T]) PopBack() bool {
	if q.Empty() {
		return false
	}

	q.rear = (q.rear - 1) % q.capacity
	q.size--

	return true
}

func (q *Dequeue[T]) PopFront() bool {
	if q.Empty() {
		return false
	}

	q.front = (q.front + 1) % q.capacity
	q.size--

	return true
}
