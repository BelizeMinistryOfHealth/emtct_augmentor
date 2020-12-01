package pregnancy

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bearbin/go-age"
)

func (p Pregnancies) FindLatest(patientId int) (*Pregnancy, error) {
	stmt := `
	SELECT 
	       pregnancy_id, patient_id, lmp, edd, end_time
	FROM pregnancies
	WHERE 
	      patient_id = $1
	ORDER BY lmp DESC
	LIMIT 1;
`
	row := p.EmtctDb.QueryRow(stmt, patientId)
	var pregnancy Pregnancy
	err := row.Scan(
		&pregnancy.PregnancyId,
		&pregnancy.PatientId,
		&pregnancy.Lmp,
		&pregnancy.Edd,
		&pregnancy.EndTime)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &pregnancy, nil
	default:
		return nil, fmt.Errorf("error retrieving pregnancy from emtctdb: %w", err)
	}
}

func (p Pregnancies) FindPregnanciesInBhisByYear(year int) ([]Pregnancy, error) {
	stmt := `
	SELECT patient_id, pregnancy_id, last_menstrual_period_date, estimated_delivery_date, end_time
	FROM acsis_hc_pregnancies
	WHERE last_menstrual_period_date BETWEEN $1 AND $2;
`
	leftYear := fmt.Sprintf("%d-01-01", year)
	rightYear := fmt.Sprintf("%d-01-01", year+1)
	rows, err := p.AcsisDb.Query(stmt, leftYear, rightYear)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("error querying for pregnancies by year: %w", err)
	}
	var ps []Pregnancy
	for rows.Next() {
		var pr Pregnancy
		err := rows.Scan(
			&pr.PatientId,
			&pr.PregnancyId,
			&pr.Lmp,
			&pr.Edd,
			&pr.EndTime)
		if err != nil {
			return nil, fmt.Errorf("error scanning pregnancy from acsis: %w", err)
		}
		ps = append(ps, pr)
	}
	return ps, nil
}

func (p Pregnancies) FindExistingPregnanciesByYear(year int) ([]Pregnancy, error) {
	stmt := `
	SELECT pregnancy_id, patient_id, lmp, edd, end_time
	FROM pregnancies
	WHERE lmp BETWEEN $1 AND $2
`
	rows, err := p.EmtctDb.Query(stmt, fmt.Sprintf("%d-01-01", year), fmt.Sprintf("%d-01-01", year+1))
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("error retrieving pregnancies from emtct db: %w", err)
	}
	var ps []Pregnancy
	for rows.Next() {
		var pr Pregnancy
		err := rows.Scan(
			&pr.PregnancyId,
			&pr.PatientId,
			&pr.Lmp,
			&pr.Edd,
			&pr.EndTime)
		if err != nil {
			return nil, fmt.Errorf("error scanning pregnancy from emtct db: %w", err)
		}
		ps = append(ps, pr)
	}
	return ps, nil
}

func (p Pregnancies) Create(ctx context.Context, ps []Pregnancy) error {
	// Being transaction
	tx, err := p.EmtctDb.BeginTx(ctx, nil)
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

// FindCurrentPregnancy returns details about the patient's current pregnancy.
// Under the hood it is 4 separate queries:
// 1. Finds obstetric details. This is information in the acsis_hc_pregnancies table, and contains the LMP and EDD.
// 2. Finds the latest anc encounter.
// 3. Retrieves apgar information
// 4. Retrieves previous diagnoses
// These are all separate queries because the database is not designed in a way to make it possible to retrieve
// all this information using joins. This is partly due to there not being any link between the pregnancies table and
// the encounters table that allows us to retrieve the corresponding anc encounter for a pregnancy.
func (d *Pregnancies) FindCurrentPregnancy(patientId int) (*Vitals, error) {
	v, err := d.FindObstetricDetails(patientId)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving pregnancy details from acsis: %+v", err)
	}
	// Find the anc encounter. This is the most recent anc encounter in patient's docket.
	anc, err := d.FindLatestAntenatalEncounter(patientId, v.Lmp)
	if err != nil {
		return nil, fmt.Errorf("could not find current pregnancy details because no antenatal encounter was found: %+v", err)
	}

	stmt := `SELECT
       CASE
           WHEN ft.facility_type_id = 14 THEN 'Private'
           ELSE 'Public'
       END AS care_provider,
       f.name as facility,
       e.begin_time as date_of_booking,
       CASE
           WHEN bs.name IS NULL THEN ''
           ELSE bs.name
           END AS birth_status,
       ahipd.apgar_fifth_minute,
       ahipd.apgar_first_minute,
       p.birth_date
FROM acsis_hc_patients p
INNER JOIN acsis_adt_encounters e ON p.patient_id=e.patient_id AND e.encounter_id=$2
INNER JOIN acsis_hc_facilities f ON e.facility_id = f.facility_id
INNER JOIN acsis_hc_facility_types ft ON f.facility_type_id = ft.facility_type_id
LEFT JOIN acsis_adt_labour_encounter_details aaled ON e.encounter_id = aaled.mch_encounter_id
LEFT JOIN acsis_hc_birth_statuses bs ON aaled.birth_status_id=bs.birth_status_id
LEFT JOIN acsis_hc_infant_patient_details ahipd ON p.patient_id = ahipd.mothers_patient_id
WHERE p.patient_id=$1
ORDER BY e.begin_time DESC
LIMIT 1;
`
	var vitals Vitals
	var dob *time.Time
	var careProvider string
	var facility string
	row := d.AcsisDb.QueryRow(stmt, patientId, anc.Id)
	var apgarFirst sql.NullInt32
	var apgarFifth sql.NullInt32
	err = row.Scan(
		&careProvider,
		&facility,
		&vitals.DateOfBooking,
		&vitals.BirthStatus,
		&apgarFifth,
		&apgarFirst,
		&dob,
	)
	vitals.PrenatalCareProvider = fmt.Sprintf("%s - %s", careProvider, facility)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		ageAtLmp := 0
		if v.Lmp != nil {
			ageAtLmp = age.AgeAt(*dob, *v.Lmp)
		}
		vitals.AgeAtLmp = ageAtLmp
		vitals.Para = v.Para
		vitals.Id = v.Id
		vitals.Cs = v.Cs
		vitals.Planned = v.Planned
		vitals.Lmp = v.Lmp
		vitals.Edd = v.Edd
		vitals.PregnancyOutcome, err = d.abortiveOutcome(vitals)
		if err != nil {
			return nil, fmt.Errorf("error while calculating abortive outcome when retrieving pregnancy details from acsis: %+v", err)
		}
		p, err := d.findPreviousPregnancyDiagnosis(patientId)
		if err != nil {
			return nil, fmt.Errorf("error while retrieving current pregnancy info from acsis: %+v", err)
		}
		if p != nil {
			vitals.DiagnosisDate = &p.Date
		}
		if apgarFirst.Valid {
			vitals.ApgarFirstMinute = int(apgarFirst.Int32)
		}
		if apgarFifth.Valid {
			vitals.ApgarFifthMinute = int(apgarFifth.Int32)
		}
		vitals.GestationalAge = anc.GestationalAge

		return &vitals, nil
	default:
		return nil, fmt.Errorf("error while retrieving pregnancy details from acsis: %+v", err)
	}

}

// FindObstetricDetails retrieves the patient's LMP, EDD from the acsis_hc_pregnancies
// and the Para/Cs/Planned data as a second query from the acsis_hc_obstetric_patient_details
// table because there is no guarantee that the patient will have data in the latter table.
// When no data is available in that table, we get the wrong results.
func (d *Pregnancies) FindObstetricDetails(patientId int) (*Vitals, error) {
	stmt := `SELECT
       hp.pregnancy_id,
       hp.last_menstrual_period_date,
       hp.estimated_delivery_date
FROM acsis_hc_patients p
INNER JOIN acsis_hc_pregnancies hp ON p.patient_id = hp.patient_id
WHERE p.patient_id=$1
ORDER BY hp.last_modified_time DESC
LIMIT 1;`

	var vitals Vitals
	row := d.AcsisDb.QueryRow(stmt, patientId)
	err := row.Scan(&vitals.Id,
		&vitals.Lmp,
		&vitals.Edd)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		v, err := d.findObstetricPatientDetails(patientId)
		if err != nil {
			return nil, fmt.Errorf("error retrieving obstetric patient detials when consolidating the obstetric details: %+v", err)
		}
		if v != nil {
			vitals.Para = v.Para
			vitals.Cs = v.Cs
			vitals.Planned = v.Planned
		}

		return &vitals, nil
	default:
		return nil, fmt.Errorf("error querying obstetric details from acsis: %+v", err)
	}
}

// findObstetricPatientDetails retrieves the number of liveborn pregnancies for the patient,
// and previous C/S and planned pregnancies. Not all patients will have data in this table.
// If we do a join of this table when trying to retrieve other obstetric information from the
// pregnancies table, the results are skewed because we can not guarantee that the patient will
// exist in the acsis_hc_obstetric_patient_details_table
func (d *Pregnancies) findObstetricPatientDetails(patientId int) (*Vitals, error) {
	stmt := `
	SELECT
       ahopd.number_liveborn_pregnancies,
       ahopd.number_caesarean_sections,
       ahopd.previous_pregnancy_planned
	FROM acsis_hc_patients p
	INNER JOIN acsis_hc_obstetric_patient_details ahopd on p.obstetric_patient_details_id = ahopd.obstetric_patient_details_id
	WHERE p.patient_id=$1 
	LIMIT 1;
`
	var vitals Vitals
	row := d.AcsisDb.QueryRow(stmt, patientId)
	err := row.Scan(
		&vitals.Para,
		&vitals.Cs,
		&vitals.Planned)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &vitals, nil
	default:
		return nil, fmt.Errorf("error querying obstetric patient details from acsis: %+v", err)
	}

}

func (d *Pregnancies) FindLatestAntenatalEncounter(patientId int, lmp *time.Time) (*AntenatalEncounter, error) {
	stmt := `SELECT e.encounter_id,
           e.patient_id,
           amed.mch_encounter_details_id, 
           amed.estimated_delivery_date,
           e.begin_time,
           amed.number_of_antenatal_visits
        FROM acsis_hc_patients p
        INNER JOIN acsis_adt_encounters e ON p.patient_id=e.patient_id AND e.encounter_type='M'
        INNER JOIN acsis_adt_mch_encounter_details amed ON e.encounter_details_id=amed.mch_encounter_details_id
        WHERE p.patient_id=$1
        ORDER BY e.begin_time DESC
        LIMIT 1;`

	var anc AntenatalEncounter
	row := d.AcsisDb.QueryRow(stmt, patientId)
	err := row.Scan(&anc.Id,
		&anc.PatientId,
		&anc.MchEncounterDetailsId,
		&anc.EstimatedDeliveryDate,
		&anc.BeginDate,
		&anc.NumberAntenatalVisits,
	)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		//Gestational Age at booking is the difference b/w LMP and begin time
		if lmp == nil {
			return &anc, nil
		}
		anc.GestationalAge = int(anc.BeginDate.Sub(*lmp).Hours() / 24)
		return &anc, nil
	default:
		return nil, fmt.Errorf("error querying for mch details from acsis: %+v", err)
	}
}

func (d *Pregnancies) abortiveOutcome(v Vitals) (string, error) {
	if v.ApgarFifthMinute > 0 && v.ApgarFirstMinute > 0 {
		return "Live Birth", nil
	}

	if v.GestationalAge >= 22 && v.GestationalAge <= 27 {
		return "Still Birth 22", nil
	}

	if v.GestationalAge > 27 {
		return "Still Birth 28", nil
	}

	// Otherwise it is an abortion.. and we need to do a query for this:
	stmt := `
	SELECT aai10d.name
	FROM acsis_adt_encounters AS e
	INNER JOIN acsis_adt_encounter_diagnoses aaed on e.encounter_id = aaed.encounter_id
	INNER JOIN acsis_adt_icd10_diseases aai10d on aaed.disease_id = aai10d.disease_id
	INNER JOIN acsis_hc_pregnancies ahp on e.patient_id = ahp.patient_id
	WHERE e.patient_id = $1 
		AND aaed.diagnosis_time < ahp.last_menstrual_period_date
		AND (aai10d.code ILIKE 'O06%' OR aai10d.code ILIKE 'O03%' OR aai10d.code ILIKE 'O05%'
		OR aai10d.code ILIKE 'O04%')
	ORDER BY aaed.diagnosis_time DESC
	LIMIT 1;
`
	var diagnosis string
	row := d.AcsisDb.QueryRow(stmt, v.PatientId)
	err := row.Scan(&diagnosis)
	switch err {
	case sql.ErrNoRows:
		return "", nil
	case nil:
		return "Abortion", nil
	default:
		return "", fmt.Errorf("error querying acsis for abortion diagnosis when determining abortive outcome: %+v", err)
	}

}

type pregnancyDiagnosis struct {
	PatientId   int
	EncounterId int
	Date        time.Time
}

func (d *Pregnancies) findPreviousPregnancyDiagnosis(patientId int) (*pregnancyDiagnosis, error) {
	stmt := `
	SELECT e.encounter_id, ed.diagnosis_time
    FROM acsis_adt_encounters e
    INNER JOIN acsis_adt_encounter_diagnoses ed ON e.encounter_id=ed.encounter_id
    WHERE ed.disease_id=32657
    AND e.patient_id=$1
    ORDER BY ed.diagnosis_time DESC
	LIMIT 2;
`
	row := d.AcsisDb.QueryRow(stmt, patientId)
	var pregnancy pregnancyDiagnosis
	err := row.Scan(&pregnancy.EncounterId, &pregnancy.Date)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		pregnancy.PatientId = patientId
		return &pregnancy, nil
	default:
		return nil, fmt.Errorf("error retrieving pregnancy diagnosis from acsis: %+v", err)
	}

}

// FindDiagnosesDuringPregnancy fetches the diagnoses for a patient between the LMP and EDD.
// First it fetches the obstetric details so it can retrieve the pregnancy id. This id is used
// in filtering out the result so that the proper LMP and EDD are used when comparing diagnoses date.
func (d *Pregnancies) FindDiagnosesDuringPregnancy(patientId int) ([]Diagnosis, error) {
	// Retrieve the obstetric details so we can use the current pregnancy's id.
	obs, err := d.FindObstetricDetails(patientId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving current pregnancy's diagnoses from acsis: %+v", err)
	}

	if obs == nil {
		return nil, fmt.Errorf("error: no obstetric details found for current pregnancy in acsis")
	}

	stmt := `SELECT
			aaed.encounter_diagnosis_id,
       		e.patient_id,
       		aai10d.name,
       		aaed.diagnosis_time
		FROM acsis_adt_encounters AS e
		INNER JOIN acsis_adt_encounter_diagnoses aaed on e.encounter_id = aaed.encounter_id
		INNER JOIN acsis_adt_icd10_diseases aai10d on aaed.disease_id = aai10d.disease_id
		INNER JOIN acsis_hc_pregnancies ahp on e.patient_id = ahp.patient_id
		WHERE e.patient_id=$1 AND ahp.pregnancy_id=$2
		AND aaed.diagnosis_time < ahp.estimated_delivery_date
		AND aaed.diagnosis_time > ahp.last_menstrual_period_date
		ORDER BY aaed.diagnosis_time DESC`
	var diagnoses []Diagnosis
	rows, err := d.AcsisDb.Query(stmt, patientId, obs.Id)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("error querying diagnoses before pregnancy from acsis: %+v", err)
	}
	for rows.Next() {
		var diagnosis Diagnosis
		err := rows.Scan(
			&diagnosis.Id,
			&diagnosis.PatientId,
			&diagnosis.Name,
			&diagnosis.Date)
		if err != nil {
			return nil, fmt.Errorf("error scanning diagnosis for patient(%d) %+v", patientId, err)
		}
		diagnoses = append(diagnoses, diagnosis)
	}
	return diagnoses, nil
}

// FindDiagnosesBeforePregnancy returns all diagnoses before the current pregnancy.
// It retrieves the obstetric details as a separate query so it can use the pregnancy id to filter diagnoses
// where the diagnosis time is before the lmp.
func (d *Pregnancies) FindDiagnosesBeforePregnancy(patientId int) ([]Diagnosis, error) {
	// Retrieve the obstetric details so we can use the current pregnancy's id.
	obs, err := d.AcsisDb.FindObstetricDetails(patientId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving diagnoses outside pregnancy from acsis: %+v", err)
	}

	if obs == nil {
		return nil, fmt.Errorf("error: no obstetric details found for current pregnancy in acsis")
	}
	stmt := `SELECT aaed.encounter_diagnosis_id,
       		e.patient_id,
			aai10d.name, 
			aaed.diagnosis_time 
		FROM acsis_adt_encounters AS e
		INNER JOIN acsis_adt_encounter_diagnoses aaed on e.encounter_id = aaed.encounter_id
		INNER JOIN acsis_adt_icd10_diseases aai10d on aaed.disease_id = aai10d.disease_id
		WHERE e.patient_id=$1
		AND aaed.diagnosis_time < (SELECT ahp.last_menstrual_period_date
		      FROM acsis_hc_pregnancies ahp WHERE ahp.pregnancy_id = $2 LIMIT 1)  
		ORDER BY aaed.diagnosis_time DESC`
	var diagnoses []Diagnosis
	rows, err := d.AcsisDb.Query(stmt, patientId, obs.Id)
	if err != nil {
		return nil, fmt.Errorf("error querying diagnoses before pregnancy from acsis: %+v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var diagnosis Diagnosis
		err := rows.Scan(
			&diagnosis.Id,
			&diagnosis.PatientId,
			&diagnosis.Name,
			&diagnosis.Date)
		if err != nil {
			return nil, fmt.Errorf("error scanning diagnosis for patient(%d) %+v", patientId, err)
		}
		diagnoses = append(diagnoses, diagnosis)
	}
	return diagnoses, nil
}

func (d *Pregnancies) FindObstetricHistory(patientId int) ([]ObstetricHistory, error) {
	stmt := `SELECT
				b.birth_id,
				b.mother_id,
				ahbs.name, 
				b.last_modified_time
			FROM acsis_hc_births b
			INNER JOIN acsis_hc_birth_statuses ahbs on b.birth_status_id = ahbs.birth_status_id
			WHERE mother_id=$1`
	var obstetricHistory []ObstetricHistory
	rows, err := d.AcsisDb.Query(stmt, patientId)
	if err != nil {
		return nil, fmt.Errorf("error querying acsis for obstetric history: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var history ObstetricHistory
		err := rows.Scan(&history.Id,
			&history.PatientId,
			&history.ObstetricEvent,
			&history.Date)
		if err != nil {
			return nil, fmt.Errorf("error scanning patient's obstetric history: %+v", err)
		}
		obstetricHistory = append(obstetricHistory, history)
	}
	return obstetricHistory, nil
}
