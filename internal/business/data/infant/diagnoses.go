package infant

import (
	"fmt"
)

func (d *Infants) FindInfantDiagnoses(infantId int) ([]Diagnoses, error) {
	stmt := `
		SELECT
			aed.disease_id,
			e.patient_id,
			aai10d.name as diagnosis,
			aed.notes,
			ap.first_name || ' ' || ap.last_name as doctor,
			aed.diagnosis_time
		FROM acsis_adt_encounters e
			INNER JOIN acsis_adt_encounter_diagnoses aed ON e.encounter_id=aed.encounter_id
			INNER JOIN acsis_adt_icd10_diseases aai10d on aed.disease_id = aai10d.disease_id
			INNER JOIN acsis_hr_staff hs ON aed.doctor_id=hs.staff_id
			INNER JOIN acsis_people ap on hs.person_id = ap.person_id
		WHERE e.patient_id=$1 
		ORDER BY aed.diagnosis_time DESC;
`
	rows, err := d.Acsis.Query(stmt, infantId)
	if err != nil {
		return nil, fmt.Errorf("error querying for infant diagnoses from acsis: %+v", err)
	}
	defer rows.Close()
	var diagnoses []Diagnoses
	for rows.Next() {
		var d Diagnoses
		err := rows.Scan(&d.DiagnosisId,
			&d.PatientId,
			&d.Diagnosis,
			&d.Comments,
			&d.Doctor,
			&d.Date)
		if err != nil {
			return nil, fmt.Errorf("error scanning infant diagnosis: %+v", err)
		}
		diagnoses = append(diagnoses, d)
	}

	return diagnoses, nil
}
