package bench

import (
	"context"
	"sync"
	"time"

	"github.com/Coding-Seal/arch-model/internal/domain"
	"github.com/Coding-Seal/arch-model/pkg/dequeue"
)

type Bench struct {
	eventCh     chan<- domain.Event // TODO: should initialize in Register
	patientIDCh <-chan int

	queue *dequeue.Dequeue[int]
	mu    sync.Mutex
}

func New(numberOfSeats int) *Bench {
	return &Bench{
		queue: dequeue.New[int](numberOfSeats),
	}
}

func (b *Bench) handleNewPatient(patientID int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.queue.Full() {
		lastPatientID, _ := b.queue.Back()
		b.queue.PopBack()
		b.publishPatientGone(lastPatientID)
	}

	_ = b.queue.PushBack(patientID)
	b.publishPatientInQueue(patientID)
}

func (b *Bench) NextPatientID() (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	patientID, ok := b.queue.Back()
	if ok {
		b.queue.PopFront()
	}

	return patientID, domain.ErrEmptyQueue
}

func (b *Bench) publishPatientGone(patientID int) {
	b.eventCh <- domain.PatientGoneEvent{
		Timestamp: time.Now(),
		PatientID: patientID,
	}
}

func (b *Bench) publishPatientInQueue(patientID int) {
	b.eventCh <- domain.PatientInQueueEvent{
		Timestamp: time.Now(),
		PatientID: patientID,
	}
}

func (b *Bench) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case newPatientID := <-b.patientIDCh:
				b.handleNewPatient(newPatientID)
			case <-ctx.Done():
				return
			}
		}
	}()
}
