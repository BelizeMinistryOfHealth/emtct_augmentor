package contactTracing

import (
	"database/sql"
	"fmt"
)

func (d *ContactTracings) Create(c ContactTracing) error {
	stmt := `
	INSERT INTO 
	    contact_tracing (id, patient_id, test, test_result, comments, date, created_by, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
`
	_, err := d.Exec(stmt,
		c.Id,
		c.PatientId,
		c.Test,
		c.TestResult,
		c.Comments,
		c.Date,
		c.CreatedBy,
		c.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("error inserting contact tracing to the database: %+v", err)
	}
	return nil
}

func (d *ContactTracings) FindByPatientId(patientId int) ([]ContactTracing, error) {
	stmt := `
	SELECT 
	       id, patient_id, test, test_result, comments, date, created_by, created_at, updated_by, updated_at
	FROM contact_tracing
	WHERE patient_id=$1;
`
	rows, err := d.Query(stmt, patientId)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("error retrieving contact tracing from database: %+v", err)
	}
	var contacts []ContactTracing
	for rows.Next() {
		var c ContactTracing
		var updatedBy sql.NullString
		err := rows.Scan(
			&c.Id,
			&c.PatientId,
			&c.Test,
			&c.TestResult,
			&c.Comments,
			&c.Date,
			&c.CreatedBy,
			&c.CreatedAt,
			&updatedBy,
			&c.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning contact tracing record: %+v", err)
		}
		if updatedBy.Valid {
			c.UpdatedBy = updatedBy.String
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (d *ContactTracings) Edit(c ContactTracing) error {
	stmt := `
	UPDATE 
	    contact_tracing SET test=$1, test_result=$2, comments=$3, date=$4, updated_by=$5, updated_at=$6
	WHERE id = $7
`

	_, err := d.Exec(stmt,
		c.Test,
		c.TestResult,
		c.Comments,
		c.Date,
		c.UpdatedBy,
		c.UpdatedAt,
		c.Id,
	)
	if err != nil {
		return fmt.Errorf("error updating contact tracing in database: %+v", err)
	}
	return nil
}
