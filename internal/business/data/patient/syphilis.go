package patient

import (
	"fmt"
	"google.golang.org/api/iterator"
	"moh.gov.bz/mch/emtct/internal/models"
	"time"
)

func (p *Patients) FindSyphilisTreatment(patientId string, beginDate *time.Time, endDate *time.Time) ([]models.Prescription, error) {
	return p.FindRx(patientId, 510, beginDate, endDate)
}

func (p *Patients) FindRx(patientId string, rx int, beginDate *time.Time, endDate *time.Time) ([]models.Prescription, error) {
	colRef := p.firestore.Client.Collection(p.collections.Prescriptions)
	iter := colRef.Where("id", "==", rx).Where("patientId", "==", patientId).Documents(p.ctx())
	var rxs []models.Prescription
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve rxs: %w", err)
		}
		var p models.Prescription
		if err := doc.DataTo(&p); err != nil {
			return nil, fmt.Errorf("failed to convert rx data from firebase: %w", err)
		}
		if beginDate == nil && endDate == nil {
			rxs = append(rxs, p)
		} else {
			if p.PrescribedTime.After(*beginDate) && p.PrescribedTime.Before(*endDate) {
				rxs = append(rxs, p)
			}
		}
	}

	return rxs, nil
}
