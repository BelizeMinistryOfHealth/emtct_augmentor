package contraceptives

import (
	"context"
	"moh.gov.bz/mch/emtct/internal/db"
	"time"
)

type Contraceptives struct {
	db         *db.FirestoreClient
	collection string
}

func New(db *db.FirestoreClient) Contraceptives {
	return Contraceptives{db: db, collection: "contraceptives"}
}

func (c *Contraceptives) ctx() context.Context {
	return c.db.Ctx
}

type ContraceptiveUsed struct {
	ID            string     `json:"id" firestore:"id"`
	PatientId     int        `json:"patientId" firestore:"patientId"`
	PregnancyId   int        `json:"pregnancyId" firestore:"pregnancyId"`
	Contraceptive string     `json:"contraceptive" firestore:"contraceptive"`
	Comments      string     `json:"comments" firestore:"comments"`
	DateUsed      time.Time  `json:"dateUsed" firestore:"dateUsed"`
	CreatedAt     time.Time  `json:"createdAt" firestore:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt" firestore:"updatedAt"`
	CreatedBy     string     `json:"createdBy" firestore:"createdBy"`
	UpdatedBy     *string    `json:"updatedBy" firestore:"updatedBy"`
}
