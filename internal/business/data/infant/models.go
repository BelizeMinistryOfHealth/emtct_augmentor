package infant

import (
	"database/sql"
	"time"

	"moh.gov.bz/mch/emtct/internal/business/data/person"
)

type Infants struct {
	Acsis *sql.DB
}

func New(db *sql.DB) Infants {
	return Infants{Acsis: db}
}

type Diagnoses struct {
	DiagnosisId int       `json:"diagnosisId"`
	PatientId   int       `json:"patientId"`
	Diagnosis   string    `json:"diagnosis"`
	Doctor      string    `json:"doctor"`
	Comments    string    `json:"comments"`
	Date        time.Time `json:"date"`
}

type Infant struct {
	Infant person.Person `json:"infant"`
	Mother person.Person `json:"mother"`
}

type HivScreening struct {
	Id                     string     `json:"id"`
	PatientId              int        `json:"patientId"`
	MotherId               int        `json:"motherId"`
	TestName               string     `json:"testName"`
	ScreeningDate          time.Time  `json:"screeningDate"`
	DateSampleReceivedAtHq *time.Time `json:"dateSampleReceivedAtHq,omitEmpty"`
	SampleCode             string     `json:"sampleCode"`
	DateSampleShipped      *time.Time `json:"dateSampleShipped"`
	Destination            string     `json:"destination"`
	DateResultReceived     *time.Time `json:"dateResultReceived,omitEmpty"`
	DateSampleTaken        *time.Time `json:"dateSampleTaken,omitEmpty"`
	DueDate                *time.Time `json:"dueDate,omitEmpty"`
	Result                 string     `json:"result"`
	DateResultShared       *time.Time `json:"dateResultShared,omitEmpty"`
	Timely                 bool       `json:"timely"`
	CreatedAt              time.Time  `json:"createdAt"`
	UpdatedAt              *time.Time `json:"updatedAt"`
	CreatedBy              string     `json:"createdBy"`
	UpdatedBy              *string    `json:"updatedBy"`
}
