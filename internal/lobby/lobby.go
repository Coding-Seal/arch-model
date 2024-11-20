package lobby

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/Coding-Seal/arch-model/internal/domain"
)

type Lobby struct {
	eventCh     chan<- domain.Event
	patientIDCh chan<- int
}

func (l *Lobby) publishNewPatient(patient domain.Patient) {
	l.eventCh <- domain.NewPatientEvent{
		Timestamp: time.Now(),
		Patient:   patient,
	}
}

func (l *Lobby) newRandomPatient() domain.Patient {
	return domain.Patient{
		ID:   rand.Int(),
		Name: "Andy", // TODO: use fakeIt or smth like that
	}
}

func (l *Lobby) sendPatientToBench(patient domain.Patient) {
	l.patientIDCh <- patient.ID
}

func (l *Lobby) GeneratePatient() {
	patient := l.newRandomPatient()
	l.publishNewPatient(patient)
	l.sendPatientToBench(patient)
}

func (l *Lobby) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second / 2) // FIXME: should configure in New
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				l.GeneratePatient()
			}
		}
	}()
}
