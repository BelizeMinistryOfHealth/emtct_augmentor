package homeVisits

import (
	"database/sql"
	"time"
)

type HomeVisits struct {
	*sql.DB
}

func New(db *sql.DB) HomeVisits {
	return HomeVisits{db}
}

type HomeVisit struct {
	Id             string     `json:"id"`
	PatientId      int        `json:"patientId"`
	MchEncounterId int        `json:"mchEncounterId"`
	Reason         string     `json:"reason"`
	Comments       string     `json:"comments"`
	DateOfVisit    time.Time  `json:"dateOfVisit"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
	CreatedBy      string     `json:"createdBy"`
	UpdatedBy      *string    `json:"updatedBy"`
}
