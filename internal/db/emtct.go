package db

import (
	"database/sql"
	"fmt"

	"moh.gov.bz/mch/emtct/internal/config"
	"moh.gov.bz/mch/emtct/internal/models"
)

type EmtctDb struct {
	*sql.DB
}

// NewConnection creates a new database connection
func NewConnection(cnf *config.DbConf) (*EmtctDb, error) {
	connstr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", cnf.DbUsername, cnf.DbDatabase, cnf.DbPassword, cnf.DbHost)
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

func (d *EmtctDb) FindDiagnoses(patientId string) ([]models.Diagnosis, error) {
	stmt := `SELECT id, patient_id, diagnosis_date, diagnosis_name FROM diagnoses WHERE patient_id=$1`
	var diagnoses []models.Diagnosis
	rows, err := d.DB.Query(stmt, patientId)
	if err != nil {
		return nil, fmt.Errorf("error execcuting query for extracting diagnoses: %+v", err)
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
