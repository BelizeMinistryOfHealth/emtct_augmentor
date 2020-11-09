package db

import (
	"database/sql"
	"fmt"
	"time"

	"moh.gov.bz/mch/emtct/internal/models"
)

func (d *EmtctDb) CreateHivScreening(v models.HivScreening) error {
	stmt := `
	INSERT INTO hiv_Screening 
    	(id, patient_id, test_name, screening_date, date_sample_received_at_hq, sample_code,
		date_sample_shipped, destination, date_result_received, result, date_result_shared, created_at, created_by, 
		date_sample_taken, mch_encounter_id, timely, due_date)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`
	_, err := d.DB.Exec(stmt,
		v.Id,
		v.PatientId,
		v.TestName,
		v.ScreeningDate,
		v.DateSampleReceivedAtHq,
		v.SampleCode,
		v.DateSampleShipped,
		v.Destination,
		v.DateResultReceived,
		v.Result,
		v.DateResultShared,
		v.CreatedAt,
		v.CreatedBy,
		v.DateSampleTaken,
		v.MchEncounterId,
		v.Timely,
		v.DueDate)
	if err != nil {
		return fmt.Errorf("error inserting new hiv screening into database: %+v", err)
	}

	return nil
}

func (d *EmtctDb) FindHivScreeningsByPatient(patientId int) ([]models.HivScreening, error) {
	stmt := `
	SELECT 
		id, patient_id, mch_encounter_id, test_name, screening_date, date_sample_received_at_hq, sample_code,
		date_sample_shipped, date_sample_taken, destination, date_result_received, result, date_result_shared, 
		created_at, created_by, updated_at, updated_by, timely, due_date 
	FROM hiv_screening 
	WHERE patient_id=$1
`

	var screenings []models.HivScreening

	rows, err := d.DB.Query(stmt, patientId)
	if err != nil {
		return screenings, fmt.Errorf("error querying hiv screenings: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var s models.HivScreening
		err := rows.Scan(
			&s.Id,
			&s.PatientId,
			&s.MchEncounterId,
			&s.TestName,
			&s.ScreeningDate,
			&s.DateSampleReceivedAtHq,
			&s.SampleCode,
			&s.DateSampleShipped,
			&s.DateSampleTaken,
			&s.Destination,
			&s.DateResultReceived,
			&s.Result,
			&s.DateResultShared,
			&s.CreatedAt,
			&s.CreatedBy,
			&s.UpdatedAt,
			&s.UpdatedBy,
			&s.Timely,
			&s.DueDate)
		if err != nil {
			return screenings, fmt.Errorf("error scanning hiv screening row: %+v", err)
		}
		screenings = append(screenings, s)
	}

	return screenings, nil
}

func (d *EmtctDb) EditHivScreening(v models.HivScreening) (*models.HivScreening, error) {
	stmt := `
	UPDATE hiv_screening 
	SET test_name=$1, result=$2, sample_code=$3, destination=$4, screening_date=$5, date_sample_received_at_hq=$6, 
	    date_sample_shipped=$7, date_result_received=$8, date_result_shared=$9, updated_at=$10, updated_by=$11, 
	    date_sample_taken=$12, timely=$13
	WHERE id=$14
`
	updatedAt := time.Now()
	_, err := d.DB.Exec(stmt,
		v.TestName,
		v.Result,
		v.SampleCode,
		v.Destination,
		v.ScreeningDate,
		v.DateSampleReceivedAtHq,
		v.DateSampleShipped,
		v.DateResultReceived,
		v.DateResultShared,
		updatedAt,
		v.UpdatedBy,
		v.DateSampleTaken,
		v.Timely,
		v.Id)
	if err != nil {
		return nil, fmt.Errorf("error updating hiv screening in database: %+v", err)
	}
	v.UpdatedAt = &updatedAt
	return &v, nil
}

func (d *EmtctDb) FindHivScreeningById(id string) (*models.HivScreening, error) {
	stmt := `
	SELECT 
		id, patient_id, mch_encounter_id, test_name, result, sample_code, destination, screening_date,
		date_sample_received_at_hq, date_sample_shipped, date_sample_taken, date_result_received, date_result_shared, 
		updated_at, updated_by, timely, due_date
	FROM hiv_screening 
	WHERE id=$1`
	var screening models.HivScreening
	row := d.DB.QueryRow(stmt, id)
	err := row.Scan(
		&screening.Id,
		&screening.PatientId,
		&screening.MchEncounterId,
		&screening.TestName,
		&screening.Result,
		&screening.SampleCode,
		&screening.Destination,
		&screening.ScreeningDate,
		&screening.DateSampleReceivedAtHq,
		&screening.DateSampleShipped,
		&screening.DateSampleTaken,
		&screening.DateResultReceived,
		&screening.DateResultShared,
		&screening.UpdatedAt,
		&screening.UpdatedBy,
		&screening.Timely,
		&screening.DueDate)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &screening, nil
	default:
		return nil, fmt.Errorf("error retrieving hiv screening from database: %+v", err)
	}
}
