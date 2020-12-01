package admissions

import (
	"database/sql"
	"time"
)

type Admissions struct {
	*sql.DB
}

func New(db *sql.DB) Admissions {
	return Admissions{db}
}

type HospitalAdmission struct {
	Id             string     `json:"id"`
	PatientId      int        `json:"patientId"`
	MchEncounterId int        `json:"mchEncounterId"`
	DateAdmitted   time.Time  `json:"dateAdmitted"`
	Facility       string     `json:"facility"`
	Reason         string     `json:"reason"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
	CreatedBy      string     `json:"createdBy"`
	UpdatedBy      *string    `json:"updatedBy"`
}
