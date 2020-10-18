package models

import "time"

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
