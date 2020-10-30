package db

import (
	"database/sql"
	"fmt"

	"moh.gov.bz/mch/emtct/internal/models"
)

func (d *EmtctDb) ContraceptivesUsedByPatientId(patientId int) ([]models.ContraceptiveUsed, error) {
	stmt := `SELECT id, patient_id, contraceptive, comments, date_used, created_at, created_by, updated_at, updated_by FROM contraceptive_used WHERE patient_id=$1`
	var contraceptives []models.ContraceptiveUsed

	rows, err := d.DB.Query(stmt, patientId)
	if err != nil {
		return contraceptives, fmt.Errorf("error when executing query to retrieve contraceptives by patient id: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.ContraceptiveUsed
		err := rows.Scan(
			&c.Id,
			&c.PatientId,
			&c.Contraceptive,
			&c.Comments,
			&c.DateUsed,
			&c.CreatedAt,
			&c.CreatedBy,
			&c.UpdatedAt,
			&c.UpdatedBy)
		if err != nil {
			return contraceptives, fmt.Errorf("error scanning contraceptive result from the database: %+v", err)
		}
		contraceptives = append(contraceptives, c)
	}

	return contraceptives, nil
}

func (d *EmtctDb) CreateContraceptiveUsed(c models.ContraceptiveUsed) error {
	stmt := `INSERT INTO contraceptive_used (id, patient_id, contraceptive, comments, date_used, created_by, created_at) 
VALUES($1, $2, $3, $4, $5, $6, $7)`
	_, err := d.DB.Exec(stmt, c.Id, c.PatientId, c.Contraceptive, c.Comments, c.DateUsed, c.CreatedBy, c.CreatedAt)
	if err != nil {
		return fmt.Errorf("error inserting a new contraceptive into the database: %+v", err)
	}
	return nil
}

func (d *EmtctDb) EditContraceptiveUsed(c models.ContraceptiveUsed) error {
	stmt := `UPDATE contraceptive_used SET contraceptive=$1, comments=$2, date_used=$3, updated_by=$4, updated_at=$5 WHERE id=$6`
	_, err := d.DB.Exec(stmt, c.Contraceptive, c.Comments, c.DateUsed, c.UpdatedBy, c.UpdatedAt, c.Id)
	if err != nil {
		return fmt.Errorf("error updating contraceptive in the database: %+v", err)
	}
	return nil
}

func (d *EmtctDb) FindContraceptiveById(id string) (*models.ContraceptiveUsed, error) {
	stmt := `SELECT id, patient_id, contraceptive, comments, created_at, created_by, updated_at, updated_by FROM 
contraceptive_used WHERE id=$1`
	var contraceptive models.ContraceptiveUsed
	row := d.DB.QueryRow(stmt, id)
	err := row.Scan(
		&contraceptive.Id,
		&contraceptive.PatientId,
		&contraceptive.Contraceptive,
		&contraceptive.Comments,
		&contraceptive.CreatedAt,
		&contraceptive.CreatedBy,
		&contraceptive.UpdatedAt,
		&contraceptive.UpdatedBy)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &contraceptive, nil
	default:
		return nil, fmt.Errorf("error retrieving contraceptive from database: %+v", err)
	}
}