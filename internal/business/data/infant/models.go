package infant

import (
	"context"
	"moh.gov.bz/mch/emtct/internal/db"
	"time"

	"moh.gov.bz/mch/emtct/internal/business/data/person"
)

type DbCollections struct {
	HivScreening string
}

type Infants struct {
	firestore   *db.FirestoreClient
	collections DbCollections
}

func New(firestore *db.FirestoreClient) Infants {
	colls := DbCollections{HivScreening: "hivScreenings"}
	return Infants{firestore: firestore, collections: colls}
}

func (i *Infants) ctx() context.Context {
	return i.firestore.Ctx
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
	Id                     string     `json:"id" firestore:"id"`
	PatientId              string     `json:"patientId" firestore:"patientId"`
	MotherId               string     `json:"motherId" firestore:"motherId"`
	PregnancyId            int        `json:"pregnancyId" firestore:"pregnancyId"`
	TestName               string     `json:"testName" firestore:"testName"`
	ScreeningDate          time.Time  `json:"screeningDate" firestore:"screeningDate"`
	DateSampleReceivedAtHq *time.Time `json:"dateSampleReceivedAtHq,omitEmpty" firestore:"dateSampleReceivedAtHq"`
	SampleCode             string     `json:"sampleCode" firestore:"sampleCode"`
	DateSampleShipped      *time.Time `json:"dateSampleShipped" firestore:"dateSampleShipped"`
	Destination            string     `json:"destination" firestore:"destination"`
	DateResultReceived     *time.Time `json:"dateResultReceived,omitEmpty" firestore:"dateResultReceived"`
	DateSampleTaken        *time.Time `json:"dateSampleTaken,omitEmpty" firestore:"dateSampleTaken"`
	DueDate                *time.Time `json:"dueDate,omitEmpty" firestore:"dueDate"`
	Result                 string     `json:"result" firestore:"result"`
	DateResultShared       *time.Time `json:"dateResultShared,omitEmpty" firestore:"dateResultShared"`
	Timely                 bool       `json:"timely" firestore:"timely"`
	CreatedAt              time.Time  `json:"createdAt" firestore:"createdAt"`
	UpdatedAt              *time.Time `json:"updatedAt" firestore:"updatedAt"`
	CreatedBy              string     `json:"createdBy" firestore:"createdBy"`
	UpdatedBy              *string    `json:"updatedBy" firestore:"updatedBy"`
}
