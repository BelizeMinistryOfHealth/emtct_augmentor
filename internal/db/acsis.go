package db

import (
	"database/sql"
	"fmt"
	"time"

	"moh.gov.bz/mch/emtct/internal/config"
	"moh.gov.bz/mch/emtct/internal/models"
)

type AcsisDb struct {
	*sql.DB
}

const (
	layoutISO = "2006-01-02"
)

func NewAcsisConnection(cnf *config.DbConf) (*AcsisDb, error) {
	connstr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", cnf.Username, cnf.Database, cnf.Password, cnf.Host)
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}
	return &AcsisDb{db}, nil
}

// FindByPatientId searches for a patient who is currently pregnant and is HIV+.
// A patient is considered pregnant if she has a record in the acsis_hc_pregnancies
func (d *AcsisDb) FindByPatientId(id int) (*models.Patient, error) {
	// we only want patients who are currently pregnant or how gave birth no later than 18 months ago.
	startDate := time.Now().AddDate(-1, -6, 0)
	dateLimit := startDate.Format(layoutISO)
	stmt := `SELECT p.patient_id, ahp.pregnancy_id, altri.released_time,
		l.first_name, l.last_name, l.middle_name,
		p.birth_date, p.ssi_number,p.birth_place, concat(l2.first_name, ' ', l2.last_name) as next_of_kin,
       	ac2.phone1 as next_of_kin_phone,
		ae.name as ethnicity, ahsl.name as education,
		concat(ac.address1, ' ', ac.address2, ',', am.name, ',', aterr.name) as address
		FROM acsis_hc_patients as p
		INNER JOIN acsis_people as l on p.person_id = l.person_id
		INNER JOIN acsis_adt_encounters AS e ON e.patient_id = p.patient_id
		INNER JOIN acsis_lab_test_requests altr on e.encounter_id = altr.encounter_id
		INNER JOIN acsis_lab_test_request_items altri on altr.test_request_id = altri.test_request_id
		INNER JOIN acsis_lab_tests alt on altri.test_id = alt.test_id
		INNER JOIN acsis_lab_test_request_results_collected altrrc on altri.test_request_item_id = altrrc.test_request_item_id
		INNER JOIN acsis_lab_test_results a on altrrc.test_result_id = a.test_id AND a.test_id IN (2, 152, 5015, 5033, 5032)
		INNER JOIN acsis_lab_user_defined_list_items aludli on altrrc.user_defined_list_value = aludli.user_defined_list_item_id
		INNER JOIN acsis_contacts ac on l.contact_id = ac.contact_id
		INNER JOIN acsis_municipalities am on ac.municipality_id = am.municipality_id
		INNER JOIN acsis_territories aterr ON ac.territory_id = aterr.territory_id
		INNER JOIN acsis_hc_pregnancies ahp on p.patient_id = ahp.patient_id AND ahp.active IS TRUE
		LEFT JOIN acsis_adt_next_of_kins aanok on p.next_of_kin_id = aanok.next_of_kin_id
		LEFT JOIN acsis_people l2 on aanok.person_id = l2.person_id
		LEFT JOIN acsis_contacts ac2 ON l2.contact_id = ac2.contact_id
		LEFT JOIN acsis_ethnicities ae on p.ethnicity_id = ae.ethnicity_id
		LEFT JOIN acsis_hc_schooling_levels ahsl on p.schooling_level_id = ahsl.schooling_level_id
		WHERE altri.test_id IN (2) -- the HIV Test
		AND altrrc.user_defined_list_value IS NOT NULL
		AND ahp.last_menstrual_period_date > $2
		AND a.test_result_id = 348
		AND p.patient_id = $1
		ORDER BY released_time DESC LIMIT 1;`

	var patient models.Patient
	row := d.DB.QueryRow(stmt, id, dateLimit)
	err := row.Scan(&patient.Id,
		&patient.PregnancyId,
		&patient.HivDiagnosisDate,
		&patient.FirstName,
		&patient.LastName,
		&patient.MiddleName,
		&patient.Dob,
		&patient.Ssn,
		&patient.CountryOfBirth,
		&patient.NextOfKin,
		&patient.NextOfKinPhone,
		&patient.Ethnicity,
		&patient.Education,
		&patient.Address)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		patient.Hiv = true
		return &patient, nil
	default:
		return nil, fmt.Errorf("error querying acsis db for patient: %+v", err)
	}

}

func (d *AcsisDb) FindDiagnosesBeforePregnancy(patientId int) ([]models.Diagnosis, error) {
	stmt := `SELECT aaed.encounter_diagnosis_id,
       		e.patient_id,
			aai10d.name, 
			aaed.diagnosis_time 
		FROM acsis_adt_encounters AS e
		INNER JOIN acsis_adt_encounter_diagnoses aaed on e.encounter_id = aaed.encounter_id
		INNER JOIN acsis_adt_icd10_diseases aai10d on aaed.disease_id = aai10d.disease_id
		INNER JOIN acsis_hc_pregnancies ahp on e.patient_id = ahp.patient_id
		WHERE e.patient_id=$1
		AND aaed.diagnosis_time < ahp.last_menstrual_period_date
		ORDER BY aaed.diagnosis_time DESC`
	var diagnoses []models.Diagnosis
	rows, err := d.Query(stmt, patientId)
	if err != nil {
		return nil, fmt.Errorf("error querying diagnoses before pregnancy from acsis: %+v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var diagnosis models.Diagnosis
		err := rows.Scan(
			&diagnosis.Id,
			&diagnosis.PatientId,
			&diagnosis.Name,
			&diagnosis.Date)
		if err != nil {
			return nil, fmt.Errorf("error scanning diagnosis for patient(%d) %+v", patientId, err)
		}
		diagnoses = append(diagnoses, diagnosis)
	}
	return diagnoses, nil
}

func (d *AcsisDb) FindDiagnosesDuringPregnancy(patientId int) ([]models.Diagnosis, error) {
	stmt := `SELECT
			aaed.encounter_diagnosis_id,
       		e.patient_id,
       		aai10d.name,
       		aaed.diagnosis_time
		FROM acsis_adt_encounters AS e
		INNER JOIN acsis_adt_encounter_diagnoses aaed on e.encounter_id = aaed.encounter_id
		INNER JOIN acsis_adt_icd10_diseases aai10d on aaed.disease_id = aai10d.disease_id
		INNER JOIN acsis_hc_pregnancies ahp on e.patient_id = ahp.patient_id
		WHERE e.patient_id=$1
		AND aaed.diagnosis_time < ahp.estimated_delivery_date
		AND aaed.diagnosis_time > ahp.last_menstrual_period_date
		ORDER BY aaed.diagnosis_time DESC`
	var diagnoses []models.Diagnosis
	rows, err := d.Query(stmt, patientId)
	if err != nil {
		return nil, fmt.Errorf("error querying diagnoses before pregnancy from acsis: %+v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var diagnosis models.Diagnosis
		err := rows.Scan(
			&diagnosis.Id,
			&diagnosis.PatientId,
			&diagnosis.Name,
			&diagnosis.Date)
		if err != nil {
			return nil, fmt.Errorf("error scanning diagnosis for patient(%d) %+v", patientId, err)
		}
		diagnoses = append(diagnoses, diagnosis)
	}
	return diagnoses, nil
}

func (d *AcsisDb) FindObstetricHistory(patientId int) ([]models.ObstetricHistory, error) {
	stmt := `SELECT
				b.birth_id,
				b.mother_id,
				ahbs.name, 
				b.last_modified_time
			FROM acsis_hc_births b
			INNER JOIN acsis_hc_birth_statuses ahbs on b.birth_status_id = ahbs.birth_status_id
			WHERE mother_id=$1`
	var obstetricHistory []models.ObstetricHistory
	rows, err := d.Query(stmt, patientId)
	if err != nil {
		return nil, fmt.Errorf("error querying acsis for obstetric history: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var history models.ObstetricHistory
		err := rows.Scan(&history.Id,
			&history.PatientId,
			&history.ObstetricEvent,
			&history.Date)
		if err != nil {
			return nil, fmt.Errorf("error scanning patient's obstetric history: %+v", err)
		}
		obstetricHistory = append(obstetricHistory, history)
	}
	return obstetricHistory, nil
}
