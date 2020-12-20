package patient

import (
	"context"
	"database/sql"
	"moh.gov.bz/mch/emtct/internal/db"
	"time"
)

type DbCollections struct {
	Patient       string
	Arvs          string
	Prescriptions string
	Pregnancies   string
	Diagnoses     string
}

type Patients struct {
	acsis       *sql.DB
	emtctDb     *sql.DB
	firestore   *db.FirestoreClient
	collections DbCollections
}

func New(acsisDb, emtctDb *sql.DB, firestore *db.FirestoreClient) Patients {
	collections := DbCollections{
		Patient:       "patients",
		Arvs:          "arvs",
		Prescriptions: "prescriptions",
		Pregnancies:   "pregnancies",
		Diagnoses:     "diagnoses",
	}
	return Patients{acsis: acsisDb, emtctDb: emtctDb, collections: collections, firestore: firestore}
}

func (p *Patients) ctx() context.Context {
	return p.firestore.Ctx
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

func (p *Patient) Index(vs []Patient) int {
	for i, v := range vs {
		if v.Id == p.Id {
			return i
		}
	}
	return -1
}

func (p *Patient) Include(vs []Patient) bool {
	return p.Index(vs) >= 0
}

type BasicInfo struct {
	Id         string    `json:"patientId"`
	FirstName  string    `json:"firstName"`
	MiddleName string    `json:"middleName"`
	LastName   string    `json:"lastName"`
	Dob        time.Time `json:"dob"`
	Ssn        string    `json:"ssn"`
}
