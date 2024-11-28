package doctor

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Coding-Seal/arch-model/internal/domain"
	"github.com/Coding-Seal/arch-model/pkg/logger"
)

type accessGetter interface {
	PublishAccess() chan<- domain.Event
}

type doctorRegisterer interface {
	RegisterDoctor() (<-chan int, int)
}

type Doctor struct {
	cancel context.CancelFunc
	done   chan struct{}

	id                      int
	nextAppointmentDuration time.Duration

	eventCh     chan<- domain.Event
	patientIDCh <-chan int

	log *logger.Logger
}

func New(log *slog.Logger, firstAppointment time.Duration) *Doctor {
	return &Doctor{
		done:                    make(chan struct{}, 1),
		nextAppointmentDuration: firstAppointment,
		log:                     logger.New(log, "DOCTOR"),
	}
}

func (d *Doctor) Register(accessGetter accessGetter, doctorRegisterer doctorRegisterer) {
	d.patientIDCh, d.id = doctorRegisterer.RegisterDoctor()
	d.log.SetServiceName(fmt.Sprintf("DOCTOR_%d", d.id))
	d.log.Info("registered nurse")
	d.eventCh = accessGetter.PublishAccess()
	d.log.Info("registered event chan")
}

func (d *Doctor) publishAppointmentFinished(patientID int) {
	d.log.Debug("publish event", slog.String("eventType", domain.APPOINTMENT_FINISHED.String()))
	d.eventCh <- domain.AppointmentFinishedEvent{
		Timestamp: time.Now(),
		DoctorID:  d.id,
		PatientID: patientID,
	}
}

func (d *Doctor) publishAppointmentStarted(patientID int) {
	d.log.Debug("publish event", slog.String("eventType", domain.APPOINTMENT_STARTED.String()))
	d.eventCh <- domain.AppointmentStartedEvent{
		Timestamp: time.Now(),
		DoctorID:  d.id,
		PatientID: patientID,
	}
}

func (d *Doctor) handleNewPatient(ctx context.Context, patientID int) {
	d.publishAppointmentStarted(patientID)
	ch := time.After(d.nextAppointmentDuration)
	select {
	case <-ctx.Done():
		return
	case <-ch:
		d.nextAppointmentDuration *= 2
		d.publishAppointmentFinished(patientID)

		return
	}
}

func (d *Doctor) Run(ctx context.Context) {
	d.log.Info("started serving")
	ctx, cancel := context.WithCancel(ctx)
	d.cancel = cancel

	go func() {
		defer d.log.Info("stopped serving")
		defer func() { d.done <- struct{}{} }()

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

func (d *Doctor) Stop() {
	d.cancel()
	<-d.done
}
