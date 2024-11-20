package doctor

import (
	"context"
	"time"

	"github.com/Coding-Seal/arch-model/internal/domain"
)

type Doctor struct {
	id                      int           // TODO: initialize in New
	nextAppointmentDuration time.Duration // TODO: initialize in New

	eventCh     chan<- domain.Event // TODO: initialize in Register
	patientIDCh <-chan int          // TODO: initialize in Register
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

func (d *Doctor) HandleNewPatient(ctx context.Context, patientID int) {
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
				d.HandleNewPatient(ctx, patientID)
			}
		}
	}()
}
