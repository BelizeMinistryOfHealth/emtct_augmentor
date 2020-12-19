package patient

import (
	"fmt"
	"google.golang.org/api/iterator"
	"moh.gov.bz/mch/emtct/internal/models"
	"time"
)

const (
	layoutISO = "2006-01-02"
)

func (p *Patients) FindArvsByPatient(patientId string, beginDate, endDate time.Time) ([]models.Prescription, error) {
	colRef := p.firestore.Client.Collection(p.collections.Arvs)
	iter := colRef.
		Where("patientId", "==", patientId).
		Where("prescribedTime", ">=", beginDate).
		Where("prescribedTime", "<=", endDate).
		Documents(p.ctx())
	var rxs []models.Prescription
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve arvs: %w", err)
		}
		var p models.Prescription
		if err := doc.DataTo(&p); err != nil {
			return nil, fmt.Errorf("failed to convert arvs data from firebase: %w", err)
		}
		rxs = append(rxs, p)
	}

	return rxs, nil
}
