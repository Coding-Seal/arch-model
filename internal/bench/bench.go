package bench

import (
	"log/slog"
	"sync"

	"github.com/Coding-Seal/arch-model/internal/domain"
	"github.com/Coding-Seal/arch-model/pkg/dequeue"
	"github.com/Coding-Seal/arch-model/pkg/logger"
)

type Bench struct {
	queue *dequeue.Dequeue[int]
	mu    sync.Mutex
	log   *logger.Logger
}

func New(log *slog.Logger, numberOfSeats int) *Bench {
	return &Bench{
		log:   logger.New(log, "BENCH"),
		queue: dequeue.New[int](numberOfSeats),
	}
}

func (b *Bench) PushPatientID(patientID int) (int, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if patientID < 1 {
		panic("should never happen")
	}

	if b.queue.Full() {
		lastPatientID, _ := b.queue.Back()
		b.log.Debug("queue is full, some one have to go",
			slog.Int("lastPatientID", lastPatientID), slog.Int("newPatientID", patientID))
		b.queue.PopBack()

		return lastPatientID, false
	}

	_ = b.queue.PushBack(patientID)
	b.log.Debug("pushed patient", slog.Int("patientID", patientID))

	return 0, true
}

func (b *Bench) NextPatientID() (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	patientID, ok := b.queue.Front()
	if !ok {
		b.log.Warning("tried to advance queue, it is empty")

		return patientID, domain.ErrEmptyQueue
	}

	if patientID < 1 {
		panic("should never happen")
	}

	b.queue.PopFront()
	b.log.Debug("advanced queue", slog.Int("patientID", patientID))

	return patientID, nil
}
