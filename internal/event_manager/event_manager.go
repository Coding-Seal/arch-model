package eventmanager

import (
	"context"
	"log/slog"

	"github.com/Coding-Seal/arch-model/internal/domain"
	"github.com/Coding-Seal/arch-model/pkg/logger"
)

type EventManager struct {
	cancel context.CancelFunc
	done   chan struct{}

	receiver    chan domain.Event
	subscribers map[domain.EventType][]chan<- domain.Event

	log *logger.Logger
}

func New(log *slog.Logger) *EventManager {
	return &EventManager{
		done:        make(chan struct{}, 1),
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
	ctx, cancel := context.WithCancel(ctx)
	m.cancel = cancel

	m.log.Info("started serving")

	go func() {
		defer m.log.Info("stopped serving")
		defer func() { m.done <- struct{}{} }()

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

func (m *EventManager) Stop() {
	m.cancel()
	<-m.done
}
