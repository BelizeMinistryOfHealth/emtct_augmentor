package hiv

import (
	"fmt"
)

func (h *HIV) FindHivDiagnoses(patientId int) ([]Diagnosis, error) {
	stmt := `
		SELECT aaed.encounter_diagnosis_id,
       		e.patient_id,
			aai10d.name, 
			aaed.diagnosis_time 
		FROM acsis_adt_encounters AS e
			INNER JOIN acsis_adt_encounter_diagnoses aaed on e.encounter_id = aaed.encounter_id
			INNER JOIN acsis_adt_icd10_diseases aai10d on aaed.disease_id = aai10d.disease_id
		WHERE e.patient_id=$1
			AND  aaed.disease_id IN (473, 474, 475, 476, 477, 9921, 32590, 33195) -- the HIV Test
		ORDER BY aaed.diagnosis_time DESC;
`

	var diagnoses []Diagnosis
	rows, err := h.AcsisDb.Query(stmt, patientId)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("error retrieving the patient's hiv diagnoses from acsis: %+v", err)
	}
	for rows.Next() {
		var d Diagnosis
		err := rows.Scan(&d.Id, &d.PatientId, &d.Name, &d.Date)
		if err != nil {
			return nil, fmt.Errorf("error scanning patient's hiv diagnosis from acsis: %+v", err)
		}
		diagnoses = append(diagnoses, d)
	}
	return diagnoses, nil
}
