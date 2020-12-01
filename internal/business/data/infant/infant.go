package infant

import (
	"fmt"
	"time"

	"moh.gov.bz/mch/emtct/internal/business/data/pregnancy"
)

func (d *Infants) FindInfant(infantId int) (*Infant, error) {
	stmt := `
	SELECT 
	       b.patient_id,
	       ppl.first_name,
	       ppl.middle_name,
	       ppl.last_name,
	       pt.birth_date,
	       mppl.first_name as mfirst_name,
	       mppl.middle_name as mmiddle_name,
	       mppl.last_name as mlast_name,
	       mpt.birth_date as mdob,
	       mpt.patient_id as mother_id
    FROM acsis_hc_births b
	INNER JOIN acsis_hc_patients pt ON pt.patient_id=b.patient_id
	INNER JOIN acsis_people ppl ON pt.person_id = ppl.person_id
	INNER JOIN acsis_hc_patients mpt ON b.mother_id=mpt.patient_id
	INNER JOIN acsis_people mppl ON mppl.person_id=mpt.person_id
	WHERE pt.patient_id=$1
	ORDER BY pt.birth_date DESC
	LIMIT 1;
`
	var infant Infant
	row := d.Acsis.QueryRow(stmt, infantId)
	err := row.Scan(
		&infant.Infant.PatientId,
		&infant.Infant.FirstName,
		&infant.Infant.MiddleName,
		&infant.Infant.LastName,
		&infant.Infant.Dob,
		&infant.Mother.FirstName,
		&infant.Mother.MiddleName,
		&infant.Mother.LastName,
		&infant.Mother.Dob,
		&infant.Mother.PatientId)
	if err != nil {
		return nil, fmt.Errorf("error querying infant basic information from acsis: %+v", err)
	}

	return &infant, nil
}

// FindPregnancyInfant returns the infant that was born after the mother's current LMP.
func (d *Infants) FindPregnancyInfant(pregnancy pregnancy.Pregnancy) (*Infant, error) {
	// Find pregnancy that corresponds to this id
	stmt := `
	SELECT 
	       b.patient_id,
	       ppl.first_name,
	       ppl.middle_name,
	       ppl.last_name,
	       pt.birth_date,
	       mppl.first_name as mfirst_name,
	       mppl.middle_name as mmiddle_name,
	       mppl.last_name as mlast_name,
	       mpt.birth_date as mdob
    FROM acsis_hc_births b
	INNER JOIN acsis_hc_patients pt ON pt.patient_id=b.patient_id
	INNER JOIN acsis_people ppl ON pt.person_id = ppl.person_id
	INNER JOIN acsis_hc_patients mpt ON b.mother_id=mpt.patient_id
	INNER JOIN acsis_people mppl ON mppl.person_id=mpt.person_id
	WHERE b.mother_id=$1 AND pt.birth_date BETWEEN $2 AND $3
	ORDER BY pt.birth_date DESC
	LIMIT 1;
`
	var infant Infant
	rightYear := pregnancy.Lmp.Add(time.Hour * 24 * 7 * 54)
	row := d.Acsis.QueryRow(stmt,
		pregnancy.PatientId,
		pregnancy.Lmp.Format("2006-01-02"),
		rightYear.Format("2006-01-02"))
	err := row.Scan(
		&infant.Infant.PatientId,
		&infant.Infant.FirstName,
		&infant.Infant.MiddleName,
		&infant.Infant.LastName,
		&infant.Infant.Dob,
		&infant.Mother.FirstName,
		&infant.Mother.MiddleName,
		&infant.Mother.LastName,
		&infant.Mother.Dob)
	if err != nil {
		return nil, fmt.Errorf("error querying infant basic information from acsis: %+v", err)
	}
	infant.Mother.PatientId = pregnancy.PatientId

	return &infant, nil
}
