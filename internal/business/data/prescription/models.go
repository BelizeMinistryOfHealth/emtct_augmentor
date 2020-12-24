package prescription

import "time"

// SyphilisTreatment describes the treatment given to a patient's contact.
// It is very similar to the Prescription struct, but we do not capture the person's
// name.
type SyphilisTreatment struct {
	ID          string     `json:"id" firestore:"id"`
	PatientId   int        `json:"patientId" firestore:"patientId"`
	PregnancyId int        `json:"pregnancyId" firestore:"pregnancyId"`
	Medication  string     `json:"medication" firestore:"medication"`
	Dosage      string     `json:"dosage" firestore:"dosage"`
	Comments    string     `json:"comments" firestore:"comments"`
	Date        time.Time  `json:"date" firestore:"date"`
	CreatedBy   string     `json:"createdBy" firestore:"createdBy"`
	CreatedAt   time.Time  `json:"createdAt" firestore:"createdAt"`
	UpdatedBy   string     `json:"updatedBy" firestore:"updatedBy"`
	UpdatedAt   *time.Time `json:"updatedAt" firestore:"updatedAt"`
}
