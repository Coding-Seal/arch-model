package ui

import (
	"github.com/Coding-Seal/arch-model/internal/domain"
)

type System struct {
	Doctors map[int]*Doctor
	Lobbies map[int]*Lobby
	Nurse   *Nurse
	Bench   *Bench
	Utils   Utils
}
type Utils struct {
	PreviousStateID int
	NextStateID     int
	EventString     string
}

func NewSystem(numDoctors, numLobbies, benchCap int, currentStateID int, EventString string) *System {
	doctors := make(map[int]*Doctor, numDoctors)
	lobbies := make(map[int]*Lobby, numLobbies)

	for i := range numDoctors {
		doctors[i] = &Doctor{ID: i}
	}

	for i := range numLobbies {
		lobbies[i+1] = &Lobby{ID: i + 1}
	}

	return &System{
		Doctors: doctors,
		Lobbies: lobbies,
		Nurse:   &Nurse{},
		Bench: &Bench{
			Cap:      benchCap,
			Patients: make([]int, benchCap),
			Top:      -1,
		},
		Utils: Utils{
			PreviousStateID: currentStateID - 1,
			NextStateID:     currentStateID + 1,
			EventString:     EventString,
		},
	}
}

type Doctor struct {
	ID        int
	Busy      bool
	PatientID int
}

type Lobby struct {
	ID        int
	Activated bool
	PatientID int
}

type Nurse struct {
	Activated bool
}

type Bench struct {
	Cap      int
	State    int
	Patients []int
	Top      int
}

const (
	Nothing = iota
	Denial
	Good
	Advanced
)

func (s *System) clearLobbies() {
	for _, l := range s.Lobbies {
		l.Activated = false
	}
}

func (s *System) activateLobby(LobbyID, PatientID int) {
	lobby, ok := s.Lobbies[LobbyID]
	if !ok {
		lobby = &Lobby{ID: LobbyID}
		s.Lobbies[LobbyID] = lobby
	}

	lobby.Activated = true
}

func (s *System) toggleDoctor(DoctorID int, PatientID int) {
	doc, ok := s.Doctors[DoctorID]
	if !ok {
		doc = &Doctor{ID: DoctorID}
		s.Doctors[DoctorID] = doc
	}

	doc.Busy = !doc.Busy
	if doc.Busy {
		doc.PatientID = PatientID
	} else {
		doc.PatientID = 0
	}
}

func (s *System) ApplyEvent(e domain.Event) {
	s.clearLobbies()
	s.Bench.State = Nothing
	s.Nurse.Activated = false

	switch event := e.(type) {
	case domain.NewPatientEvent:
		s.activateLobby(event.LobbyID, event.Patient.ID)
		s.Lobbies[event.LobbyID].PatientID = event.Patient.ID
	case domain.PatientInQueueEvent:
		s.Bench.Top++
		s.Bench.Patients[s.Bench.Top] = event.PatientID
		s.activateLobby(event.LobbyID, event.PatientID)
		s.Bench.State = Good
	case domain.PatientGoneEvent:
		s.Bench.Patients = s.Bench.Patients[:len(s.Bench.Patients)-1]
		s.Bench.Patients = append(s.Bench.Patients, 0)
		s.Bench.Top--
		s.activateLobby(event.LobbyID, event.PatientID)
		s.Bench.State = Denial
	case domain.AppointmentFinishedEvent:
		s.toggleDoctor(event.DoctorID, event.PatientID)
	case domain.AppointmentStartedEvent:
		s.Nurse.Activated = true
		s.Bench.Patients = s.Bench.Patients[1:]
		s.Bench.Patients = append(s.Bench.Patients, 0)
		s.Bench.Top--
		s.Bench.State = Advanced
		s.toggleDoctor(event.DoctorID, event.PatientID)
	}
}
