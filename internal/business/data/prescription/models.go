package prescription

import "time"

type Prescription struct {
	Id             int       `json:"id"`
	PatientId      int       `json:"patientId"`
	TotalDoses     int       `json:"totalDoses"`
	Pharmaceutical string    `json:"pharmaceutical"`
	Frequency      string    `json:"frequency"`
	Strength       string    `json:"strength"`
	Comments       string    `json:"comments"`
	PrescribedTime time.Time `json:"prescribedTime"`
}
