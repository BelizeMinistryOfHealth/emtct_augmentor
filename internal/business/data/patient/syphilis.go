package patient

import (
	"database/sql"
	"fmt"
	"time"

	"moh.gov.bz/mch/emtct/internal/business/data/prescription"
)

func (p *Patients) FindSyphilisTreatment(patientId int, beginDate *time.Time, endDate *time.Time) ([]prescription.Prescription, error) {
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
`
	args := []interface{}{patientId}
	if beginDate != nil && endDate != nil {
		stmt = fmt.Sprintf("%s AND adep.prescribed_time BETWEEN $2 AND $3", stmt)
		args = append(args, beginDate.Format(layoutISO))
		args = append(args, endDate.Format(layoutISO))
	}
	stmt = fmt.Sprintf("%s ORDER BY adep.prescribed_time DESC", stmt)
	rows, err := p.acsis.Query(stmt, args...)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("error retrieving syphilis from acsis: %+v", err)
	}

	var prescriptions []prescription.Prescription
	var totalDoses sql.NullInt64
	for rows.Next() {
		var prescription prescription.Prescription
		err := rows.Scan(&prescription.Id,
			&totalDoses,
			&prescription.Pharmaceutical,
			&prescription.Frequency,
			&prescription.Strength,
			&prescription.Comments,
			&prescription.PrescribedTime)
		if err != nil {
			return prescriptions, fmt.Errorf("error scanning syphillis prescription from acsis: %+v", err)
		}
		prescription.PatientId = patientId
		if totalDoses.Valid {
			prescription.TotalDoses = int(totalDoses.Int64)
		}

		prescriptions = append(prescriptions, prescription)
	}

	return prescriptions, nil

}
