package labs

import (
	"time"

	"moh.gov.bz/mch/emtct/internal/db"
)

type Labs struct {
	AcsisDb *db.AcsisDb
}

func New(acsisDb *db.AcsisDb) Labs {
	return Labs{AcsisDb: acsisDb}
}

type LabResult struct {
	Id                     int        `json:"id"`
	PatientId              int        `json:"patientId"`
	TestResult             string     `json:"testResult"`
	TestName               string     `json:"testName"`
	TestRequestId          int        `json:"testRequestId"`
	TestRequestItemId      int        `json:"testRequestItemId"`
	DateSampleTaken        *time.Time `json:"dateSampleTaken"`
	ResultDate             *time.Time `json:"resultDate"`
	ReleasedTime           *time.Time `json:"releasedTime"`
	DateOrderReceivedByLab *time.Time `json:"dateOrderReceivedByLab"`
}

type SampleTimeliness string

const (
	Timely       SampleTimeliness = "Timely"
	NotTimely    SampleTimeliness = "NotTimely"
	NotAvailable SampleTimeliness = "N/A"
)

type SyphilisScreening struct {
	Id                 int              `json:"id"`
	PatientId          int              `json:"patientId"`
	TestName           string           `json:"testName"`
	ScreeningDate      time.Time        `json:"screeningDate"`
	DateResultReceived *time.Time       `json:"dateResultReceived,omitEmpty"`
	DateSampleTaken    *time.Time       `json:"dateSampleTaken,omitEmpty"`
	DueDate            *time.Time       `json:"dueDate,omitEmpty"`
	Result             string           `json:"result"`
	DateResultShared   *time.Time       `json:"dateResultShared,omitEmpty"`
	Timely             SampleTimeliness `json:"timely"`
}
