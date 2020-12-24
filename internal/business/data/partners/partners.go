package partners

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/api/iterator"

	"moh.gov.bz/mch/emtct/internal/business/data/prescription"
	"moh.gov.bz/mch/emtct/internal/db"
)

type Partners struct {
	db         *db.FirestoreClient
	collection string
}

func New(db *db.FirestoreClient) Partners {
	return Partners{db: db, collection: "partnerSyphilisTreatments"}
}

func (p *Partners) ctx() context.Context {
	return p.db.Ctx
}

func (p *Partners) Create(treatment prescription.SyphilisTreatment) error {
	ref := p.db.Client.Collection(p.collection)
	_, err := ref.Doc(treatment.ID).Set(p.ctx(), treatment)
	if err != nil {
		return fmt.Errorf("failed to create partner's syphilis treatment: %w", err)
	}
	return nil
}

func (p *Partners) FindByPregnancyId(pregnancyId int) ([]prescription.SyphilisTreatment, error) {
	ref := p.db.Client.Collection(p.collection)
	iter := ref.Where("pregnancyId", "==", pregnancyId).Documents(p.ctx())
	var ts []prescription.SyphilisTreatment
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to fetch pregnancy's contact treatments: %w", err)
		}
		var t prescription.SyphilisTreatment
		if err := doc.DataTo(&t); err != nil {
			return nil, fmt.Errorf("failed to transform partner's syphilis treatment data: %w", err)
		}
		ts = append(ts, t)
	}
	return ts, nil
}

func (p *Partners) Update(treatment prescription.SyphilisTreatment) error {
	ref := p.db.Client.Collection(p.collection).Doc(treatment.ID)
	_, err := ref.Update(p.ctx(), []firestore.Update{
		{
			Path:  "medication",
			Value: treatment.Medication,
		},
		{
			Path:  "dosage",
			Value: treatment.Dosage,
		},
		{
			Path:  "comments",
			Value: treatment.Comments,
		},
		{
			Path:  "date",
			Value: treatment.Date,
		},
		{
			Path:  "updatedBy",
			Value: treatment.UpdatedBy,
		},
		{
			Path:  "updatedAt",
			Value: treatment.UpdatedAt,
		},
	})
	return err
}
