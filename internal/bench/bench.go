package bench

import (
	"sync"

	"github.com/Coding-Seal/arch-model/internal/domain"
	"github.com/Coding-Seal/arch-model/pkg/dequeue"
)

var errEmptyQueue = "should never occur empty queue"

type eventSender interface {
	PatientLeft(patient *domain.Patient)
	PatientCame(patient *domain.Patient)
}

type Bench struct {
	eventSender eventSender

	queue *dequeue.Dequeue[*domain.Patient]
	mu    sync.Mutex
}

func New(eventSender eventSender, numberOfSeats int) *Bench {
	return &Bench{
		eventSender: eventSender,
		queue:       dequeue.New[*domain.Patient](numberOfSeats),
	}
}

func (b *Bench) HandleNewPatient(patient *domain.Patient) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.queue.Full() {
		lastPatient, _ := b.queue.Back()
		b.queue.PopBack()
		b.eventSender.PatientLeft(lastPatient)
	}

	_ = b.queue.PushBack(patient)
	b.eventSender.PatientCame(patient)
}

func (b *Bench) NextPatient() *domain.Patient {
	b.mu.Lock()
	defer b.mu.Unlock()

	patient, ok := b.queue.Back()
	if !ok {
		panic(errEmptyQueue)
	}
	return patient
}
