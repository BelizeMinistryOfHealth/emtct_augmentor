package admissions

import (
	"context"
	"moh.gov.bz/mch/emtct/internal/db"
	"time"
)

type Admissions struct {
	db         *db.FirestoreClient
	collection string
}

func New(db *db.FirestoreClient) Admissions {
	return Admissions{db: db, collection: "admissions"}
}

func (a Admissions) ctx() context.Context {
	return a.db.Ctx
}

type HospitalAdmission struct {
	ID           string     `json:"id" firestore:"id"`
	PatientId    int        `json:"patientId" firestore:"patientId"`
	PregnancyId  int        `json:"pregnancyId" firestore:"pregnancyId"`
	DateAdmitted time.Time  `json:"dateAdmitted" firestore:"dateAdmitted"`
	Facility     string     `json:"facility" firestore:"facility"`
	Reason       string     `json:"reason" firestore:"reason"`
	CreatedAt    time.Time  `json:"createdAt" firestore:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt" firestore:"updatedAt"`
	CreatedBy    string     `json:"createdBy" firestore:"createdBy"`
	UpdatedBy    *string    `json:"updatedBy" firestore:"updatedBy"`
}
