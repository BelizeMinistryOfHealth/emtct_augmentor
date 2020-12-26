package infant

import (
	"fmt"
	"google.golang.org/api/iterator"
	"time"
)

func (d *Infants) SaveHivScreening(v HivScreening) error {
	_, err := d.firestore.Client.Collection(d.collections.HivScreening).Doc(v.Id).Set(d.ctx(), v)
	if err != nil {
		return fmt.Errorf("failed to create new hiv screening in firestore: %w", err)
	}
	return nil
}

func (d *Infants) FindHivScreeningsByPatient(patientId string) ([]HivScreening, error) {

	ref := d.firestore.Client.Collection(d.collections.HivScreening)
	iter := ref.Where("patientId", "==", patientId).Documents(d.ctx())
	var screenings []HivScreening
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var s HivScreening
		if err := doc.DataTo(&s); err != nil {
			return nil, fmt.Errorf("failed to convert hiv screening: %w", err)
		}
		screenings = append(screenings, s)
	}

	return screenings, nil
}

func (d *Infants) FindHivScreeningById(id string) (*HivScreening, error) {
	ref := d.firestore.Client.Collection(d.collections.HivScreening)
	snap, err := ref.Doc(id).Get(d.ctx())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve hiv screening: %w", err)
	}
	var s HivScreening
	if err := snap.DataTo(&s); err != nil {
		return nil, fmt.Errorf("failed to convert hiv screening data: %w", err)
	}
	return &s, nil
}

// IsHivScreeningTimely indicates if an hiv screening was done in a timely manner.
// The timeliness depends on the type of test and when the sample was taken:
// PCR 1: sample must be taken 3 days or less after birth.
// PCR 2: sample must be taken no later than 6 weeks after birth
// PCR 3: sample must be taken no later than 90 days after birth
// ELISA: sample must be taken no longer than 18 months after birth
func IsHivScreeningTimely(birthDate time.Time, testName string, dateSampleTaken time.Time) bool {
	diff := dateSampleTaken.Sub(birthDate).Hours() / 24
	switch testName {
	case "PCR 1":
		return diff < 4
	case "PCR 2":
		return diff < (6 * 7)
	case "PCR 3":
		return diff < 91
	case "ELISA":
		return diff <= (18 * 7 * 4)
	default:
		return false
	}
}

// HivScreeningDueDate calculates the due date for taking a sample for an HIV screening.
// PCR 1: sample must be taken 3 days or less after birth.
// PCR 2: sample must be taken no later than 6 weeks after birth
// PCR 3: sample must be taken no later than 90 days after birth
// ELISA: sample must be taken no longer than 18 months after birth
func HivScreeningDueDate(testName string, birthDate time.Time) time.Time {
	switch testName {
	case "PCR 1":
		return birthDate.AddDate(0, 0, 3)
	case "PCR 2":
		return birthDate.AddDate(0, 0, 42)
	case "PCR 3":
		return birthDate.AddDate(0, 0, 90)
	default:
		return birthDate.AddDate(0, 18, 0)
	}
}
