package pregnancy

import (
	"time"
)

type ObstetricHistory struct {
	Id             string    `json:"id"`
	PatientId      string    `json:"patientId"`
	Date           time.Time `json:"date"`
	ObstetricEvent string    `json:"event"`
}
