package partners

import (
	"database/sql"
	"fmt"

	"moh.gov.bz/mch/emtct/internal/business/data/prescription"
	"moh.gov.bz/mch/emtct/internal/db"
)

type Partners struct {
	emtctdb *db.EmtctDb
}

func New(db *db.EmtctDb) Partners {
	return Partners{emtctdb: db}
}

func (p *Partners) AddPartnerSyphilisTreatment(treatment prescription.SyphilisTreatment) error {
	stmt := `
	INSERT INTO syphilis_treatment_partner 
    	(id, patient_id, medication_name, dosage, comments, date, created_by, created_at)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8);
`
	_, err := p.emtctdb.Exec(stmt,
		treatment.Id,
		treatment.PatientId,
		treatment.Medication,
		treatment.Dosage,
		treatment.Comments,
		treatment.Date,
		treatment.CreatedBy,
		treatment.CreatedAt)
	if err != nil {
		return fmt.Errorf("error inserting syphilis treatment for partner into database: %w", err)
	}
	return nil
}

func (p *Partners) FindPartnerSyphilisTreatments(patientId int) ([]prescription.SyphilisTreatment, error) {
	stmt := `
	SELECT 
	       id, patient_id, medication_name, dosage, comments, date, created_by, created_at, updated_by, updated_at
	FROM syphilis_treatment_partner
	WHERE patient_id=$1
	ORDER BY date DESC;
`
	rows, err := p.emtctdb.Query(stmt, patientId)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("error querying partner syphilis treatment from database: %w", err)
	}
	var treatments []prescription.SyphilisTreatment
	for rows.Next() {
		var t prescription.SyphilisTreatment
		var updatedBy sql.NullString
		err := rows.Scan(
			&t.Id,
			&t.PatientId,
			&t.Medication,
			&t.Dosage,
			&t.Comments,
			&t.Date,
			&t.CreatedBy,
			&t.CreatedAt,
			&updatedBy,
			&t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning syphilis treatment for partner query results: %+v", err)
		}
		if updatedBy.Valid {
			t.UpdatedBy = updatedBy.String
		}
		treatments = append(treatments, t)
	}
	return treatments, nil
}

func (p *Partners) UpdatePartnerSyphilisTreatment(treatment prescription.SyphilisTreatment) error {
	stmt := `
	UPDATE syphilis_treatment_partner 
	SET medication_name=$1, dosage=$2, comments=$3, updated_by=$4, updated_at=$5, date=$6
	WHERE id=$7;
`
	_, err := p.emtctdb.Exec(stmt,
		treatment.Medication,
		treatment.Dosage,
		treatment.Comments,
		treatment.UpdatedBy,
		treatment.UpdatedAt,
		treatment.Date,
		treatment.Id)
	if err != nil {
		return fmt.Errorf("error updating a partner's syphilis treatment in the database: %w", err)
	}
	return nil
}
