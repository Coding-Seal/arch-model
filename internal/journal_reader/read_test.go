package journalreader_test

import (
	"bytes"
	"encoding/json"

	"github.com/Coding-Seal/arch-model/internal/domain"
	journal_reader "github.com/Coding-Seal/arch-model/internal/journal_reader"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Read", func() {
	DescribeTable("simple test", func(e domain.Event) {
		b, err := json.Marshal(e)
		Expect(err).ShouldNot(HaveOccurred())

		buff := bytes.NewBuffer(b)
		reader := journal_reader.NewJournal(buff)
		events, err := reader.Read()

		Expect(err).ShouldNot(HaveOccurred())
		Expect(events).To(ConsistOf(e))
	},
		Entry("1", domain.NewPatientEvent{
			Patient: domain.Patient{ID: 1, Name: "John"},
			LobbyID: 1, EventTimer: domain.EventTimer{EventType: domain.NEW_PATIENT},
		}),
		Entry("2", domain.PatientGoneEvent{
			PatientID: 1, LobbyID: 1,
			EventTimer: domain.EventTimer{EventType: domain.PATIENT_GONE},
		}),
		Entry("3", domain.PatientInQueueEvent{
			PatientID: 1, LobbyID: 1,
			EventTimer: domain.EventTimer{EventType: domain.PATIENT_IN_QUEUE},
		}),
		Entry("4", domain.AppointmentStartedEvent{
			PatientID: 1, DoctorID: 1,
			EventTimer: domain.EventTimer{EventType: domain.APPOINTMENT_STARTED},
		}),
		Entry("5", domain.AppointmentFinishedEvent{
			PatientID: 1, DoctorID: 1,
			EventTimer: domain.EventTimer{EventType: domain.APPOINTMENT_FINISHED},
		}),
	)
})
