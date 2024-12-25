package domain

import (
	"sync"
)

type SeqID struct {
	id int
	mu sync.Mutex
}

func NewSeqID() *SeqID {
	return &SeqID{
		id: 1,
	}
}

func (s *SeqID) Get() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	old := s.id
	s.id++

	return old
}

type Patient struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
