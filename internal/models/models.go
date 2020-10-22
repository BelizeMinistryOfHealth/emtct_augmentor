package models

import (
	"database/sql"
	"time"
)

type Patient struct {
	Id               string    `json:"patientId"`
	FirstName        string    `json:"firstName"`
	MiddleName       string    `json:"middleName"`
	LastName         string    `json:"lastName"`
	Dob              time.Time `json:"dob"`
	Ssn              string    `json:"ssn"`
	CountryOfBirth   string    `json:"countryOfBirth"`
	DistrictAddress  string    `json:"district"`
	CommunityAddress string    `json:"community"`
	Education        string    `json:"education"`
	Ethnicity        string    `json:"ethnicity"`
	Hiv              bool      `json:"hiv"`
	NextOfKin        string    `json:"nextOfKin"`
	NextOfKinPhone   string    `json:"nextOfKinPhone"`
}

type ObstetricHistory struct {
	Id             string    `json:"id"`
	PatientId      string    `json:"patientId"`
	Date           time.Time `json:"date"`
	ObstetricEvent string    `json:"event"`
}

type Diagnosis struct {
	Id        string    `json:"id"`
	PatientId string    `json:"patientId"`
	Date      time.Time `json:"date"`
	Name      string    `json:"name"`
}

type PregnancyVitals struct {
	Id                   int            `json:"id"`
	PatientId            int            `json:"patientId"`
	GestationalAge       int            `json:"gestationalAge"`
	Para                 int            `json:"para"`
	Cs                   bool           `json:"cs"`
	AbortiveOutcome      sql.NullString `json:"abortiveOutcome"`
	DiagnosisDate        time.Time      `json:"diagnosisDate"`
	Planned              bool           `json:"planned"`
	AgeAtLmp             int            `json:"ageAtLmp"`
	Lmp                  time.Time      `json:"lmp"`
	Edd                  time.Time      `json:"edd"`
	DateOfBooking        time.Time      `json:"dateOfBooking"`
	PrenatalCareProvider string         `json:"prenatalCareProvider"`
	TotalChecks          int            `json:"totalChecks"`
}

func FindCurrentPregnancy(ps []PregnancyVitals) *PregnancyVitals {
	for _, v := range ps {
		if v.Edd.After(time.Now()) {
			return &v
		}
	}
	return nil
}
