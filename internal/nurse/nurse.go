package nurse

import (
	"context"

	"github.com/Coding-Seal/arch-model/internal/domain"
)

type PatientGetter interface {
	NextPatientID() (int, error)
}

type Nurse struct {
	PatentGoneCh          <-chan domain.Event
	PatentInQueueCh       <-chan domain.Event
	AppointmentFinishedCh <-chan domain.Event

	doctors        []*doctor
	previousDoctor int

	patientGetter PatientGetter

	numPatients int32
}

type doctor struct {
	ch   chan<- int
	busy bool
}

func (n *Nurse) RegisterDoctor() (<-chan int, int) {
	ch := make(chan int)
	doctorID := len(n.doctors)
	n.doctors = append(n.doctors, &doctor{ch: ch, busy: false})

	return ch, doctorID
}

func (n *Nurse) findAvailableDoctor() (*doctor, error) {
	var i = (n.previousDoctor + 1) % len(n.doctors)
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
			case <-n.PatentInQueueCh:
				n.handlePatientInQueueEvent()
			case <-n.PatentGoneCh:
				n.handlePatientGoneEvent()
			case event := <-n.PatentInQueueCh:
				n.handleAppointmentFinishedEvent(event)
			}
		}
	}()
}
