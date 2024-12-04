package domain

import (
	"time"
)

type EventType int

const (
	NONE EventType = iota
	NEW_PATIENT
	PATIENT_GONE
	PATIENT_IN_QUEUE
	APPOINTMENT_FINISHED
	APPOINTMENT_STARTED
)

func (t EventType) String() string {
	switch t {
	case NEW_PATIENT:
		return "NewPatient"
	case PATIENT_GONE:
		return "PatientGone"
	case PATIENT_IN_QUEUE:
		return "PatientInQueue"
	case APPOINTMENT_FINISHED:
		return "AppointmentFinished"
	case APPOINTMENT_STARTED:
		return "AppointmentStarted"
	default:
		return "None"
	}
}

type Event interface {
	Time() time.Time
	Type() EventType
}

type EventTimer struct {
	Timestamp time.Time `json:"timestamp"`
	EventType EventType `json:"eventType"`
}

func (e EventTimer) Type() EventType {
	return e.EventType
}

func (e EventTimer) Time() time.Time {
	return e.Timestamp
}

type NewPatientEvent struct {
	EventTimer
	Patient Patient `json:"patient"`
	LobbyID int     `json:"lobbyId"`
}

func NewNewPatientEvent(patient Patient, lobbyID int) NewPatientEvent {
	return NewPatientEvent{
		EventTimer: EventTimer{
			Timestamp: time.Now(),
			EventType: NEW_PATIENT,
		},
		Patient: patient,
		LobbyID: lobbyID,
	}
}

type PatientGoneEvent struct {
	EventTimer
	PatientID int `json:"patientId"`
	LobbyID   int `json:"lobbyId"`
}

func NewPatientGoneEvent(patientId int, lobbyID int) PatientGoneEvent {
	return PatientGoneEvent{
		EventTimer: EventTimer{
			Timestamp: time.Now(),
			EventType: PATIENT_GONE,
		},
		PatientID: patientId,
		LobbyID:   lobbyID,
	}
}

type PatientInQueueEvent struct {
	EventTimer
	PatientID int `json:"patientId"`
	LobbyID   int `json:"lobbyId"`
}

func NewPatientInQueueEvent(patientId int, lobbyID int) PatientInQueueEvent {
	return PatientInQueueEvent{
		EventTimer: EventTimer{
			Timestamp: time.Now(),
			EventType: PATIENT_IN_QUEUE,
		},
		PatientID: patientId,
		LobbyID:   lobbyID,
	}
}

type AppointmentFinishedEvent struct {
	EventTimer
	PatientID int `json:"patientId"`
	DoctorID  int `json:"doctorId"`
}

func NewAppointmentFinishedEvent(patientId int, doctorID int) AppointmentFinishedEvent {
	return AppointmentFinishedEvent{
		EventTimer: EventTimer{
			Timestamp: time.Now(),
			EventType: APPOINTMENT_FINISHED,
		},
		PatientID: patientId,
		DoctorID:  doctorID,
	}
}

type AppointmentStartedEvent struct {
	EventTimer
	PatientID int `json:"patientId"`
	DoctorID  int `json:"doctorId"`
}

func NewAppointmentStartedEvent(patientId int, doctorID int) AppointmentStartedEvent {
	return AppointmentStartedEvent{
		EventTimer: EventTimer{
			Timestamp: time.Now(),
			EventType: APPOINTMENT_STARTED,
		},
		PatientID: patientId,
		DoctorID:  doctorID,
	}
}
