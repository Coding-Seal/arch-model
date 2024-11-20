package domain

import (
	"time"
)

type EventType int

const (
	NONE = iota
	NEW_PATIENT
	PATIENT_GONE
	PATIENT_IN_QUEUE
	APPOINTMENT_FINISHED
	APPOINTMENT_STARTED
)

type Event interface {
	Time() time.Time
	Type() EventType
}

type NewPatientEvent struct {
	Timestamp time.Time
	Patient   Patient
}

func (e NewPatientEvent) Type() EventType {
	return NEW_PATIENT
}

func (e NewPatientEvent) Time() time.Time {
	return e.Timestamp
}

type PatientGoneEvent struct {
	Timestamp time.Time
	PatientID int
}

func (e PatientGoneEvent) Time() time.Time {
	return e.Timestamp
}

func (e PatientGoneEvent) Type() EventType {
	return PATIENT_GONE
}

type PatientInQueueEvent struct {
	Timestamp time.Time
	PatientID int
}

func (e PatientInQueueEvent) Time() time.Time {
	return e.Timestamp
}

func (e PatientInQueueEvent) Type() EventType {
	return PATIENT_IN_QUEUE
}

type AppointmentFinishedEvent struct {
	Timestamp time.Time
	PatientID int
	DoctorID  int
}

func (e AppointmentFinishedEvent) Time() time.Time {
	return e.Timestamp
}

func (e AppointmentFinishedEvent) Type() EventType {
	return APPOINTMENT_FINISHED
}

type AppointmentStartedEvent struct {
	Timestamp time.Time
	PatientID int
	DoctorID  int
}

func (e AppointmentStartedEvent) Time() time.Time {
	return e.Timestamp
}

func (e AppointmentStartedEvent) Type() EventType {
	return APPOINTMENT_STARTED
}
