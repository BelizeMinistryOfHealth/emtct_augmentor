package contactTracing

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"google.golang.org/api/iterator"
)

func (d *ContactTracings) Create(c ContactTracing) error {
	ref := d.db.Client.Collection(d.collection)
	_, err := ref.Doc(c.ID).Set(d.ctx(), c)
	if err != nil {
		return fmt.Errorf("failed to create contact tracing for patient: %w", err)
	}
	return nil
}

func (d *ContactTracings) FindByPregnancyId(pregnancyId int) ([]ContactTracing, error) {
	ref := d.db.Client.Collection(d.collection)
	iter := ref.Where("pregnancyId", "==", pregnancyId).Documents(d.ctx())
	var ts []ContactTracing
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve pregnancy's contact tracing: %w", err)
		}
		var t ContactTracing
		if err := doc.DataTo(&t); err != nil {
			return nil, fmt.Errorf("failed to transform contact tracing data: %w", err)
		}
		ts = append(ts, t)
	}
	return ts, nil
}

func (d *ContactTracings) Update(c ContactTracing) error {
	ref := d.db.Client.Collection(d.collection).Doc(c.ID)
	_, err := ref.Update(d.ctx(), []firestore.Update{
		{
			Path:  "test",
			Value: c.Test,
		},
		{
			Path:  "testResult",
			Value: c.TestResult,
		},
		{
			Path:  "comments",
			Value: c.Comments,
		},
		{
			Path:  "date",
			Value: c.Date,
		},
		{
			Path:  "updatedBy",
			Value: c.UpdatedBy,
		},
		{
			Path:  "updatedAt",
			Value: c.UpdatedAt,
		},
	})
	return err
}
