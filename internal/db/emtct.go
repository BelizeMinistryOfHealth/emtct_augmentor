package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

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

// FindCurrentPregnancy returns the current pregnancy for the specified patient.
// The pregnancy is deemed current if the EDD is in the future.
func (d *EmtctDb) FindCurrentPregnancy(patientId string) (*models.PregnancyVitals, error) {
	stmt := `SELECT id, patient_id, gestational_age, para, cs, abortive_outcome, diagnosis_date, planned,
age_at_lmp, lmp, edd, date_of_booking, prenatal_care_provider, total_checks FROM pregnancies WHERE patient_id=$1`
	var pregnancies []models.PregnancyVitals
	id, err := strconv.Atoi(patientId)
	if err != nil {
		return nil, fmt.Errorf("error: patient id is not a number: %+v", err)
	}
	rows, err := d.DB.Query(stmt, id)
	if err != nil {
		return nil, fmt.Errorf("error executing query for extracting the current pregnancy: %+v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var p models.PregnancyVitals
		err := rows.Scan(
			&p.Id,
			&p.PatientId,
			&p.GestationalAge,
			&p.Para,
			&p.Cs,
			&p.AbortiveOutcome,
			&p.DiagnosisDate,
			&p.Planned,
			&p.AgeAtLmp,
			&p.Lmp,
			&p.Edd,
			&p.DateOfBooking,
			&p.PrenatalCareProvider,
			&p.TotalChecks)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pregnancy for patient(%s): %+v", patientId, err)
		}
		pregnancies = append(pregnancies, p)
	}

	p := models.FindCurrentPregnancy(pregnancies)
	return p, nil
}

func (d *EmtctDb) FindPregnancyDiagnoses(patientId string) ([]models.Diagnosis, error) {
	pregnancy, err := d.FindCurrentPregnancy(patientId)
	if err != nil {
		return nil, fmt.Errorf("error fetching the current pregnancy when retrieving the diagnoses: %+v", err)
	}
	if pregnancy == nil {
		return []models.Diagnosis{}, nil
	}
	edd := pregnancy.Edd

	diagnoses, err := d.FindDiagnoses(patientId)
	if err != nil {
		return nil, fmt.Errorf("error fetching diagnoses for the current pregnancy: %+v", err)
	}

	var pregnancyDiagnoses []models.Diagnosis
	for _, v := range diagnoses {
		if v.Date.Before(edd) && v.Date.After(edd.Add(-time.Hour*24*30*9)) {
			pregnancyDiagnoses = append(pregnancyDiagnoses, v)
		}
	}

	return pregnancyDiagnoses, nil
}
