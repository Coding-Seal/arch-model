package eventmanager

import (
	"context"

	"github.com/Coding-Seal/arch-model/internal/domain"
)

type EventManager struct {
	receiver chan domain.Event

	subscribers map[domain.EventType][]chan<- domain.Event
}

func (m *EventManager) PublishAccess() chan<- domain.Event {
	return m.receiver
}

func (m *EventManager) Subscribe(t domain.EventType) <-chan domain.Event {
	ch := make(chan domain.Event)
	m.subscribers[t] = append(m.subscribers[t], ch)

	return ch
}

func (m *EventManager) notify(event domain.Event) {
	for _, sub := range m.subscribers[event.Type()] {
		sub <- event
	}
}

func (m *EventManager) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-m.receiver:
				m.notify(event)
			}
		}
	}()
}
