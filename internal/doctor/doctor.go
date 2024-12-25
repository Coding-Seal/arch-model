package doctor

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
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

	id                       int
	firstAppointmentDuration time.Duration

	eventCh     chan<- domain.Event
	patientIDCh <-chan int

	log *logger.Logger
}

func New(log *slog.Logger, firstAppointment time.Duration) *Doctor {
	return &Doctor{
		done:                     make(chan struct{}, 1),
		firstAppointmentDuration: firstAppointment,
		log:                      logger.New(log, "DOCTOR"),
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
	d.log.Debug("appointment finished", slog.Int("patientID", patientID))
	d.eventCh <- domain.NewAppointmentFinishedEvent(patientID, d.id)
}

func (d *Doctor) publishAppointmentStarted(patientID int) {
	d.log.Debug("publish event", slog.String("eventType", domain.APPOINTMENT_STARTED.String()))
	d.log.Debug("appointment started", slog.Int("patientID", patientID))
	d.eventCh <- domain.NewAppointmentStartedEvent(patientID, d.id)
}

func (d *Doctor) handleNewPatient(ctx context.Context, patientID int) {
	d.publishAppointmentStarted(patientID)

	ticker := time.NewTicker(getAppointmentDuration(d.firstAppointmentDuration, patientID))
	select {
	case <-ctx.Done():
		d.publishAppointmentFinished(patientID)
		ticker.Stop()

		return
	case <-ticker.C:
		d.publishAppointmentFinished(patientID)

		return
	}
}

// func getAppointmentDuration(con time.Duration, id int) time.Duration {
// 	return con * time.Duration(math.Pow(2, float64(id)))
// }

func getAppointmentDuration(con time.Duration, _ int) time.Duration {
	r := rand.ExpFloat64()

	return time.Duration(r * float64(con))
}

func (d *Doctor) Run(ctx context.Context) {
	d.log.Info("started serving")

	ctx, cancel := context.WithCancel(ctx)
	d.cancel = cancel

	go func() {
		defer d.log.Info("stopped serving")
		defer func() { close(d.done) }()

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
