package contraceptives

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"google.golang.org/api/iterator"
)

func (d *Contraceptives) Create(c ContraceptiveUsed) error {
	ref := d.db.Client.Collection(d.collection)
	_, err := ref.Doc(c.ID).Set(d.ctx(), c)
	if err != nil {
		return err
	}
	return nil
}

func (d *Contraceptives) Update(c ContraceptiveUsed) error {
	ref := d.db.Client.Collection(d.collection).Doc(c.ID)
	_, err := ref.Update(d.ctx(), []firestore.Update{
		{
			Path:  "contraceptive",
			Value: c.Contraceptive,
		},
		{
			Path:  "comments",
			Value: c.Comments,
		},
		{
			Path:  "dateUsed",
			Value: c.DateUsed,
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

func (d *Contraceptives) FindById(id string) (*ContraceptiveUsed, error) {
	ref := d.db.Client.Collection(d.collection)
	snap, err := ref.Doc(id).Get(d.ctx())
	if err != nil {
		return nil, err
	}
	var c ContraceptiveUsed
	if err := snap.DataTo(&c); err != nil {
		return nil, fmt.Errorf("failed to transform contraceptive data: %w", err)
	}
	return &c, nil
}

func (d *Contraceptives) FindByPregnancyId(pregnancyId int) ([]ContraceptiveUsed, error) {
	ref := d.db.Client.Collection(d.collection)
	iter := ref.Where("pregnancyId", "==", pregnancyId).Documents(d.ctx())
	var cs []ContraceptiveUsed
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to fetch pregnancy's contraceptives: %w", err)
		}
		var c ContraceptiveUsed
		if err := doc.DataTo(&c); err != nil {
			return nil, fmt.Errorf("failed to transform contraceptive data: %w", err)
		}
		cs = append(cs, c)
	}
	return cs, nil
}
