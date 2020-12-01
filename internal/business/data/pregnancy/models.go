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
