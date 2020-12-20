package patient

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"google.golang.org/api/iterator"
	"moh.gov.bz/mch/emtct/internal/models"
	"time"
)

func (p *Patients) FindLabTestsInPeriod(patientId string, startDate, endDate time.Time) ([]models.LabResult, error) {
	colRef := p.firestore.Client.Collection(p.collections.LabTests)
	iter := colRef.Where("patientId", "==", patientId).
		Where("dateOrderReceivedByLab", ">=", startDate).
		Where("dateOrderReceivedByLab", "<=", endDate).
		OrderBy("dateOrderReceivedByLab", firestore.Desc).
		Documents(p.ctx())
	var labs []models.LabResult
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve lab tests: %w", err)
		}
		var l models.LabResult
		if err := doc.DataTo(&l); err != nil {
			return nil, fmt.Errorf("failed to transform lab test: %w", err)
		}
		labs = append(labs, l)
	}
	return labs, nil
}
