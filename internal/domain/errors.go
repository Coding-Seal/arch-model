package domain

import "errors"

var (
	ErrEmptyQueue     = errors.New("no patients on bench")
	ErrAllDoctorsBusy = errors.New("no available doctors")
)
