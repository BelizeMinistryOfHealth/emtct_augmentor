package infant

import (
	"database/sql"
	"fmt"

	"moh.gov.bz/mch/emtct/internal/models"
)

func (d *Infants) FindInfantSyphilisTreatment(patientId int) ([]models.Prescription, error) {
	stmt := `
SELECT
		    adep.encounter_pharmaceutical_id,
			adep.total_doses,
		   	aap.name as prescription,
		   	acfu.name as frequency, 
			aap.strength || ' ' || aapu.name as strength,
		   	adep.prescribing_physician_special_instructions || ' ' || adep.notes AS comments,
		   	adep.prescribed_time
		FROM acsis_hc_patients p
			INNER JOIN acsis_adt_encounters e ON p.patient_id = e.patient_id
			INNER JOIN acsis_adt_encounter_pharmaceuticals adep ON adep.encounter_id=e.encounter_id
			INNER JOIN acsis_adt_pharmaceuticals aap ON adep.pharmaceutical_id=aap.pharmaceutical_id
			INNER JOIN acsis_coe_frequency_units acfu ON acfu.frequency_unit_id =adep.frequency_unit_id
			INNER JOIN acsis_adt_pharmaceutical_units aapu ON aapu.pharmaceutical_unit_id=aap.strength_unit_id
		WHERE p.patient_id=$1 
		  AND aap.pharmaceutical_id=510
		ORDER BY adep.prescribed_time DESC;
`
	rows, err := d.Acsis.Query(stmt, patientId)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("error retrieving syphilis treatment for infant from acsis: %+v", err)
	}
	var prescriptions []models.Prescription
	for rows.Next() {
		var p models.Prescription
		var totalDoses sql.NullInt64
		err := rows.Scan(&p.Id, &totalDoses, &p.Pharmaceutical, &p.Frequency, &p.Strength, &p.Comments, &p.PrescribedTime)
		if err != nil {
			return nil, fmt.Errorf("error scanning syphilis prescriptions for infant from acsis: %+v", err)
		}
		p.PatientId = patientId
		if totalDoses.Valid {
			p.TotalDoses = int(totalDoses.Int64)
		}
		prescriptions = append(prescriptions, p)
	}
	return prescriptions, nil
}
