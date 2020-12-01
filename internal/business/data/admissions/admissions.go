package admissions

import (
	"database/sql"
	"fmt"
)

func (a *Admissions) FindByPatientId(patientId int) ([]HospitalAdmission, error) {
	stmt := `
	SELECT 
	       id, patient_id, date_admitted, facility, reason, created_at, created_by, updated_at, updated_by, mch_encounter_id 
	FROM 
	     hospital_admission 
	WHERE 
	      patient_id=$1`
	var admissions []HospitalAdmission
	rows, err := a.Query(stmt, patientId)
	defer rows.Close()
	if err != nil {
		return admissions, fmt.Errorf("error when executing query to retrieve hospital admissions for a patient: %+v", err)
	}

	for rows.Next() {
		var h HospitalAdmission
		err := rows.Scan(
			&h.Id,
			&h.PatientId,
			&h.DateAdmitted,
			&h.Facility,
			&h.Reason,
			&h.CreatedAt,
			&h.CreatedBy,
			&h.UpdatedAt,
			&h.UpdatedBy,
			&h.MchEncounterId)
		if err != nil {
			return admissions, fmt.Errorf("error scanning hotpsital admissions result from the database: %+v", err)
		}
		admissions = append(admissions, h)
	}
	return admissions, nil
}

func (a *Admissions) FindById(id string) (*HospitalAdmission, error) {
	stmt := `
	SELECT id, patient_id, date_admitted, facility, reason, created_at, created_by, updated_at, updated_by, mch_encounter_id
	FROM hospital_admission 
	WHERE id=$1`
	var admission HospitalAdmission
	row := a.QueryRow(stmt, id)
	err := row.Scan(
		&admission.Id,
		&admission.PatientId,
		&admission.DateAdmitted,
		&admission.Facility,
		&admission.Reason,
		&admission.CreatedAt,
		&admission.CreatedBy,
		&admission.UpdatedAt,
		&admission.UpdatedBy,
		&admission.MchEncounterId)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &admission, nil
	default:
		return nil, fmt.Errorf("error retrieving hospital admission from database: %+v", err)
	}
}

func (a *Admissions) Create(h HospitalAdmission) error {
	stmt := `
	INSERT INTO hospital_admission 
	    (id, patient_id, date_admitted, facility, reason, created_at, created_by, mch_encounter_id) 
	Values
	       ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := a.Exec(stmt, h.Id, h.PatientId, h.DateAdmitted, h.Facility, h.Reason, h.CreatedAt, h.CreatedBy, h.MchEncounterId)
	if err != nil {
		return fmt.Errorf("error inserting a new hospital admission into the database: %+v", err)
	}
	return nil
}

func (a *Admissions) Edit(h HospitalAdmission) error {
	stmt := `
	UPDATE hospital_admission 
	SET date_admitted=$1, facility=$2, reason=$3, updated_at=$4, updated_by=$5 
	WHERE id=$6;
`
	_, err := a.Exec(stmt, h.DateAdmitted, h.Facility, h.Reason, h.UpdatedAt, h.UpdatedBy, h.Id)
	if err != nil {
		return fmt.Errorf("error updating a hospital admission in the database: %w", err)
	}
	return nil
}
