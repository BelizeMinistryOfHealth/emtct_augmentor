package db

import (
	"database/sql"
	"fmt"
	"time"

	"moh.gov.bz/mch/emtct/internal/config"
	"moh.gov.bz/mch/emtct/internal/models"
)

type EmtctDb struct {
	*sql.DB
}

// NewConnection creates a new database connection
func NewConnection(cnf *config.DbConf) (*EmtctDb, error) {
	connstr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", cnf.Username, cnf.Database, cnf.Password, cnf.Host)
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}

	return &EmtctDb{db}, nil
}

func (d *EmtctDb) FindPatientById(id string) (*models.Patient, error) {
	stmt := `SELECT id, first_name, middle_name, last_name, dob, ssn, country_of_birth, district_address, 
			community_address, education, ethnicity, hiv, next_of_kin, next_of_kin_phone FROM patients WHERE id = $1`
	var patient models.Patient
	row := d.DB.QueryRow(stmt, id)
	err := row.Scan(&patient.Id,
		&patient.FirstName,
		&patient.MiddleName,
		&patient.LastName,
		&patient.Dob,
		&patient.Ssn,
		&patient.CountryOfBirth,
		&patient.DistrictAddress,
		&patient.CommunityAddress,
		&patient.Education,
		&patient.Ethnicity,
		&patient.Hiv,
		&patient.NextOfKin,
		&patient.NextOfKinPhone)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &patient, nil
	default:
		return nil, fmt.Errorf("error querying for patients: %+v", err)
	}

}

func (d *EmtctDb) CreatePatient(p models.Patient) error {
	stmt := `INSERT INTO patients (id, first_name, middle_name, last_name, dob, ssn, country_of_birth, district_address, community_address, education, ethnicity, hiv, next_of_kin, next_of_kin_phone) 
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	_, err := d.DB.Exec(stmt,
		p.Id,
		p.FirstName,
		p.MiddleName,
		p.LastName,
		p.Dob,
		p.Ssn,
		p.CountryOfBirth,
		p.DistrictAddress,
		p.CommunityAddress,
		p.Education,
		p.Ethnicity,
		p.Hiv,
		p.NextOfKin,
		p.NextOfKinPhone,
	)
	if err != nil {
		return fmt.Errorf("error inserting patient: %+v", err)
	}
	return nil
}

func (d *EmtctDb) FindObstetricHistory(patientId string) ([]models.ObstetricHistory, error) {
	stmt := `SELECT id, patient_id, event_date, event_name FROM obstetric_history WHERE patient_id=$1`
	var obstetricHistory []models.ObstetricHistory
	rows, err := d.DB.Query(stmt, patientId)
	if err != nil {
		return nil, fmt.Errorf("error executing query for extracting obstetric history %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var history models.ObstetricHistory
		err := rows.Scan(
			&history.Id,
			&history.PatientId,
			&history.Date,
			&history.ObstetricEvent)
		if err != nil {
			return nil, fmt.Errorf("error scanning patient's (%s) obstetric history: %+v", patientId, err)
		}
		obstetricHistory = append(obstetricHistory, history)
	}
	return obstetricHistory, nil
}

func (d *EmtctDb) CreateObstetricHistory(h models.ObstetricHistory) error {
	stmt := `INSERT INTO obstetric_history (id, patient_id, event_date, event_name) VALUES($1, $2, $3, $4)`
	_, err := d.DB.Exec(stmt,
		h.Id,
		h.PatientId,
		h.Date,
		h.ObstetricEvent)
	if err != nil {
		return fmt.Errorf("error inserting obstetric history: %+v", err)
	}
	return nil
}

func (d *EmtctDb) FindDiagnoses(patientId string) ([]models.Diagnosis, error) {
	stmt := `SELECT id, patient_id, diagnosis_date, diagnosis_name FROM diagnoses WHERE patient_id=$1`
	var diagnoses []models.Diagnosis
	rows, err := d.DB.Query(stmt, patientId)
	if err != nil {
		return nil, fmt.Errorf("error executing query for extracting diagnoses: %+v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var diagnosis models.Diagnosis
		err := rows.Scan(
			&diagnosis.Id,
			&diagnosis.PatientId,
			&diagnosis.Date,
			&diagnosis.Name)
		if err != nil {
			return nil, fmt.Errorf("error scanning diagnosis for patient(%s) %+v", patientId, err)
		}
		diagnoses = append(diagnoses, diagnosis)
	}
	return diagnoses, nil
}

func (d *EmtctDb) CreateDiagnosis(di models.Diagnosis) error {
	stmt := `INSERT INTO diagnoses (id, patient_id, diagnosis_date, diagnosis_name) VALUES($1, $2, $3, $4)`
	_, err := d.DB.Exec(stmt,
		di.Id,
		di.PatientId,
		di.Date,
		di.Name)
	if err != nil {
		return fmt.Errorf("error inserting diagnosis: %+v", err)
	}
	return nil
}

func (d *EmtctDb) CreatePregnancy(p models.PregnancyVitals) error {
	stmt := `INSERT INTO pregnancies 
(id, patient_id, gestational_age, para, cs, abortive_outcome, diagnosis_date, planned, age_at_lmp, lmp, edd, date_of_booking, prenatal_care_provider, total_checks)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err := d.DB.Exec(stmt,
		p.Id,
		p.PatientId,
		p.GestationalAge,
		p.Para,
		p.Cs,
		p.PregnancyOutcome,
		p.DiagnosisDate,
		p.Planned,
		p.AgeAtLmp,
		p.Lmp,
		p.Edd,
		p.DateOfBooking,
		p.PrenatalCareProvider,
		p.TotalChecks)
	if err != nil {
		return fmt.Errorf("error inserting pregnancy: %+v", err)
	}
	return nil
}

func (d *EmtctDb) CreateLabResult(l models.LabResult) error {
	stmt := `INSERT INTO lab_results (id, patient_id, test_result, test_name, date_sample_taken, result_date)
VALUES($1, $2, $3, $4, $5, $6)`
	_, err := d.DB.Exec(stmt,
		l.Id,
		l.PatientId,
		l.TestResult,
		l.TestName,
		l.DateSampleTaken,
		l.ResultDate)
	if err != nil {
		return fmt.Errorf("error inserting lab result: %+v", err)
	}
	return nil
}

func (d *EmtctDb) SaveLabTest(l models.LabTest) error {
	stmt := `
	INSERT INTO lab_tests (id, patient_id, test_result, test_name, test_request_id,
	                       test_request_item_id, date_sample_taken, result_date, 
	                       released_time, date_order_received_by_lab)
    VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)  
`
	_, err := d.Exec(stmt,
		l.Id,
		l.PatientId,
		l.TestResult,
		l.TestName,
		l.TestRequestId,
		l.TestRequestItemId,
		l.DateSampleTaken,
		l.ResultDate,
		l.ReleasedTime,
		l.DateOrderReceivedByLab)
	if err != nil {
		return fmt.Errorf("error inserting new lab test into emtct: %w", err)
	}
	return nil
}

func (d *EmtctDb) ListLabTests(patientId int, lmp time.Time) ([]models.LabTest, error) {
	stmt := `
SELECT 
	id, patient_id, test_result, test_name, test_request_id, test_request_item_id,
	date_sample_taken, result_date, released_time, date_order_received_by_lab
FROM lab_tests
WHERE patient_id=$1 AND date_order_received_by_lab BETWEEN $2 AND $3
`
	endDate := lmp.Add(time.Hour * 24 * 7 * 52)
	rows, err := d.Query(stmt, patientId, lmp.Format(layoutISO), endDate.Format(layoutISO))
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("error retrieving lab tests from etmtct: %w", err)
	}
	var tests []models.LabTest
	for rows.Next() {
		var t models.LabTest
		err := rows.Scan(
			&t.Id,
			&t.PatientId,
			&t.TestResult,
			&t.TestName,
			&t.TestRequestId,
			&t.TestRequestItemId,
			&t.DateSampleTaken,
			&t.ResultDate,
			&t.ReleasedTime,
			&t.DateOrderReceivedByLab)
		if err != nil {
			return nil, fmt.Errorf("error scanning test result from etmct: %w", err)
		}
		tests = append(tests, t)
	}
	return tests, nil
}
