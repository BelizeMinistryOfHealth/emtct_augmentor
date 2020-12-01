package patient

import (
	"database/sql"
	"time"
)

type Patients struct {
	Acsis *sql.DB
}

func New(db *sql.DB) Patients {
	return Patients{db}
}

type Patient struct {
	Id               string     `json:"patientId"`
	PregnancyId      string     `json:"pregnancyId"`
	FirstName        string     `json:"firstName"`
	MiddleName       string     `json:"middleName"`
	LastName         string     `json:"lastName"`
	Dob              time.Time  `json:"dob"`
	Ssn              string     `json:"ssn"`
	CountryOfBirth   string     `json:"countryOfBirth"`
	DistrictAddress  string     `json:"district"`
	CommunityAddress string     `json:"community"`
	Address          string     `json:"address"`
	Education        string     `json:"education"`
	Ethnicity        string     `json:"ethnicity"`
	Hiv              bool       `json:"hiv"`
	HivDiagnosisDate *time.Time `json:"hivDiagnosisDate"`
	NextOfKin        string     `json:"nextOfKin"`
	NextOfKinPhone   string     `json:"nextOfKinPhone"`
}

type BasicInfo struct {
	Id         string    `json:"patientId"`
	FirstName  string    `json:"firstName"`
	MiddleName string    `json:"middleName"`
	LastName   string    `json:"lastName"`
	Dob        time.Time `json:"dob"`
	Ssn        string    `json:"ssn"`
}
