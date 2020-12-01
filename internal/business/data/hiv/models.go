package hiv

import (
	"time"

	"moh.gov.bz/mch/emtct/internal/db"
)

type HIV struct {
	AcsisDb *db.AcsisDb
}

func New(db *db.AcsisDb) HIV {
	return HIV{AcsisDb: db}
}

type Diagnosis struct {
	Id        string    `json:"id"`
	PatientId string    `json:"patientId"`
	Date      time.Time `json:"date"`
	Name      string    `json:"name"`
}
