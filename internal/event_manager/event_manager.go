package eventmanager

import (
	"context"
	"log/slog"

	"github.com/Coding-Seal/arch-model/internal/domain"
	"github.com/Coding-Seal/arch-model/pkg/logger"
)

type EventManager struct {
	receiver    chan domain.Event
	subscribers map[domain.EventType][]chan<- domain.Event

	log *logger.Logger
}

func New(log *slog.Logger) *EventManager {
	return &EventManager{
		receiver:    make(chan domain.Event),
		subscribers: make(map[domain.EventType][]chan<- domain.Event),

		log: logger.New(log, "EVENT_MANAGER"),
	}
}

func (m *EventManager) PublishAccess() chan<- domain.Event {
	return m.receiver
}

func (m *EventManager) Subscribe(t domain.EventType) <-chan domain.Event {
	m.log.Info("someone subscribed to events", slog.String("eventType", t.String()))

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
	m.log.Info("started")

	go func() {
		defer m.log.Info("stopped")

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
