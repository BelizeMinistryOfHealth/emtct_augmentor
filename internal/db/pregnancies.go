package db

import (
	"context"
	"fmt"

	"moh.gov.bz/mch/emtct/internal/models"
)

func (d *AcsisDb) FindLatestPregnancy(patientId int) (*models.Pregnancy, error) {
	stmt := `
	SELECT 
		pregnancy_id, last_menstrual_period_date, estimated_delivery_date, end_time
	FROM acsis_hc_pregnancies 
	WHERE patient_id = $1
	ORDER BY last_menstrual_period_date DESC LIMIT 1;
`
	row := d.QueryRow(stmt, patientId)
	var pregnancy models.Pregnancy
	err := row.Scan(
		&pregnancy.PregnancyId,
		&pregnancy.Lmp,
		&pregnancy.Edd,
		&pregnancy.EndTime)
	if err != nil {
		return nil, fmt.Errorf("error retrieving pregnancy: %w", err)
	}
	pregnancy.PatientId = patientId
	return &pregnancy, nil
}

func (d *AcsisDb) FindPregnanciesByYear(year int) ([]models.Pregnancy, error) {
	stmt := `
	SELECT patient_id, pregnancy_id, last_menstrual_period_date, estimated_delivery_date, end_time
	FROM acsis_hc_pregnancies
	WHERE last_menstrual_period_date BETWEEN $1 AND $2;
`
	leftYear := fmt.Sprintf("%d-01-01", year)
	rightYear := fmt.Sprintf("%d-01-01", year+1)
	rows, err := d.Query(stmt, leftYear, rightYear)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("error querying for pregnancies by year: %w", err)
	}
	var ps []models.Pregnancy
	for rows.Next() {
		var p models.Pregnancy
		err := rows.Scan(
			&p.PatientId,
			&p.PregnancyId,
			&p.Lmp,
			&p.Edd,
			&p.EndTime)
		if err != nil {
			return nil, fmt.Errorf("error scanning pregnancy from acsis: %w", err)
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func (e *EmtctDb) FindExistingPregnanciesByYear(year int) ([]models.Pregnancy, error) {
	stmt := `
	SELECT pregnancy_id, patient_id, lmp, edd, end_time
	FROM pregnancies
	WHERE lmp BETWEEN $1 AND $2
`
	rows, err := e.Query(stmt, fmt.Sprintf("%d-01-01", year), fmt.Sprintf("%d-01-01", year+1))
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("error retrieving pregnancies from emtct db: %w", err)
	}
	var ps []models.Pregnancy
	for rows.Next() {
		var p models.Pregnancy
		err := rows.Scan(
			&p.PregnancyId,
			&p.PatientId,
			&p.Lmp,
			&p.Edd,
			&p.EndTime)
		if err != nil {
			return nil, fmt.Errorf("error scanning pregnancy from emtct db: %w", err)
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func (d *EmtctDb) InsertPregnancies(ctx context.Context, ps []models.Pregnancy) error {

	// Being transaction
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction for inserting pregnancies: %w", err)
	}

	stmt := `INSERT INTO pregnancies (pregnancy_id, patient_id, lmp, edd, end_time) VALUES($1, $2, $3, $4, $5)`
	for _, p := range ps {
		_, err := tx.ExecContext(ctx, stmt, p.PregnancyId, p.PatientId, p.Lmp, p.Edd, p.EndTime)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("")
		}
	}

	// Commit the transactions if there are no errors
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit the transaction for inserting pregnancies: %w", err)
	}

	return nil
}
