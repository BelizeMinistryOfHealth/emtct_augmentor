package homeVisits

import (
	"context"
	"moh.gov.bz/mch/emtct/internal/db"
	"time"
)

type HomeVisits struct {
	db         *db.FirestoreClient
	collection string
}

func New(db *db.FirestoreClient) HomeVisits {
	return HomeVisits{db: db, collection: "homeVisits"}
}

func (h HomeVisits) ctx() context.Context {
	return h.db.Ctx
}

type HomeVisit struct {
	Id          string     `json:"id" firestore:"id"`
	PatientId   int        `json:"patientId" firestore:"patientId"`
	PregnancyId int        `json:"pregnancyId" firestore:"pregnancyId"`
	Reason      string     `json:"reason" firestore:"reason"`
	Comments    string     `json:"comments" firestore:"comments"`
	DateOfVisit time.Time  `json:"dateOfVisit" firestore:"dateOfVisit"`
	CreatedAt   time.Time  `json:"createdAt" firestore:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt" firestore:"updatedAt"`
	CreatedBy   string     `json:"createdBy" firestore:"createdBy"`
	UpdatedBy   *string    `json:"updatedBy" firestore:"updatedBy"`
}
