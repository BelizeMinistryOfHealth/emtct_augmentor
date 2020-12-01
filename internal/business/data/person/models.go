package person

import "time"

type Person struct {
	PatientId  int        `json:"patientId"`
	FirstName  string     `json:"firstName"`
	LastName   string     `json:"lastName"`
	MiddleName string     `json:"middleName"`
	Dob        *time.Time `json:"dob"`
}
