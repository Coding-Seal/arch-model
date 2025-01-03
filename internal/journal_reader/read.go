package journalreader

import (
	"io"
	"slices"

	"github.com/Coding-Seal/arch-model/internal/domain"
	"github.com/Coding-Seal/arch-model/pkg/jsonl"
)

type Journal struct {
	sc *jsonl.Scanner
}

func NewJournal(r io.Reader) *Journal {
	return &Journal{
		sc: jsonl.NewScanner(r),
	}
}

func (j *Journal) Read() ([]domain.Event, error) {
	data := make([]domain.Event, 0)

	var newPatient domain.NewPatientEvent

	var patientGone domain.PatientGoneEvent

	var patientInQueue domain.PatientInQueueEvent

	var appointmentFinished domain.AppointmentFinishedEvent

	var appointmentStarted domain.AppointmentStartedEvent

	for j.sc.Scan() {
		err := j.sc.Json(&newPatient)
		if err == nil && newPatient.Type() == domain.NEW_PATIENT {
			data = append(data, newPatient)
		}

		err = j.sc.Json(&patientGone)
		if err == nil && patientGone.Type() == domain.PATIENT_GONE {
			data = append(data, patientGone)
		}

		err = j.sc.Json(&patientInQueue)
		if err == nil && patientInQueue.Type() == domain.PATIENT_IN_QUEUE {
			data = append(data, patientInQueue)
		}

		err = j.sc.Json(&appointmentFinished)
		if err == nil && appointmentFinished.Type() == domain.APPOINTMENT_FINISHED {
			data = append(data, appointmentFinished)
		}

		err = j.sc.Json(&appointmentStarted)
		if err == nil && appointmentStarted.Type() == domain.APPOINTMENT_STARTED {
			data = append(data, appointmentStarted)
		}

		if err != nil {
			return nil, err
		}
	}

	slices.SortFunc(data, func(a, b domain.Event) int {
		return a.ID() - b.ID()
	})

	return data, nil
}
