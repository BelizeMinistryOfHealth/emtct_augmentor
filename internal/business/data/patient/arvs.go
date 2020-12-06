package patient

import (
	"database/sql"
	"fmt"
	"time"

	"moh.gov.bz/mch/emtct/internal/business/data/prescription"
)

const (
	layoutISO = "2006-01-02"
)

func (p *Patients) FindArvsByPatient(patientId int, beginDate, endDate time.Time) ([]prescription.Prescription, error) {
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
		WHERE p.patient_id=$1 AND adep.prescribed_time BETWEEN $2 AND $3
			AND (aap.name ILIKE '%Lamivudine%'
			OR aap.name ILIKE '%Zidovudine%'
			OR aap.name ILIKE '%Nevirapine%')
		ORDER BY adep.prescribed_time DESC;
`
	rows, err := p.acsis.Query(stmt,
		patientId,
		beginDate.Format(layoutISO),
		endDate.Format(layoutISO))
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("error retrieving arvs from acsis: %w", err)
	}
	var arvs []prescription.Prescription
	var totalDoses sql.NullInt64
	for rows.Next() {
		var arv prescription.Prescription
		err := rows.Scan(&arv.Id,
			&totalDoses,
			&arv.Pharmaceutical,
			&arv.Frequency,
			&arv.Strength,
			&arv.Comments,
			&arv.PrescribedTime)
		if err != nil {
			return arvs, fmt.Errorf("error scanning arv prescription from acsis: %+v", err)
		}
		arv.PatientId = patientId
		if totalDoses.Valid {
			arv.TotalDoses = int(totalDoses.Int64)
		}

		arvs = append(arvs, arv)
	}

	return arvs, nil
}
