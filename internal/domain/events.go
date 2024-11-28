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

type NewPatientEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Patient   Patient   `json:"patient_id"`
	LobbyID   int       `json:"lobby_id"`
}

func (e NewPatientEvent) Type() EventType {
	return NEW_PATIENT
}

func (e NewPatientEvent) Time() time.Time {
	return e.Timestamp
}

type PatientGoneEvent struct {
	Timestamp time.Time `json:"timestamp"`
	PatientID int       `json:"patient_id"`
	LobbyID   int       `json:"lobby_id"`
}

func (e PatientGoneEvent) Time() time.Time {
	return e.Timestamp
}

func (e PatientGoneEvent) Type() EventType {
	return PATIENT_GONE
}

type PatientInQueueEvent struct {
	Timestamp time.Time `json:"timestamp"`
	PatientID int       `json:"patient_id"`
	LobbyID   int       `json:"lobby_id"`
}

func (e PatientInQueueEvent) Time() time.Time {
	return e.Timestamp
}

func (e PatientInQueueEvent) Type() EventType {
	return PATIENT_IN_QUEUE
}

type AppointmentFinishedEvent struct {
	Timestamp time.Time `json:"timestamp"`
	PatientID int       `json:"patient_id"`
	DoctorID  int       `json:"doctor_id"`
}

func (e AppointmentFinishedEvent) Time() time.Time {
	return e.Timestamp
}

func (e AppointmentFinishedEvent) Type() EventType {
	return APPOINTMENT_FINISHED
}

type AppointmentStartedEvent struct {
	Timestamp time.Time `json:"timestamp"`
	PatientID int       `json:"patient_id"`
	DoctorID  int       `json:"doctor_id"`
}

func (e AppointmentStartedEvent) Time() time.Time {
	return e.Timestamp
}

func (e AppointmentStartedEvent) Type() EventType {
	return APPOINTMENT_STARTED
}
