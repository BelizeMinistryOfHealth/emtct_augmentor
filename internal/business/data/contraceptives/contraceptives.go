package contraceptives

import (
	"database/sql"
	"fmt"
)

func (d *Contraceptives) Create(c ContraceptiveUsed) error {
	stmt := `
		INSERT INTO contraceptive_used 
    		(id, patient_id, contraceptive, comments, date_used, created_by, created_at, mch_encounter_id) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := d.Exec(stmt, c.Id, c.PatientId, c.Contraceptive, c.Comments, c.DateUsed, c.CreatedBy, c.CreatedAt, c.MchEncounterId)
	if err != nil {
		return fmt.Errorf("error inserting a new contraceptive into the database: %+v", err)
	}
	return nil
}

func (d *Contraceptives) Edit(c ContraceptiveUsed) error {
	stmt := `
	UPDATE contraceptive_used 
	SET contraceptive=$1, comments=$2, date_used=$3, updated_by=$4, updated_at=$5 
	WHERE id=$6;
`
	_, err := d.Exec(stmt, c.Contraceptive, c.Comments, c.DateUsed, c.UpdatedBy, c.UpdatedAt, c.Id)
	if err != nil {
		return fmt.Errorf("error updating contraceptive in the database: %+v", err)
	}
	return nil
}

func (d *Contraceptives) FindById(id string) (*ContraceptiveUsed, error) {
	stmt := `
		SELECT 
		       id, patient_id, contraceptive, comments, created_at, created_by, updated_at, 
		       updated_by, mch_encounter_id 
		FROM contraceptive_used 
		WHERE 
		      id=$1;
`
	var contraceptive ContraceptiveUsed
	row := d.DB.QueryRow(stmt, id)
	err := row.Scan(
		&contraceptive.Id,
		&contraceptive.PatientId,
		&contraceptive.Contraceptive,
		&contraceptive.Comments,
		&contraceptive.CreatedAt,
		&contraceptive.CreatedBy,
		&contraceptive.UpdatedAt,
		&contraceptive.UpdatedBy,
		&contraceptive.MchEncounterId)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &contraceptive, nil
	default:
		return nil, fmt.Errorf("error retrieving contraceptive from database: %+v", err)
	}
}

func (d *Contraceptives) FindByPatientId(patientId int) ([]ContraceptiveUsed, error) {
	stmt := `
		SELECT 
		       id, patient_id, contraceptive, comments, date_used, created_at, created_by, 
		       updated_at, updated_by, mch_encounter_id 
		FROM 
		     contraceptive_used 
		WHERE patient_id=$1`
	var contraceptives []ContraceptiveUsed

	rows, err := d.Query(stmt, patientId)
	if err != nil {
		return contraceptives, fmt.Errorf("error when executing query to retrieve contraceptives by patient id: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c ContraceptiveUsed
		err := rows.Scan(
			&c.Id,
			&c.PatientId,
			&c.Contraceptive,
			&c.Comments,
			&c.DateUsed,
			&c.CreatedAt,
			&c.CreatedBy,
			&c.UpdatedAt,
			&c.UpdatedBy,
			&c.MchEncounterId)
		if err != nil {
			return contraceptives, fmt.Errorf("error scanning contraceptive result from the database: %+v", err)
		}
		contraceptives = append(contraceptives, c)
	}

	return contraceptives, nil
}
