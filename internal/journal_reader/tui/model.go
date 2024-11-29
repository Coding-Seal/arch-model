package journalreader

import "github.com/Coding-Seal/arch-model/internal/domain"

type Model struct {
	System   System
	Patients map[int]domain.Patient
	Events   []domain.Event

	EventsDisplay []string
	Cursor        int
}

type System struct {
	Doctors []Doctor
	Lobbies []Lobby
	Nurse   Nurse
	Bench   Bench
}

type Doctor struct {
	ID          int
	Busy        bool
	PatientName string
}

type Lobby struct {
	ID          int
	Activated   bool
	PatientName string
}

type Nurse struct {
	Activated bool
}

type Bench struct {
	Cap     int
	Present int
	State   int
}

const (
	Nothing = iota
	Fail
	Pass
)
