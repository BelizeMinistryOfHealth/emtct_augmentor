package contraceptives

import (
	"database/sql"
	"time"
)

type Contraceptives struct {
	*sql.DB
}

func New(db *sql.DB) Contraceptives {
	return Contraceptives{db}
}

type ContraceptiveUsed struct {
	Id             string     `json:"id"`
	PatientId      int        `json:"patientId"`
	MchEncounterId int        `json:"mchEncounterId"`
	Contraceptive  string     `json:"contraceptive"`
	Comments       string     `json:"comments"`
	DateUsed       time.Time  `json:"dateUsed"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
	CreatedBy      string     `json:"createdBy"`
	UpdatedBy      *string    `json:"updatedBy"`
}
