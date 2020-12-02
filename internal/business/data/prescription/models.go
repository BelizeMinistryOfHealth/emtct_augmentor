package prescription

import "time"

type Prescription struct {
	Id             int       `json:"id"`
	PatientId      int       `json:"patientId"`
	TotalDoses     int       `json:"totalDoses"`
	Pharmaceutical string    `json:"pharmaceutical"`
	Frequency      string    `json:"frequency"`
	Strength       string    `json:"strength"`
	Comments       string    `json:"comments"`
	PrescribedTime time.Time `json:"prescribedTime"`
}

// SyphilisTreatment describes the treatment given to a patient's contact.
// It is very similar to the Prescription struct, but we do not capture the person's
// name.
type SyphilisTreatment struct {
	Id         string     `json:"id"`
	PatientId  int        `json:"patientId"`
	Medication string     `json:"medication"`
	Dosage     string     `json:"dosage"`
	Comments   string     `json:"comments"`
	Date       time.Time  `json:"date"`
	CreatedBy  string     `json:"createdBy"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedBy  string     `json:"updatedBy"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}
