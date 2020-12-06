package patient

import (
	"context"
	"database/sql"
	"fmt"
	"time"
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
	row := p.acsis.QueryRow(stmt, patientId)
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
	row := p.acsis.QueryRow(stmt, id)
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

func (p *Patients) FindInBhisByYear(yr int) ([]Patient, error) {

	stmt := `
	SELECT 
		p.patient_id, ahp.pregnancy_id,
		l.first_name, l.last_name, l.middle_name,
		p.birth_date, p.ssi_number,p.birth_place, concat(l2.first_name, ' ', l2.last_name) as next_of_kin,
       	ac2.phone1 as next_of_kin_phone,
		ae.name as ethnicity, ahsl.name as education,
		concat(ac.address1, ' ', ac.address2, ',', am.name, ',', aterr.name) as address,
	    am.name as community,
	       terr.name as district
	FROM acsis_hc_pregnancies ahp
	    INNER JOIN acsis_hc_patients as p on p.patient_id = ahp.patient_id AND ahp.active IS TRUE
		INNER JOIN acsis_people as l on p.person_id = l.person_id
		INNER JOIN acsis_contacts ac on l.contact_id = ac.contact_id
	    INNER JOIN acsis_territories as terr ON ac.territory_id=terr.territory_id
		INNER JOIN acsis_municipalities am on ac.municipality_id = am.municipality_id
		INNER JOIN acsis_territories aterr ON ac.territory_id = aterr.territory_id
		LEFT JOIN acsis_adt_next_of_kins aanok on p.next_of_kin_id = aanok.next_of_kin_id
		LEFT JOIN acsis_people l2 on aanok.person_id = l2.person_id
		LEFT JOIN acsis_contacts ac2 ON l2.contact_id = ac2.contact_id
		LEFT JOIN acsis_ethnicities ae on p.ethnicity_id = ae.ethnicity_id
		LEFT JOIN acsis_hc_schooling_levels ahsl on p.schooling_level_id = ahsl.schooling_level_id
	WHERE ahp.last_menstrual_period_date BETWEEN $1 AND $2;
`
	leftYear := fmt.Sprintf("%d-01-01", yr)
	rightYear := fmt.Sprintf("%d-01-01", yr+1)
	rows, err := p.acsis.Query(stmt, leftYear, rightYear)

	if err != nil {
		return nil, fmt.Errorf("error querying acsis for pregnant patients by year: %w", err)
	}
	var patients []Patient

	for rows.Next() {
		var patient Patient
		var nok sql.NullString
		var nokPhone sql.NullString
		err := rows.Scan(&patient.Id,
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
			&patient.Address,
			&patient.CommunityAddress,
			&patient.DistrictAddress)
		if err != nil {
			return nil, fmt.Errorf("error scanning pregnancy data from acsis: %w", err)
		}
		patients = append(patients, patient)
	}

	return patients, nil
}

//func (p *Patients) FindExistingPatientsByYear(yr int) ([]Patient, error) {
//	stmt := `
//	SELECT id, first_name, last_name, middle_name, dob, education, ssn, country_of_birth, district_address,
//	       community_address, education, ethnicity, hiv, next_of_kin, next_of_kin_phone, hiv, hiv_diagnosis_date
//	FROM patients
//	WHERE created_at
//`
//	rows, err := p.emtctDb.Query(stmt)
//}

// Create creates a patient in the emtct database. The patients from acsis are mirrored into this database
// to facilitate faster queries.
func (p *Patients) Create(ctx context.Context, patients []Patient) error {
	// Begin transaction
	tx, err := p.emtctDb.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start a transaction for inserting patients: %w", err)
	}

	stmt := `
	INSERT INTO patients (
	                      id, first_name, middle_name, last_name, next_of_kin_phone, next_of_kin,
	                      dob, ssn, country_of_birth, district_address, community_address,
	                      education, ethnicity, hiv, hiv_diagnosis_date, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	ON CONFLICT ON CONSTRAINT patients_pkey
	DO 
	UPDATE SET first_name=$2, middle_name=$3, last_name=$4, next_of_kin_phone=$5, next_of_kin=$6, dob=$7,
	           ssn=$8, country_of_birth=$9, district_address=$10, community_address=$11, education=$12, ethnicity=$13,
	           hiv=$14, hiv_diagnosis_date=$15;
`
	for _, patient := range patients {
		_, err := tx.ExecContext(ctx, stmt,
			patient.Id,
			patient.FirstName,
			patient.MiddleName,
			patient.LastName,
			patient.NextOfKinPhone,
			patient.NextOfKin,
			patient.Dob,
			patient.Ssn,
			patient.CountryOfBirth,
			patient.DistrictAddress,
			patient.CommunityAddress,
			patient.Education,
			patient.Ethnicity,
			patient.Hiv,
			patient.HivDiagnosisDate,
			time.Now())

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error inserting new patient into emtct db: %w", err)
		}
	}

	// Commit the transactions
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit the transaction for inserting patients: %w", err)
	}
	return nil
}
