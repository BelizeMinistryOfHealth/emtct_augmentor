package homeVisits

import (
	"fmt"
	"google.golang.org/api/iterator"
)

func (h *HomeVisits) Save(v HomeVisit) error {
	ref := h.db.Client.Collection(h.collection)
	_, err := ref.Doc(v.Id).Set(h.ctx(), v)
	if err != nil {
		return fmt.Errorf("failed to create home visit: %w", err)
	}
	return nil
}

func (h *HomeVisits) FindById(id string) (*HomeVisit, error) {
	ref := h.db.Client.Collection(h.collection)
	snap, err := ref.Doc(id).Get(h.ctx())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve home visit: %w", err)
	}
	var v HomeVisit
	err = snap.DataTo(&v)
	if err != nil {
		return nil, fmt.Errorf("failed to convert home visit data: %w", err)
	}
	return &v, nil
}

func (h *HomeVisits) FindByPatientId(patientId int) ([]HomeVisit, error) {
	ref := h.db.Client.Collection(h.collection)
	iter := ref.Where("patientId", "==", patientId).Documents(h.ctx())
	var visits []HomeVisit
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to fetch patient's home visits: %w", err)
		}
		var v HomeVisit
		if err := doc.DataTo(&v); err != nil {
			return nil, fmt.Errorf("failed to transform home visit data: %w", err)
		}
		visits = append(visits, v)
	}
	return visits, nil
}

func (h *HomeVisits) FindByPregnancyId(pregnancyId int) ([]HomeVisit, error) {
	ref := h.db.Client.Collection(h.collection)
	iter := ref.Where("pregnancyId", "==", pregnancyId).Documents(h.ctx())
	var visits []HomeVisit
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to fetch patient's home visits: %w", err)
		}
		var v HomeVisit
		if err := doc.DataTo(&v); err != nil {
			return nil, fmt.Errorf("failed to transform home visit data: %w", err)
		}
		visits = append(visits, v)
	}
	return visits, nil
}
