package infant

import (
	"database/sql"
	"fmt"
	"time"
)

func (d *Infants) CreateHivScreening(v HivScreening) error {
	stmt := `
	INSERT INTO hiv_Screening 
    	(id, patient_id, test_name, screening_date, date_sample_received_at_hq, sample_code,
		date_sample_shipped, destination, date_result_received, result, date_result_shared, created_at, created_by, 
		date_sample_taken, mother_id, timely, due_date)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`
	_, err := d.Acsis.Exec(stmt,
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
		v.MotherId,
		v.Timely,
		v.DueDate)
	if err != nil {
		return fmt.Errorf("error inserting new hiv screening into database: %+v", err)
	}

	return nil
}

func (d *Infants) FindHivScreeningsByPatient(patientId int) ([]HivScreening, error) {
	stmt := `
	SELECT 
		id, patient_id, mother_id, test_name, screening_date, date_sample_received_at_hq, sample_code,
		date_sample_shipped, date_sample_taken, destination, date_result_received, result, date_result_shared, 
		created_at, created_by, updated_at, updated_by, timely, due_date 
	FROM hiv_screening 
	WHERE patient_id=$1
`

	var screenings []HivScreening

	rows, err := d.Acsis.Query(stmt, patientId)
	if err != nil {
		return screenings, fmt.Errorf("error querying hiv screenings: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var s HivScreening
		err := rows.Scan(
			&s.Id,
			&s.PatientId,
			&s.MotherId,
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

func (d *Infants) EditHivScreening(v HivScreening) (*HivScreening, error) {
	stmt := `
	UPDATE hiv_screening 
	SET test_name=$1, result=$2, sample_code=$3, destination=$4, screening_date=$5, date_sample_received_at_hq=$6, 
	    date_sample_shipped=$7, date_result_received=$8, date_result_shared=$9, updated_at=$10, updated_by=$11, 
	    date_sample_taken=$12, timely=$13
	WHERE id=$14
`
	updatedAt := time.Now()
	_, err := d.Acsis.Exec(stmt,
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

func (d *Infants) FindHivScreeningById(id string) (*HivScreening, error) {
	stmt := `
	SELECT 
		id, patient_id, mother_id, test_name, result, sample_code, destination, screening_date,
		date_sample_received_at_hq, date_sample_shipped, date_sample_taken, date_result_received, date_result_shared, 
		updated_at, updated_by, timely, due_date
	FROM hiv_screening 
	WHERE id=$1`
	var screening HivScreening
	row := d.Acsis.QueryRow(stmt, id)
	err := row.Scan(
		&screening.Id,
		&screening.PatientId,
		&screening.MotherId,
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

// IsHivScreeningTimely indicates if an hiv screening was done in a timely manner.
// The timeliness depends on the type of test and when the sample was taken:
// PCR 1: sample must be taken 3 days or less after birth.
// PCR 2: sample must be taken no later than 6 weeks after birth
// PCR 3: sample must be taken no later than 90 days after birth
// ELISA: sample must be taken no longer than 18 months after birth
func (d *Infants) IsHivScreeningTimely(birthDate time.Time, testName string, dateSampleTaken time.Time) bool {
	diff := dateSampleTaken.Sub(birthDate).Hours() / 24
	switch testName {
	case "PCR 1":
		return diff < 4
	case "PCR 2":
		return diff < (6 * 7)
	case "PCR 3":
		return diff < 91
	case "ELISA":
		return diff <= (18 * 7 * 4)
	default:
		return false
	}
}

// HivScreeningDueDate calculates the due date for taking a sample for an HIV screening.
// PCR 1: sample must be taken 3 days or less after birth.
// PCR 2: sample must be taken no later than 6 weeks after birth
// PCR 3: sample must be taken no later than 90 days after birth
// ELISA: sample must be taken no longer than 18 months after birth
func (d *Infants) HivScreeningDueDate(testName string, birthDate time.Time) time.Time {
	switch testName {
	case "PCR 1":
		return birthDate.AddDate(0, 0, 3)
	case "PCR 2":
		return birthDate.AddDate(0, 0, 42)
	case "PCR 3":
		return birthDate.AddDate(0, 0, 90)
	default:
		return birthDate.AddDate(0, 18, 0)
	}
}
