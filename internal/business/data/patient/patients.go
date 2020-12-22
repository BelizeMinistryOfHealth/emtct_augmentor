package patient

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"google.golang.org/api/iterator"
	"moh.gov.bz/mch/emtct/internal/models"
	"moh.gov.bz/mch/emtct/nums"
	"strconv"
	"time"
)

// FindByPatientId searches for a patient who is currently pregnant.
// A patient is considered pregnant if she has a record in the acsis_hc_pregnancies
func (p *Patients) FindByPatientId(id string) (*models.Patient, error) {

	ref := p.firestore.Client.Collection(p.collections.Patient)
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
	coll := p.firestore.Client.Collection(p.collections.Patient)
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

func (p *Patients) GetPregnancies(patientId string) ([]models.Pregnancy, error) {
	ref := p.firestore.Client.Collection(p.collections.Pregnancies)
	iter := ref.Where("patientId", "==", patientId).OrderBy("lmp", firestore.Desc).Documents(p.ctx())
	var pregs []models.Pregnancy
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var pr models.Pregnancy
		if err := doc.DataTo(&pr); err != nil {
			return nil, fmt.Errorf("error converting pregnancy result: %w", err)
		}
		pregs = append(pregs, pr)
	}
	return pregs, nil
}

func (p *Patients) GetPregnancy(pregnancyId int) (*models.Pregnancy, error) {
	ref := p.firestore.Client.Collection(p.collections.Pregnancies)
	snap, err := ref.Doc(strconv.Itoa(pregnancyId)).Get(p.ctx())
	if err != nil {
		return nil, fmt.Errorf("error retrieving pregnancy: %w", err)
	}
	var preg models.Pregnancy
	if err := snap.DataTo(&preg); err != nil {
		return nil, fmt.Errorf("error converting pregnancy: %w", err)
	}
	return &preg, nil
}

func (p *Patients) GetDiagnoses(patientId string) ([]models.Diagnosis, error) {
	ref := p.firestore.Client.Collection(p.collections.Diagnoses)
	iter := ref.Where("patientId", "==", patientId).OrderBy("date", firestore.Desc).Documents(p.ctx())
	var ds []models.Diagnosis
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var d models.Diagnosis
		if err := doc.DataTo(&d); err != nil {
			return nil, fmt.Errorf("error converting diagnosis: %w", err)
		}
		ds = append(ds, d)
	}
	return ds, nil
}

func (p *Patients) GetInfantForPregnancy(patientId int, lmp time.Time) (*models.Infant, error) {
	ref := p.firestore.Client.Collection(p.collections.Infants)
	iter := ref.Where("mother.patientId", "==", patientId).
		Where("dob", ">", lmp).
		Where("dob", "<", lmp.Add(time.Hour*25*365)).
		Limit(1).
		Documents(p.ctx())
	var preg models.Infant
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve infant for pregnancy: %w", err)
		}
		if err := doc.DataTo(&preg); err != nil {
			return nil, fmt.Errorf("failed to transform infant: %w", err)
		}
	}
	return &preg, nil
}

func (p *Patients) GetInfant(infantId string) (*models.Infant, error) {
	ref := p.firestore.Client.Collection(p.collections.Infants)
	snap, err := ref.Doc(infantId).Get(p.ctx())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve infant: %w", err)
	}
	var inf models.Infant
	if err := snap.DataTo(&inf); err != nil {
		return nil, fmt.Errorf("failed to convert infant data: %w", err)
	}
	return &inf, nil
}
