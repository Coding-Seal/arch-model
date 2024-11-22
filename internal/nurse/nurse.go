package nurse

import (
	"context"
	"log/slog"

	"github.com/Coding-Seal/arch-model/internal/domain"
	"github.com/Coding-Seal/arch-model/pkg/logger"
)

type patientGetter interface {
	NextPatientID() (int, error)
}

type subscriber interface {
	Subscribe(t domain.EventType) <-chan domain.Event
}

type Nurse struct {
	patentGoneCh          <-chan domain.Event
	patentInQueueCh       <-chan domain.Event
	appointmentFinishedCh <-chan domain.Event

	doctors        []*doctor
	previousDoctor int

	patientGetter patientGetter
	numPatients   int32

	log *logger.Logger
}

func New(log *slog.Logger, patientGetter patientGetter) *Nurse {
	return &Nurse{patientGetter: patientGetter, log: logger.New(log, "NURSE")}
}

type doctor struct {
	ch   chan<- int
	busy bool
}

func (n *Nurse) Register(subscriber subscriber) {
	n.log.Info("subscribed to events")
	n.patentGoneCh = subscriber.Subscribe(domain.PATIENT_GONE)
	n.patentInQueueCh = subscriber.Subscribe(domain.PATIENT_IN_QUEUE)
	n.appointmentFinishedCh = subscriber.Subscribe(domain.APPOINTMENT_FINISHED)
}

func (n *Nurse) RegisterDoctor() (<-chan int, int) {
	ch := make(chan int)
	doctorID := len(n.doctors)
	n.doctors = append(n.doctors, &doctor{ch: ch, busy: false})
	n.log.Info("registered doctor", slog.Int("doctorID", doctorID))

	return ch, doctorID
}

func (n *Nurse) findAvailableDoctor() (*doctor, error) {
	i := (n.previousDoctor + 1) % len(n.doctors)
	for ; i != n.previousDoctor; i = (i + 1) % len(n.doctors) {
		if !n.doctors[i].busy {
			n.previousDoctor = i

			return n.doctors[i], nil
		}
	}

	if !n.doctors[i].busy {
		n.previousDoctor = i

		return n.doctors[i], nil
	}

	return nil, domain.ErrAllDoctorsBusy
}

func (n *Nurse) handlePatientGoneEvent() {
	n.numPatients--
}

func (n *Nurse) handlePatientInQueueEvent() {
	n.numPatients++

	err := n.sendPatientToDoctor()
	if err != nil {
		// TODO: Log error
	}
}

func (n *Nurse) handleAppointmentFinishedEvent(event domain.Event) {
	e, ok := event.(domain.AppointmentFinishedEvent)
	if !ok {
		// TODO: Log error
	}

	if e.DoctorID >= len(n.doctors) {
		// TODO: Log error
	}

	n.doctors[e.DoctorID].busy = false

	if n.numPatients <= 0 {
		return
	}

	err := n.sendPatientToDoctor()
	if err != nil {
		// TODO: Log error
	}
}

func (n *Nurse) sendPatientToDoctor() error {
	doc, err := n.findAvailableDoctor()
	if err != nil {
		return err
	}

	patientID, err := n.patientGetter.NextPatientID()
	if err != nil {
		return err
	}

	n.numPatients--
	doc.ch <- patientID
	doc.busy = true

	return nil
}

func (n *Nurse) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-n.patentInQueueCh:
				n.handlePatientInQueueEvent()
			case <-n.patentGoneCh:
				n.handlePatientGoneEvent()
			case event := <-n.patentInQueueCh:
				n.handleAppointmentFinishedEvent(event)
			}
		}
	}()
}
