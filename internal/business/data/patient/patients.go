package patient

import (
	"database/sql"
	"fmt"
)

func (p *Patients) FindBasicInfo(patientId int) (*BasicInfo, error) {
	stmt := `
		SELECT
		    hp.patient_id,
			p.first_name,
			p.last_name,
			p.middle_name,
			hp.birth_date,
			hp.ssi_number
		FROM acsis_people p
			INNER JOIN acsis_hc_patients hp ON p.person_id=hp.person_id
		WHERE hp.patient_id=$1;
`
	row := p.Acsis.QueryRow(stmt, patientId)
	var info BasicInfo
	err := row.Scan(&info.Id,
		&info.FirstName,
		&info.LastName,
		&info.MiddleName,
		&info.Dob,
		&info.Ssn)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &info, nil
	default:
		return nil, fmt.Errorf("error retrieving patient basic info from acsis: %+v", err)

	}
}

// FindByPatientId searches for a patient who is currently pregnant.
// A patient is considered pregnant if she has a record in the acsis_hc_pregnancies
func (p *Patients) FindByPatientId(id int) (*Patient, error) {

	stmt := `SELECT p.patient_id, ahp.pregnancy_id,
		l.first_name, l.last_name, l.middle_name,
		p.birth_date, p.ssi_number,p.birth_place, concat(l2.first_name, ' ', l2.last_name) as next_of_kin,
       	ac2.phone1 as next_of_kin_phone,
		ae.name as ethnicity, ahsl.name as education,
		concat(ac.address1, ' ', ac.address2, ',', am.name, ',', aterr.name) as address
		FROM acsis_hc_patients as p
		INNER JOIN acsis_people as l on p.person_id = l.person_id
		INNER JOIN acsis_adt_encounters AS e ON e.patient_id = p.patient_id
		INNER JOIN acsis_adt_encounter_diagnoses ed ON e.encounter_id=ed.encounter_id
		INNER JOIN acsis_contacts ac on l.contact_id = ac.contact_id
		INNER JOIN acsis_municipalities am on ac.municipality_id = am.municipality_id
		INNER JOIN acsis_territories aterr ON ac.territory_id = aterr.territory_id
		INNER JOIN acsis_hc_pregnancies ahp on p.patient_id = ahp.patient_id AND ahp.active IS TRUE
		LEFT JOIN acsis_adt_next_of_kins aanok on p.next_of_kin_id = aanok.next_of_kin_id
		LEFT JOIN acsis_people l2 on aanok.person_id = l2.person_id
		LEFT JOIN acsis_contacts ac2 ON l2.contact_id = ac2.contact_id
		LEFT JOIN acsis_ethnicities ae on p.ethnicity_id = ae.ethnicity_id
		LEFT JOIN acsis_hc_schooling_levels ahsl on p.schooling_level_id = ahsl.schooling_level_id
-- 		WHERE ed.disease_id IN (473, 474, 475, 476, 477, 9921, 32590, 33195) -- the HIV Test
		WHERE p.patient_id = $1
	    ORDER BY ed.diagnosis_time DESC
		LIMIT 1;`

	var patient Patient
	row := p.Acsis.QueryRow(stmt, id)
	var nok sql.NullString
	var nokPhone sql.NullString
	err := row.Scan(&patient.Id,
		&patient.PregnancyId,
		&patient.FirstName,
		&patient.LastName,
		&patient.MiddleName,
		&patient.Dob,
		&patient.Ssn,
		&patient.CountryOfBirth,
		&nok,
		&nokPhone,
		&patient.Ethnicity,
		&patient.Education,
		&patient.Address)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		if nok.Valid {
			patient.NextOfKin = nok.String
		}
		if nokPhone.Valid {
			patient.NextOfKinPhone = nok.String
		}

		return &patient, nil
	default:
		return nil, fmt.Errorf("error querying acsis db for patient: %+v", err)
	}

}
