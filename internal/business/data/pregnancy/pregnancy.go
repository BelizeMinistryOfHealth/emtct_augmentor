package pregnancy

import (
	"database/sql"
	"fmt"

	"moh.gov.bz/mch/emtct/internal/db"
)

type Pregnancies struct {
	*db.EmtctDb
}

func (p Pregnancies) FindLatest(patientId int) (*Pregnancy, error) {
	stmt := `
	SELECT 
	       pregnancy_id, patient_id, lmp, edd, end_time
	FROM pregnancies
	WHERE 
	      patient_id = $1
	ORDER BY lmp DESC
	LIMIT 1;
`
	row := p.QueryRow(stmt, patientId)
	var pregnancy Pregnancy
	err := row.Scan(
		&pregnancy.PregnancyId,
		&pregnancy.PatientId,
		&pregnancy.Lmp,
		&pregnancy.Edd,
		&pregnancy.EndTime)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &pregnancy, nil
	default:
		return nil, fmt.Errorf("error retrieving pregnancy from emtctdb: %w", err)
	}
}
