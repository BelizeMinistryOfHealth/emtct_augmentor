package db

import (
	"database/sql"
	"fmt"
	"time"

	"moh.gov.bz/mch/emtct/internal/models"
)

func (d *EmtctDb) FindHomeVisitsByPatientId(patientId int) ([]models.HomeVisit, error) {

	stmt := `
	SELECT id, patient_id, reason, comments, date_of_visit, created_at, updated_at, created_by, updated_by, mch_encounter_id 
	FROM home_visit WHERE patient_id=$1`
	rows, err := d.DB.Query(stmt, patientId)
	if err != nil {
		return nil, fmt.Errorf("error executing query for retrieving home visits: %+v", err)
	}
	defer rows.Close()

	var homeVisits []models.HomeVisit
	for rows.Next() {
		var homeVisit models.HomeVisit
		err := rows.Scan(
			&homeVisit.Id,
			&homeVisit.PatientId,
			&homeVisit.Reason,
			&homeVisit.Comments,
			&homeVisit.DateOfVisit,
			&homeVisit.CreatedAt,
			&homeVisit.UpdatedAt,
			&homeVisit.CreatedBy,
			&homeVisit.UpdatedBy,
			&homeVisit.MchEncounterId,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning a home visit row: %+v", err)
		}

		homeVisits = append(homeVisits, homeVisit)
	}
	return homeVisits, nil
}

func (d *EmtctDb) FindHomeVisitById(id string) (*models.HomeVisit, error) {
	stmt := `
	SELECT id, patient_id, reason, comments, date_of_visit, created_at, updated_at, created_by, updated_by, mch_encounter_id 
	FROM home_visit WHERE id=$1`
	var homeVisit models.HomeVisit
	row := d.DB.QueryRow(stmt, id)
	err := row.Scan(
		&homeVisit.Id,
		&homeVisit.PatientId,
		&homeVisit.Reason,
		&homeVisit.Comments,
		&homeVisit.DateOfVisit,
		&homeVisit.CreatedAt,
		&homeVisit.UpdatedBy,
		&homeVisit.CreatedBy,
		&homeVisit.UpdatedBy,
		&homeVisit.MchEncounterId,
	)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &homeVisit, nil
	default:
		return nil, fmt.Errorf("error scanning home visit row: %+v", err)
	}
}

func (d *EmtctDb) CreateHomeVisit(v models.HomeVisit) error {
	stmt := `
	INSERT INTO home_visit 
	    (id, patient_id, reason, comments, date_of_visit, created_at, created_by, mch_encounter_id) 
	    VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := d.DB.Exec(stmt, v.Id, v.PatientId, v.Reason, v.Comments, v.DateOfVisit, v.CreatedAt, v.CreatedBy, v.MchEncounterId)
	if err != nil {
		return fmt.Errorf("error creating a home visit: %+v", err)
	}
	return nil
}

func (d *EmtctDb) EditHomeVisit(v models.HomeVisit) (*models.HomeVisit, error) {
	stmt := `
	UPDATE home_visit 
	SET reason=$1, comments=$2, date_of_visit=$3, updated_by=$4, updated_at=$5
	WHERE id=$6`
	updateddAt := time.Now()
	_, err := d.DB.Exec(stmt,
		v.Reason,
		v.Comments,
		v.DateOfVisit,
		v.UpdatedBy,
		updateddAt,
		v.Id)
	if err != nil {
		return nil, fmt.Errorf("error updating homve visit in database: %+v", err)
	}
	v.UpdatedAt = &updateddAt
	return &v, nil
}
