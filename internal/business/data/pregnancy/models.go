package pregnancy

import (
	"time"

	"moh.gov.bz/mch/emtct/internal/db"
)

type Pregnancies struct {
	EmtctDb *db.EmtctDb
	AcsisDb *db.AcsisDb
}

func New(emtctdb *db.EmtctDb, acsisdb *db.AcsisDb) Pregnancies {
	return Pregnancies{
		EmtctDb: emtctdb,
		AcsisDb: acsisdb,
	}
}

type Pregnancy struct {
	PatientId   int
	PregnancyId int
	Lmp         *time.Time
	Edd         *time.Time
	EndTime     *time.Time
}

func (p *Pregnancy) Index(vs []Pregnancy) int {
	for i, v := range vs {
		if v.PregnancyId == p.PregnancyId {
			return i
		}
	}
	return -1
}

func (p *Pregnancy) Include(vs []Pregnancy) bool {
	return p.Index(vs) >= 0
}

type Vitals struct {
	Id                   int        `json:"id"`
	PatientId            int        `json:"patientId"`
	GestationalAge       int        `json:"gestationalAge"`
	Para                 int        `json:"para"`
	Cs                   int        `json:"cs"`
	PregnancyOutcome     string     `json:"pregnancyOutcome"`
	DiagnosisDate        *time.Time `json:"diagnosisDate"`
	Planned              bool       `json:"planned"`
	AgeAtLmp             int        `json:"ageAtLmp"`
	Lmp                  *time.Time `json:"lmp"`
	Edd                  time.Time  `json:"edd"`
	DateOfBooking        *time.Time `json:"dateOfBooking"`
	PrenatalCareProvider string     `json:"prenatalCareProvider"`
	TotalChecks          int        `json:"totalChecks"`
	ApgarFirstMinute     int        `json:"apgarFirstMinute"`
	ApgarFifthMinute     int        `json:"apgarFifthMinute"`
	BirthStatus          string     `json:"birthStatus"`
}

type AntenatalEncounter struct {
	Id                    int        `json:"id"`
	PatientId             int        `json:"patientId"`
	MchEncounterDetailsId int        `json:"mchEncounterDetailsId"`
	EstimatedDeliveryDate *time.Time `json:"estimatedDeliveryDate"`
	BeginDate             *time.Time `json:"beginDate"`
	GestationalAge        int        `json:"gestationalAge"`
	NumberAntenatalVisits int        `json:"numberAntenatalVisits"`
}

type Diagnosis struct {
	Id        string    `json:"id"`
	PatientId string    `json:"patientId"`
	Date      time.Time `json:"date"`
	Name      string    `json:"name"`
}
