package db

import (
	"database/sql"
	"fmt"
	"time"

	age "github.com/bearbin/go-age"
	"moh.gov.bz/mch/emtct/internal/config"
	"moh.gov.bz/mch/emtct/internal/models"
)

type AcsisDb struct {
	*sql.DB
}

const (
	layoutISO = "2006-01-02"
)

func NewAcsisConnection(cnf *config.DbConf) (*AcsisDb, error) {
	connstr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", cnf.Username, cnf.Database, cnf.Password, cnf.Host)
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}
	return &AcsisDb{db}, nil
}

// FindByPatientId searches for a patient who is currently pregnant and is HIV+.
// A patient is considered pregnant if she has a record in the acsis_hc_pregnancies
func (d *AcsisDb) FindByPatientId(id int) (*models.Patient, error) {
	// we only want patients who are currently pregnant or how gave birth no later than 18 months ago.
	startDate := time.Now().AddDate(-1, -6, 0)
	dateLimit := startDate.Format(layoutISO)
	stmt := `SELECT p.patient_id, ahp.pregnancy_id, altri.released_time,
		l.first_name, l.last_name, l.middle_name,
		p.birth_date, p.ssi_number,p.birth_place, concat(l2.first_name, ' ', l2.last_name) as next_of_kin,
       	ac2.phone1 as next_of_kin_phone,
		ae.name as ethnicity, ahsl.name as education,
		concat(ac.address1, ' ', ac.address2, ',', am.name, ',', aterr.name) as address
		FROM acsis_hc_patients as p
		INNER JOIN acsis_people as l on p.person_id = l.person_id
		INNER JOIN acsis_adt_encounters AS e ON e.patient_id = p.patient_id
		INNER JOIN acsis_lab_test_requests altr on e.encounter_id = altr.encounter_id
		INNER JOIN acsis_lab_test_request_items altri on altr.test_request_id = altri.test_request_id
		INNER JOIN acsis_lab_tests alt on altri.test_id = alt.test_id
		INNER JOIN acsis_lab_test_request_results_collected altrrc on altri.test_request_item_id = altrrc.test_request_item_id
		INNER JOIN acsis_lab_test_results a on altrrc.test_result_id = a.test_id AND a.test_id IN (2, 152, 5015, 5033, 5032)
		INNER JOIN acsis_lab_user_defined_list_items aludli on altrrc.user_defined_list_value = aludli.user_defined_list_item_id
		INNER JOIN acsis_contacts ac on l.contact_id = ac.contact_id
		INNER JOIN acsis_municipalities am on ac.municipality_id = am.municipality_id
		INNER JOIN acsis_territories aterr ON ac.territory_id = aterr.territory_id
		INNER JOIN acsis_hc_pregnancies ahp on p.patient_id = ahp.patient_id AND ahp.active IS TRUE
		LEFT JOIN acsis_adt_next_of_kins aanok on p.next_of_kin_id = aanok.next_of_kin_id
		LEFT JOIN acsis_people l2 on aanok.person_id = l2.person_id
		LEFT JOIN acsis_contacts ac2 ON l2.contact_id = ac2.contact_id
		LEFT JOIN acsis_ethnicities ae on p.ethnicity_id = ae.ethnicity_id
		LEFT JOIN acsis_hc_schooling_levels ahsl on p.schooling_level_id = ahsl.schooling_level_id
		WHERE altri.test_id IN (2) -- the HIV Test
		AND altrrc.user_defined_list_value IS NOT NULL
		AND ahp.last_menstrual_period_date > $2
		AND a.test_result_id = 348
		AND p.patient_id = $1
		ORDER BY released_time DESC LIMIT 1;`

	var patient models.Patient
	row := d.DB.QueryRow(stmt, id, dateLimit)
	err := row.Scan(&patient.Id,
		&patient.PregnancyId,
		&patient.HivDiagnosisDate,
		&patient.FirstName,
		&patient.LastName,
		&patient.MiddleName,
		&patient.Dob,
		&patient.Ssn,
		&patient.CountryOfBirth,
		&patient.NextOfKin,
		&patient.NextOfKinPhone,
		&patient.Ethnicity,
		&patient.Education,
		&patient.Address)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		patient.Hiv = true
		return &patient, nil
	default:
		return nil, fmt.Errorf("error querying acsis db for patient: %+v", err)
	}

}

func (d *AcsisDb) FindDiagnosesBeforePregnancy(patientId int) ([]models.Diagnosis, error) {
	stmt := `SELECT aaed.encounter_diagnosis_id,
       		e.patient_id,
			aai10d.name, 
			aaed.diagnosis_time 
		FROM acsis_adt_encounters AS e
		INNER JOIN acsis_adt_encounter_diagnoses aaed on e.encounter_id = aaed.encounter_id
		INNER JOIN acsis_adt_icd10_diseases aai10d on aaed.disease_id = aai10d.disease_id
		WHERE e.patient_id=$1
		AND aaed.diagnosis_time < (SELECT ahp.last_menstrual_period_date
		      FROM acsis_hc_pregnancies ahp WHERE ahp.patient_id = e.patient_id ORDER BY
		      ahp.last_menstrual_period_date DESC LIMIT 1)  
		ORDER BY aaed.diagnosis_time DESC`
	var diagnoses []models.Diagnosis
	rows, err := d.Query(stmt, patientId)
	if err != nil {
		return nil, fmt.Errorf("error querying diagnoses before pregnancy from acsis: %+v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var diagnosis models.Diagnosis
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

func (d *AcsisDb) FindDiagnosesDuringPregnancy(patientId int) ([]models.Diagnosis, error) {
	stmt := `SELECT
			aaed.encounter_diagnosis_id,
       		e.patient_id,
       		aai10d.name,
       		aaed.diagnosis_time
		FROM acsis_adt_encounters AS e
		INNER JOIN acsis_adt_encounter_diagnoses aaed on e.encounter_id = aaed.encounter_id
		INNER JOIN acsis_adt_icd10_diseases aai10d on aaed.disease_id = aai10d.disease_id
		INNER JOIN acsis_hc_pregnancies ahp on e.patient_id = ahp.patient_id
		WHERE e.patient_id=$1
		AND aaed.diagnosis_time < ahp.estimated_delivery_date
		AND aaed.diagnosis_time > ahp.last_menstrual_period_date
		ORDER BY aaed.diagnosis_time DESC`
	var diagnoses []models.Diagnosis
	rows, err := d.Query(stmt, patientId)
	if err != nil {
		return nil, fmt.Errorf("error querying diagnoses before pregnancy from acsis: %+v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var diagnosis models.Diagnosis
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

func (d *AcsisDb) FindObstetricHistory(patientId int) ([]models.ObstetricHistory, error) {
	stmt := `SELECT
				b.birth_id,
				b.mother_id,
				ahbs.name, 
				b.last_modified_time
			FROM acsis_hc_births b
			INNER JOIN acsis_hc_birth_statuses ahbs on b.birth_status_id = ahbs.birth_status_id
			WHERE mother_id=$1`
	var obstetricHistory []models.ObstetricHistory
	rows, err := d.Query(stmt, patientId)
	if err != nil {
		return nil, fmt.Errorf("error querying acsis for obstetric history: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var history models.ObstetricHistory
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

func (d *AcsisDb) findLatestAntenatalEncounter(patientId int) (*models.AntenatalEncounter, error) {
	stmt := `SELECT e.encounter_id,
           e.patient_id,
           amed.mch_encounter_details_id, 
           amed.estimated_delivery_date,
           e.begin_time,
       		COALESCE(amed.gestational_age_by_calculation, amed.gestational_age_by_ultrasound) AS gestational_age,
           amed.number_of_antenatal_visits
        FROM acsis_hc_patients p
        INNER JOIN acsis_adt_encounters e ON p.patient_id=e.patient_id
        INNER JOIN acsis_adt_mch_encounter_details amed ON e.encounter_details_id=amed.mch_encounter_details_id
        WHERE p.patient_id=$1
        ORDER BY e.begin_time DESC
        LIMIT 1;`

	var anc models.AntenatalEncounter
	row := d.QueryRow(stmt, patientId)
	err := row.Scan(&anc.Id,
		&anc.PatientId,
		&anc.MchEncounterDetailsId,
		&anc.EstimatedDeliveryDate,
		&anc.BeginDate,
		&anc.GestationalAge,
		&anc.NumberAntenatalVisits)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &anc, nil
	default:
		return nil, fmt.Errorf("error querying for mch details from acsis: %+v", err)
	}
}

func (d *AcsisDb) findObstetricDetails(patientId int) (*models.PregnancyVitals, error) {
	stmt := `SELECT
       hp.pregnancy_id,
       ahopd.number_liveborn_pregnancies,
       ahopd.number_caesarean_sections,
       ahopd.previous_pregnancy_planned,
       hp.last_menstrual_period_date,
       hp.estimated_delivery_date
FROM acsis_hc_patients p
INNER JOIN acsis_hc_pregnancies hp ON p.patient_id = hp.patient_id
INNER JOIN acsis_hc_obstetric_patient_details ahopd on p.obstetric_patient_details_id = ahopd.obstetric_patient_details_id
WHERE p.patient_id=$1
ORDER BY hp.last_menstrual_period_date DESC
LIMIT 1;`

	var vitals models.PregnancyVitals
	row := d.QueryRow(stmt, patientId)
	err := row.Scan(&vitals.Id,
		&vitals.Para,
		&vitals.Cs,
		&vitals.Planned,
		&vitals.Lmp,
		&vitals.Edd)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &vitals, nil
	default:
		return nil, fmt.Errorf("error querying obstetric details from acsis: %+v", err)
	}

}

func (d *AcsisDb) FindCurrentPregnancy(patientId int) (*models.PregnancyVitals, error) {
	// Find the anc encounter. This is the most recent anc encounter in patient's docket.
	anc, err := d.findLatestAntenatalEncounter(patientId)
	if err != nil {
		return nil, fmt.Errorf("could not find current pregnancy details because no antenatal encounter was found: %+v", err)
	}

	stmt := `SELECT
       CASE
           WHEN ft.facility_type_id = 14 THEN 'Private'
           ELSE 'Public'
       END AS care_provider,
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
	var vitals models.PregnancyVitals
	var dob *time.Time
	row := d.QueryRow(stmt, patientId, anc.Id)
	err = row.Scan(
		&vitals.PrenatalCareProvider,
		&vitals.DateOfBooking,
		&vitals.BirthStatus,
		&vitals.ApgarFifthMinute,
		&vitals.ApgarFirstMinute,
		&dob,
	)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		v, err := d.findObstetricDetails(patientId)
		if err != nil {
			return nil, fmt.Errorf("error while retrieving pregnancy details from acsis: %+v", err)
		}
		ageAtLmp := age.AgeAt(*dob, v.Lmp)
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

		return &vitals, nil
	default:
		return nil, fmt.Errorf("error while retrieving pregnancy details from acsis: %+v", err)
	}

}

func (d *AcsisDb) abortiveOutcome(v models.PregnancyVitals) (string, error) {
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
	stmt := `SELECT aai10d.name
            FROM acsis_adt_encounters AS e
            INNER JOIN acsis_adt_encounter_diagnoses aaed on e.encounter_id = aaed.encounter_id
            INNER JOIN acsis_adt_icd10_diseases aai10d on aaed.disease_id = aai10d.disease_id
            INNER JOIN acsis_hc_pregnancies ahp on e.patient_id = ahp.patient_id
            WHERE e.patient_id = $1 
            AND e.begin_time > '2020-01-01' AND e.begin_time < '2020-11-03'
            AND aaed.diagnosis_time < ahp.last_menstrual_period_date
            AND (aai10d.code ILIKE 'O06%' OR aai10d.code ILIKE 'O03%' OR aai10d.code ILIKE 'O05%'
            OR aai10d.code ILIKE 'O04%')
            ORDER BY aaed.diagnosis_time DESC
            LIMIT 1;`
	var diagnosis string
	row := d.QueryRow(stmt, v.PatientId)
	err := row.Scan(&diagnosis)
	if err != nil {
		return "", fmt.Errorf("error querying acsis for abortion diagnosis when determining abortive outcome: %+v", err)
	}
	return "Abortion", nil
}
