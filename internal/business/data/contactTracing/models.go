package contactTracing

import (
	"context"
	"moh.gov.bz/mch/emtct/internal/db"
	"time"
)

type ContactTracings struct {
	db         *db.FirestoreClient
	collection string
}

func New(db *db.FirestoreClient) ContactTracings {
	return ContactTracings{db: db, collection: "contactTracings"}
}

func (c *ContactTracings) ctx() context.Context {
	return c.db.Ctx
}

type ContactTracing struct {
	ID          string     `json:"id" firestore:"id"`
	PatientId   int        `json:"patientId" firestore:"patientId"`
	PregnancyId int        `json:"pregnancyId" firestore:"pregnancyId"`
	Test        string     `json:"test" firestore:"test"`
	TestResult  string     `json:"testResult" firestore:"testResult"`
	Comments    string     `json:"comments" firestore:"comments"`
	Date        time.Time  `json:"date" firestore:"date"`
	CreatedBy   string     `json:"createdBy" firestore:"createdBy"`
	CreatedAt   time.Time  `json:"createdAt" firestore:"createdAt"`
	UpdatedBy   string     `json:"updatedBy" firestore:"updatedBy"`
	UpdatedAt   *time.Time `json:"updatedAt" firestore:"updatedAt"`
}
