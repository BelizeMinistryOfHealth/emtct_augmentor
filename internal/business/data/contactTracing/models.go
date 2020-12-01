package contactTracing

import (
	"database/sql"
	"time"
)

type ContactTracings struct {
	*sql.DB
}

func New(db *sql.DB) ContactTracings {
	return ContactTracings{db}
}

type ContactTracing struct {
	Id         string     `json:"id"`
	PatientId  int        `json:"patientId"`
	Test       string     `json:"test"`
	TestResult string     `json:"testResult"`
	Comments   string     `json:"comments"`
	Date       time.Time  `json:"date"`
	CreatedBy  string     `json:"createdBy"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedBy  string     `json:"updatedBy"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}
