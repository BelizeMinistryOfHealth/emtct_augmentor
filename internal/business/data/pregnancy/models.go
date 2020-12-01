package pregnancy

import "time"

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
