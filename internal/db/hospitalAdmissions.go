package db

import (
	"database/sql"
	"fmt"

	"moh.gov.bz/mch/emtct/internal/models"
)

func (d *EmtctDb) HospitalAdmissionsByPatientId(patientId int) ([]models.HospitalAdmission, error) {
	stmt := `SELECT id, patient_id, date_admitted, facility, created_at, created_by, updated_at, updated_by 
FROM hospital_admission WHERE patient_id=$1`
	var admissions []models.HospitalAdmission
	rows, err := d.DB.Query(stmt, patientId)
	if err != nil {
		return admissions, fmt.Errorf("error when executing query to retrieve hospital admissions for a patient: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var h models.HospitalAdmission
		err := rows.Scan(
			&h.Id,
			&h.PatientId,
			&h.DateAdmitted,
			&h.Facility,
			&h.CreatedAt,
			&h.CreatedBy,
			&h.UpdatedAt,
			&h.UpdatedBy)
		if err != nil {
			return admissions, fmt.Errorf("error scanning hotpsital admissions result from the database: %+v", err)
		}
		admissions = append(admissions, h)
	}
	return admissions, nil
}

// FindHospitalAdmissionById returns a hospital admission that was created by the EMTCT tool. It does not
// search for hospital admissions in the BHIS.
func (d *EmtctDb) FindHospitalAdmissionById(id string) (*models.HospitalAdmission, error) {
	stmt := `SELECT id, patient_id, date_admitted, facility, created_at, created_by, updated_at, updated_by
FROM hospital_admission WHERE id=$1`
	var admission models.HospitalAdmission
	row := d.DB.QueryRow(stmt, id)
	err := row.Scan(
		&admission.Id,
		&admission.PatientId,
		&admission.DateAdmitted,
		&admission.Facility,
		&admission.CreatedAt,
		&admission.CreatedBy,
		&admission.UpdatedAt,
		&admission.UpdatedBy)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &admission, nil
	default:
		return nil, fmt.Errorf("error retrieving hospital admission from database: %+v", err)
	}
}

func (d *EmtctDb) CreateHospitalAdmission(h models.HospitalAdmission) error {
	stmt := `INSERT INTO hospital_admission (id, patient_id, date_admitted, facility, created_at, created_by) 
Values($1, $2, $3, $4, $5, $6)`
	_, err := d.DB.Exec(stmt, h.Id, h.PatientId, h.DateAdmitted, h.Facility, h.CreatedAt, h.CreatedBy)
	if err != nil {
		return fmt.Errorf("error inserting a new hospital admission into the database: %+v", err)
	}
	return nil
}

func (d *EmtctDb) EditHospitalAdmission(h models.HospitalAdmission) error {
	stmt := `UPDATE hospital_admission SET date_admitted=$1, facility=$2, updated_at=$3, updated_by=$4 
WHERE id=$5`
	_, err := d.DB.Exec(stmt, h.DateAdmitted, h.Facility, h.UpdatedAt, h.UpdatedBy, h.Id)
	if err != nil {
		return fmt.Errorf("error updating a hospital admission in the database: %+v", err)
	}
	return nil
}
