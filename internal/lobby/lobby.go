package lobby

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/Coding-Seal/arch-model/internal/domain"
	"github.com/Coding-Seal/arch-model/pkg/logger"
)

type patientSender interface {
	PushPatientID(patientID int) (int, bool)
}

type accessGetter interface {
	PublishAccess() chan<- domain.Event
}

type Lobby struct {
	eventCh       chan<- domain.Event
	patientSender patientSender

	log *logger.Logger
}

func New(log *slog.Logger, patientSender patientSender, ID int) *Lobby {
	return &Lobby{patientSender: patientSender, log: logger.New(log, fmt.Sprintf("LOBBY_%d", ID))}
}

func (l *Lobby) Register(accessGetter accessGetter) {
	l.eventCh = accessGetter.PublishAccess()
	l.log.Info("registered event chan")
}

func (l *Lobby) publishNewPatient(patient domain.Patient) {
	l.log.Debug("sent event", slog.String("eventType", domain.NEW_PATIENT.String()))
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

func (l *Lobby) publishPatientGone(patientID int) {
	l.log.Debug("sent event", slog.String("eventType", domain.PATIENT_GONE.String()))
	l.eventCh <- domain.PatientGoneEvent{
		Timestamp: time.Now(),
		PatientID: patientID,
	}
}

func (l *Lobby) publishPatientInQueue(patientID int) {
	l.log.Debug("sent event", slog.String("eventType", domain.PATIENT_IN_QUEUE.String()))
	l.eventCh <- domain.PatientInQueueEvent{
		Timestamp: time.Now(),
		PatientID: patientID,
	}
}

func (l *Lobby) sendPatientToBench(patient domain.Patient) {
	lastID, ok := l.patientSender.PushPatientID(patient.ID)
	l.log.Debug("pushed patient into queue", slog.Int("patientID", patient.ID))
	if !ok {
		l.log.Debug("Patient gone", slog.Int("patientID", lastID))
		l.publishPatientGone(lastID)
	}

	l.publishPatientInQueue(patient.ID)
}

func (l *Lobby) generatePatient() {
	patient := l.newRandomPatient()
	l.log.Debug("generated new patient", slog.Any("patient", patient))
	l.publishNewPatient(patient)
	l.sendPatientToBench(patient)
}

func (l *Lobby) Run(ctx context.Context) {
	l.log.Info("Started")

	go func() {
		ticker := time.NewTicker(time.Second / 2) // FIXME: should configure in New
		defer ticker.Stop()
		defer l.log.Info("Stopped")
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				l.generatePatient()
			}
		}
	}()
}
