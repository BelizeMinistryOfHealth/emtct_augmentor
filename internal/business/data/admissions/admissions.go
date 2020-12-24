package admissions

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"google.golang.org/api/iterator"
)

func (a *Admissions) FindByPatientId(patientId int) ([]HospitalAdmission, error) {
	ref := a.db.Client.Collection(a.collection)
	iter := ref.Where("patientId", "==", patientId).Documents(a.ctx())
	var admissions []HospitalAdmission
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve patient's admissions: %w", err)
		}
		var ad HospitalAdmission
		if err := doc.DataTo(&ad); err != nil {
			return nil, fmt.Errorf("failed to transform admission: %w", err)
		}
		admissions = append(admissions, ad)
	}
	return admissions, nil
}

func (a *Admissions) FindById(id string) (*HospitalAdmission, error) {
	ref := a.db.Client.Collection(a.collection)
	snap, err := ref.Doc(id).Get(a.ctx())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve admission: %w", err)
	}
	var h HospitalAdmission
	if err := snap.DataTo(&h); err != nil {
		return nil, fmt.Errorf("failed to transform admission: %w", err)
	}
	return &h, nil
}

func (a *Admissions) Save(h HospitalAdmission) error {
	ref := a.db.Client.Collection(a.collection)
	_, err := ref.Doc(h.ID).Set(a.ctx(), h)
	if err != nil {
		return fmt.Errorf("failed to create admission: %w", err)
	}
	return nil
}

func (a *Admissions) Update(h HospitalAdmission) error {
	ref := a.db.Client.Collection(a.collection).Doc(h.ID)
	_, err := ref.Update(a.ctx(), []firestore.Update{
		{
			Path:  "facility",
			Value: h.Facility,
		},
		{
			Path:  "reason",
			Value: h.Reason,
		},
		{
			Path:  "dateAdmitted",
			Value: h.DateAdmitted,
		},
		{
			Path:  "updatedBy",
			Value: h.UpdatedBy,
		},
		{
			Path:  "updatedAt",
			Value: h.UpdatedAt,
		},
	})
	return err
}

func (a *Admissions) FindByPregnancyId(id int) ([]HospitalAdmission, error) {
	ref := a.db.Client.Collection(a.collection)
	iter := ref.Where("pregnancyId", "==", id).Documents(a.ctx())
	var admissions []HospitalAdmission
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve pregnancy's admissions: %w", err)
		}
		var h HospitalAdmission
		if err := doc.DataTo(&h); err != nil {
			return nil, fmt.Errorf("failed to transform admission: %w", err)
		}
		admissions = append(admissions, h)
	}
	return admissions, nil
}
