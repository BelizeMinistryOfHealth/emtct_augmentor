package patient

import (
	"fmt"
	"moh.gov.bz/mch/emtct/internal/models"
	"moh.gov.bz/mch/emtct/nums"
)

// FindByPatientId searches for a patient who is currently pregnant.
// A patient is considered pregnant if she has a record in the acsis_hc_pregnancies
func (p *Patients) FindByPatientId(id string) (*models.Patient, error) {

	ref := p.firestore.Client.Collection(p.collection)
	snap, err := ref.Doc(id).Get(p.ctx())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve patient data: %w", err)
	}

	var patient models.Patient
	err = snap.DataTo(&patient)
	if err != nil {
		return nil, fmt.Errorf("failed to convert firebase data to Patient")
	}
	return &patient, nil
}

// Create creates a patient in the etmct database. The patients from acsis are mirrored into this database
// to facilitate faster queries.
// SplitPatients splits a list of patients into groups of "batchSize".
// Useful for batching firestore writes.
func SplitPatients(ps []models.Patient, batchSize int) [][]models.Patient {
	var batch [][]models.Patient
	for i := 0; i < len(ps); i += batchSize {
		bs := ps[i:nums.Min(i+batchSize, len(ps))]
		batch = append(batch, bs)
	}
	return batch
}

// Create persists patients in Firebase.
func (p *Patients) Create(patients []models.Patient) error {
	pts := SplitPatients(patients, 500)
	for _, pt := range pts {
		err := p.saveBatch(pt)
		if err != nil {
			return fmt.Errorf("error saving batch of patients in firestore: %w", err)
		}
	}
	return nil
}

func (p *Patients) saveBatch(patients []models.Patient) error {
	batch := p.firestore.Client.Batch()
	coll := p.firestore.Client.Collection("patients")
	for _, v := range patients {
		ref := coll.Doc(v.ID)
		//dsnap, _ := ref.Get(p.ctx())
		////Only create new records
		//switch dsnap.Exists() {
		//case false:
		//	batch.Set(ref, v)
		//case true:
		//	//noop
		//
		//}
		batch.Set(ref, v)

	}
	_, err := batch.Commit(p.ctx())
	return err
}
