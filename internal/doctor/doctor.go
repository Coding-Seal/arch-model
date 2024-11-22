package doctor

import (
	"context"
	"time"

	"github.com/Coding-Seal/arch-model/internal/domain"
)

type accessGetter interface {
	PublishAccess() chan<- domain.Event
}

type doctorRegisterer interface {
	RegisterDoctor() (<-chan int, int)
}

type Doctor struct {
	id                      int
	nextAppointmentDuration time.Duration

	eventCh     chan<- domain.Event
	patientIDCh <-chan int
}

func New(firstAppointment time.Duration) *Doctor {
	return &Doctor{nextAppointmentDuration: firstAppointment}
}

func (d *Doctor) Register(accessGetter accessGetter, doctorRegisterer doctorRegisterer) {
	d.eventCh = accessGetter.PublishAccess()
	d.patientIDCh, d.id = doctorRegisterer.RegisterDoctor()
}

func (d *Doctor) publishAppointmentFinished(patientID int) {
	d.eventCh <- domain.AppointmentFinishedEvent{
		Timestamp: time.Now(),
		DoctorID:  d.id,
		PatientID: patientID,
	}
}

func (d *Doctor) publishAppointmentStarted(patientID int) {
	d.eventCh <- domain.AppointmentStartedEvent{
		Timestamp: time.Now(),
		DoctorID:  d.id,
		PatientID: patientID,
	}
}

func (d *Doctor) handleNewPatient(ctx context.Context, patientID int) {
	d.publishAppointmentStarted(patientID)
	select {
	case <-ctx.Done():
		return
	case <-time.After(d.nextAppointmentDuration):
		break
	}
	d.publishAppointmentFinished(patientID)
}

func (d *Doctor) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case patientID := <-d.patientIDCh:
				d.handleNewPatient(ctx, patientID)
			}
		}
	}()
}
