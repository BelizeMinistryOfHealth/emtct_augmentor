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

func (d *AcsisDb) FindHivDiagnoses(patientId int) ([]models.Diagnosis, error) {
	stmt := `SELECT aaed.encounter_diagnosis_id,
       		e.patient_id,
			aai10d.name, 
			aaed.diagnosis_time 
		FROM acsis_adt_encounters AS e
		INNER JOIN acsis_adt_encounter_diagnoses aaed on e.encounter_id = aaed.encounter_id
		INNER JOIN acsis_adt_icd10_diseases aai10d on aaed.disease_id = aai10d.disease_id
		WHERE e.patient_id=$1
		AND  aaed.disease_id IN (473, 474, 475, 476, 477, 9921, 32590, 33195) -- the HIV Test
		ORDER BY aaed.diagnosis_time DESC;`

	var diagnoses []models.Diagnosis
	rows, err := d.Query(stmt, patientId)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("error retrieving the patient's hiv diagnoses from acsis: %+v", err)
	}
	for rows.Next() {
		var d models.Diagnosis
		err := rows.Scan(&d.Id, &d.PatientId, &d.Name, &d.Date)
		if err != nil {
			return nil, fmt.Errorf("error scanning patient's hiv diagnosis from acsis: %+v", err)
		}
		diagnoses = append(diagnoses, d)
	}
	return diagnoses, nil
}

// FindByPatientId searches for a patient who is currently pregnant and is HIV+.
// A patient is considered pregnant if she has a record in the acsis_hc_pregnancies
func (d *AcsisDb) FindByPatientId(id int) (*models.Patient, error) {
	// we only want patients who are currently pregnant or how gave birth no later than 18 months ago.
	//startDate := time.Now().AddDate(-1, -6, 0)
	//dateLimit := startDate.Format(layoutISO)
	stmt := `SELECT p.patient_id, ahp.pregnancy_id,
		l.first_name, l.last_name, l.middle_name,
		p.birth_date, p.ssi_number,p.birth_place, concat(l2.first_name, ' ', l2.last_name) as next_of_kin,
       	ac2.phone1 as next_of_kin_phone,
		ae.name as ethnicity, ahsl.name as education,
		concat(ac.address1, ' ', ac.address2, ',', am.name, ',', aterr.name) as address
		FROM acsis_hc_patients as p
		INNER JOIN acsis_people as l on p.person_id = l.person_id
		INNER JOIN acsis_adt_encounters AS e ON e.patient_id = p.patient_id
		INNER JOIN acsis_adt_encounter_diagnoses ed ON e.encounter_id=ed.encounter_id
		INNER JOIN acsis_contacts ac on l.contact_id = ac.contact_id
		INNER JOIN acsis_municipalities am on ac.municipality_id = am.municipality_id
		INNER JOIN acsis_territories aterr ON ac.territory_id = aterr.territory_id
		INNER JOIN acsis_hc_pregnancies ahp on p.patient_id = ahp.patient_id AND ahp.active IS TRUE
		LEFT JOIN acsis_adt_next_of_kins aanok on p.next_of_kin_id = aanok.next_of_kin_id
		LEFT JOIN acsis_people l2 on aanok.person_id = l2.person_id
		LEFT JOIN acsis_contacts ac2 ON l2.contact_id = ac2.contact_id
		LEFT JOIN acsis_ethnicities ae on p.ethnicity_id = ae.ethnicity_id
		LEFT JOIN acsis_hc_schooling_levels ahsl on p.schooling_level_id = ahsl.schooling_level_id
-- 		WHERE ed.disease_id IN (473, 474, 475, 476, 477, 9921, 32590, 33195) -- the HIV Test
		WHERE p.patient_id = $1
	    ORDER BY ed.diagnosis_time DESC
		LIMIT 1;`

	var patient models.Patient
	row := d.DB.QueryRow(stmt, id)
	var nok sql.NullString
	var nokPhone sql.NullString
	err := row.Scan(&patient.Id,
		&patient.PregnancyId,
		&patient.FirstName,
		&patient.LastName,
		&patient.MiddleName,
		&patient.Dob,
		&patient.Ssn,
		&patient.CountryOfBirth,
		&nok,
		&nokPhone,
		&patient.Ethnicity,
		&patient.Education,
		&patient.Address)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		if nok.Valid {
			patient.NextOfKin = nok.String
		}
		if nokPhone.Valid {
			patient.NextOfKinPhone = nok.String
		}

		return &patient, nil
	default:
		return nil, fmt.Errorf("error querying acsis db for patient: %+v", err)
	}

}

// FindDiagnosesBeforePregnancy returns all diagnoses before the current pregnancy.
// It retrieves the obstetric details as a separate query so it can use the pregnancy id to filter diagnoses
// where the diagnosis time is before the lmp.
func (d *AcsisDb) FindDiagnosesBeforePregnancy(patientId int) ([]models.Diagnosis, error) {
	// Retrieve the obstetric details so we can use the current pregnancy's id.
	obs, err := d.FindObstetricDetails(patientId)
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
	var diagnoses []models.Diagnosis
	rows, err := d.Query(stmt, patientId, obs.Id)
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

// FindDiagnosesDuringPregnancy fetches the diagnoses for a patient between the LMP and EDD.
// First it fetches the obstetric details so it can retrieve the pregnancy id. This id is used
// in filtering out the result so that the proper LMP and EDD are used when comparing diagnoses date.
func (d *AcsisDb) FindDiagnosesDuringPregnancy(patientId int) ([]models.Diagnosis, error) {
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
	var diagnoses []models.Diagnosis
	rows, err := d.Query(stmt, patientId, obs.Id)
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

func (d *AcsisDb) FindLatestAntenatalEncounter(patientId int, lmp *time.Time) (*models.AntenatalEncounter, error) {
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

	var anc models.AntenatalEncounter
	row := d.QueryRow(stmt, patientId)
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

func (d *AcsisDb) FindAntenatalEncounterById(patientId, ancId int) (*models.AntenatalEncounter, error) {
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
        WHERE p.patient_id=$1 AND e.encounter_id=$2
        LIMIT 1;`

	var anc models.AntenatalEncounter
	row := d.QueryRow(stmt, patientId, ancId)
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

// findObstetricPatientDetails retrieves the number of liveborn pregnancies for the patient,
// and previous C/S and planned pregnancies. Not all patients will have data in this table.
// If we do a join of this table when trying to retrieve other obstetric information from the
// pregnancies table, the results are skewed because we can not guarantee that the patient will
// exist in the acsis_hc_obstetric_patient_details_table
func (d *AcsisDb) findObstetricPatientDetails(patientId int) (*models.PregnancyVitals, error) {
	stmt := `SELECT
       ahopd.number_liveborn_pregnancies,
       ahopd.number_caesarean_sections,
       ahopd.previous_pregnancy_planned
FROM acsis_hc_patients p
INNER JOIN acsis_hc_obstetric_patient_details ahopd on p.obstetric_patient_details_id = ahopd.obstetric_patient_details_id
WHERE p.patient_id=$1 LIMIT 1;`

	var vitals models.PregnancyVitals
	row := d.QueryRow(stmt, patientId)
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

// FindObstetricDetails retrieves the patient's LMP, EDD from the acsis_hc_pregnancies
// and the Para/Cs/Planned data as a second query from the acsis_hc_obstetric_patient_details
// table because there is no guarantee that the patient will have data in the latter table.
// When no data is available in that table, we get the wrong results.
func (d *AcsisDb) FindObstetricDetails(patientId int) (*models.PregnancyVitals, error) {
	stmt := `SELECT
       hp.pregnancy_id,
       hp.last_menstrual_period_date,
       hp.estimated_delivery_date
FROM acsis_hc_patients p
INNER JOIN acsis_hc_pregnancies hp ON p.patient_id = hp.patient_id
WHERE p.patient_id=$1
ORDER BY hp.last_modified_time DESC
LIMIT 1;`

	var vitals models.PregnancyVitals
	row := d.QueryRow(stmt, patientId)
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

type PregnancyDiagnosis struct {
	PatientId   int
	EncounterId int
	Date        time.Time
}

func (d *AcsisDb) findPreviousPregnancyDiagnosis(patientId int) (*PregnancyDiagnosis, error) {
	stmt := `
	SELECT e.encounter_id, ed.diagnosis_time
    FROM acsis_adt_encounters e
    INNER JOIN acsis_adt_encounter_diagnoses ed ON e.encounter_id=ed.encounter_id
    WHERE ed.disease_id=32657
    AND e.patient_id=$1
    ORDER BY ed.diagnosis_time DESC
	LIMIT 2;
`
	row := d.QueryRow(stmt, patientId)
	var pregnancy PregnancyDiagnosis
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

// FindCurrentPregnancy returns details about the patient's current pregnancy.
// Under the hood it is 4 separate queries:
// 1. Finds obstetric details. This is information in the acsis_hc_pregnancies table, and contains the LMP and EDD.
// 2. Finds the latest anc encounter.
// 3. Retrieves apgar information
// 4. Retrieves previous diagnoses
// These are all separate queries because the database is not designed in a way to make it possible to retrieve
// all this information using joins. This is partly due to there not being any link between the pregnancies table and
// the encounters table that allows us to retrieve the corresponding anc encounter for a pregnancy.
func (d *AcsisDb) FindCurrentPregnancy(patientId int) (*models.PregnancyVitals, error) {
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
	var vitals models.PregnancyVitals
	var dob *time.Time
	var careProvider string
	var facility string
	row := d.QueryRow(stmt, patientId, anc.Id)
	err = row.Scan(
		&careProvider,
		&facility,
		&vitals.DateOfBooking,
		&vitals.BirthStatus,
		&vitals.ApgarFifthMinute,
		&vitals.ApgarFirstMinute,
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
		vitals.GestationalAge = anc.GestationalAge

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

type testRequestItem struct {
	PatientId              int
	EncounterId            int
	TestRequestId          int
	TestRequestItemId      int
	TestName               string
	TestResult             string
	ReleasedTime           *time.Time
	DateOrderReceivedByLab *time.Time
}

// findCurrentTestRequestItems finds all test requests in a given encounter. This is used when
// searching for a pregnant woman's test results during pregnancy.
func (d *AcsisDb) findCurrentTestRequestItems(patientId, encounterId int) ([]testRequestItem, error) {
	stmt := `SELECT p.patient_id,
                    e.encounter_id,
                    tri.test_request_item_id,
                    tr.test_request_id,
       				tri.released_time,
       				tr.order_received_by_lab_time,
                    t.name
             FROM acsis_hc_patients p
             INNER JOIN acsis_adt_encounters e ON p.patient_id=e.patient_id AND encounter_type='M'
             INNER JOIN acsis_lab_test_requests tr ON tr.encounter_id=e.encounter_id
             INNER JOIN acsis_lab_test_request_items tri ON tr.test_request_id=tri.test_request_id
             INNER JOIN acsis_lab_tests t ON tri.test_id=t.test_id
             WHERE p.patient_id=$1 AND e.encounter_id=$2`
	var testRequests []testRequestItem
	rows, err := d.Query(stmt, patientId, encounterId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving test request items from acsis: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t testRequestItem
		err := rows.Scan(&t.PatientId,
			&t.EncounterId,
			&t.TestRequestItemId,
			&t.TestRequestId,
			&t.ReleasedTime,
			&t.DateOrderReceivedByLab,
			&t.TestName)
		if err != nil {
			return nil, fmt.Errorf("error scanning test request item from acsis: %+v", err)
		}
		testRequests = append(testRequests, t)
	}

	return testRequests, nil
}

type testResult struct {
	Id                int
	PatientId         int
	TestRequestId     int
	TestRequestItemId int
	EncounterId       int
	TestName          string
	TestResult        string
	TestLabel         string
}

func (d *AcsisDb) findTestResults(ri testRequestItem) ([]testResult, error) {
	stmt := `
	SELECT 
	    a.test_result_id,
		altr.test_request_id,
		altri.test_request_item_id,
		e.encounter_id,
		alt.name as test,
		aludli.name as result,
		a.label
	FROM acsis_hc_patients p
		INNER JOIN acsis_adt_encounters AS e ON e.patient_id = p.patient_id
		INNER JOIN acsis_lab_test_requests altr on e.encounter_id = altr.encounter_id AND encounter_type='M'
		INNER JOIN acsis_lab_test_request_items altri on altr.test_request_id = altri.test_request_id
		INNER JOIN acsis_lab_tests alt on altri.test_id = alt.test_id
		INNER JOIN acsis_lab_test_request_results_collected altrrc on altri.test_request_item_id = altrrc.test_request_item_id
		INNER JOIN acsis_lab_test_results a on altrrc.test_result_id = a.test_result_id
		INNER JOIN acsis_lab_user_defined_list_items aludli on altrrc.user_defined_list_value = aludli.user_defined_list_item_id
	WHERE p.patient_id=$1
		AND altr.test_request_id=$2
		AND e.active IS TRUE
	ORDER BY altr.last_modified_time DESC;
`
	var results []testResult
	rows, err := d.Query(stmt, ri.PatientId, ri.TestRequestId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving test results for a test request items from acsis: %+v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var r testResult
		err := rows.Scan(
			&r.Id,
			&r.TestRequestId,
			&r.TestRequestItemId,
			&r.EncounterId,
			&r.TestName,
			&r.TestResult,
			&r.TestLabel)
		if err != nil {
			return nil, fmt.Errorf("error scanning test results when fetching test results from acsis: %+v", err)
		}
		results = append(results, r)
	}
	return results, nil
}

type testSample struct {
	TestSampleId      sql.NullInt32
	CollectedTime     *time.Time
	TestRequestItemId int
	TestRequestId     int
}

func (d *AcsisDb) findTestSamples(tr testRequestItem) (*testSample, error) {
	stmt := `
	SELECT  
		alts.test_sample_id,
		alts.collected_time
	FROM acsis_lab_test_request_specimen_types altrst
		INNER JOIN acsis_lab_test_request_items altri ON altri.test_request_item_id=$1
		LEFT JOIN acsis_lab_test_samples alts ON alts.test_request_specimen_type_id=altrst.test_request_specimen_type_id
	WHERE altrst.test_request_id=$2
	ORDER BY alts.collected_time
	LIMIT 1;
`
	row := d.QueryRow(stmt, tr.TestRequestItemId, tr.TestRequestId)
	var s testSample
	err := row.Scan(&s.TestSampleId, &s.CollectedTime)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		s.TestRequestId = tr.TestRequestId
		s.TestRequestItemId = tr.TestRequestItemId
		return &s, nil
	default:
		return nil, fmt.Errorf("error retrieving lab sample info from acsis: %+v", err)
	}
}

// FindTestsDuringPregnancy returns all the tests conducted during a woman's pregnancy.
// Since we also need the date the sample was collected, the query gets more complicated.
// So we have to issue multiple queries to retrieve separate parts of the information.
// 0. Find the latest anc encounter.
// 1. Find all test request items.
// 2. Find test results for each test request item
// 3. Find the samples for each test request item
// 4. Create the response that will merge the data from all these queries.
func (d *AcsisDb) FindLabTestsDuringPregnancy(patientId int) ([]models.LabResult, error) {
	v, err := d.FindObstetricDetails(patientId)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving lab tests during pregnancy details from acsis: %+v", err)
	}
	anc, err := d.FindLatestAntenatalEncounter(patientId, v.Lmp)
	if err != nil {
		return nil, fmt.Errorf("rerror retrieving antenatal encounter when retrieving lab tests during pregnancy: %+v", err)
	}
	encounterId := anc.Id

	testItems, err := d.findCurrentTestRequestItems(patientId, encounterId)
	if err != nil {
		return nil, fmt.Errorf("error finding current test request items from acsis when retrieving lab tests during pregnancy: %+v", err)
	}
	var labResults []models.LabResult
	for _, t := range testItems {
		testResults, err := d.findTestResults(t)
		if err != nil {
			return nil, fmt.Errorf("error finding test results for a test request item when retrieving lab tests during pregnancy from acsis: %+v", err)
		}
		for _, r := range testResults {
			result := models.LabResult{
				Id:                     r.Id,
				PatientId:              patientId,
				TestName:               fmt.Sprintf("%s - %s", r.TestName, r.TestLabel),
				TestResult:             r.TestResult,
				TestRequestItemId:      r.TestRequestItemId,
				ReleasedTime:           t.ReleasedTime,
				DateOrderReceivedByLab: t.DateOrderReceivedByLab,
				DateSampleTaken:        nil,
				ResultDate:             nil,
			}
			labResults = append(labResults, result)
		}
	}
	var testSamples []testSample
	for _, t := range testItems {
		sample, err := d.findTestSamples(t)
		if err != nil {
			return nil, fmt.Errorf("error finding test samples from when retrieving lab tests during prengnacy from acsis: %+v", err)
		}
		if sample != nil {
			testSamples = append(testSamples, *sample)
		}

	}
	results := assignSamplesToResults(labResults, testSamples)
	return results, nil
}

func assignSamplesToResults(results []models.LabResult, samples []testSample) []models.LabResult {
	for _, s := range samples {
		for _, r := range results {
			if s.TestRequestItemId == r.TestRequestItemId {
				r.DateSampleTaken = s.CollectedTime
			}
		}
	}
	// Deduplicate results
	var r []models.LabResult
	for _, result := range results {
		index := findLabResultIndex(r, result.Id)
		if index < 0 {
			r = append(r, result)
		}
	}
	return r
}

func findLabResultIndex(rs []models.LabResult, id int) int {
	for i, v := range rs {
		if v.Id == id {
			return i
		}
	}
	return -1
}

func (d *AcsisDb) FindPatientSyphilisTreatment(patientId int) ([]models.Prescription, error) {
	v, err := d.FindObstetricDetails(patientId)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving syphilis treatment from acsis: %+v", err)
	}
	anc, err := d.FindLatestAntenatalEncounter(patientId, v.Lmp)
	if err != nil {
		return nil, fmt.Errorf("could not find an antenatal encounter while retrieving syphilis treatment from acsis: %+v", err)
	}
	stmt := `
		SELECT
		    adep.encounter_pharmaceutical_id,
			adep.total_doses,
		   	aap.name as prescription,
		   	acfu.name as frequency, 
			aap.strength || ' ' || aapu.name as strength,
		   	adep.prescribing_physician_special_instructions || ' ' || adep.notes AS comments,
		   	adep.prescribed_time
		FROM acsis_hc_patients p
			INNER JOIN acsis_adt_encounters e ON p.patient_id = e.patient_id
			INNER JOIN acsis_adt_encounter_pharmaceuticals adep ON adep.encounter_id=e.encounter_id
			INNER JOIN acsis_adt_pharmaceuticals aap ON adep.pharmaceutical_id=aap.pharmaceutical_id
			INNER JOIN acsis_coe_frequency_units acfu ON acfu.frequency_unit_id =adep.frequency_unit_id
			INNER JOIN acsis_adt_pharmaceutical_units aapu ON aapu.pharmaceutical_unit_id=aap.strength_unit_id
		WHERE p.patient_id=$1 AND adep.prescribed_time>=$2
		  AND aap.pharmaceutical_id=510
		ORDER BY adep.prescribed_time DESC;
`
	ancDate := anc.BeginDate.Format(layoutISO)
	rows, err := d.Query(stmt, patientId, ancDate)
	if err != nil {
		return nil, fmt.Errorf("error retrieving syphilis from acsis: %+v", err)
	}
	defer rows.Close()
	var prescriptions []models.Prescription
	var totalDoses sql.NullInt64
	for rows.Next() {
		var prescription models.Prescription
		err := rows.Scan(&prescription.Id,
			&totalDoses,
			&prescription.Pharmaceutical,
			&prescription.Frequency,
			&prescription.Strength,
			&prescription.Comments,
			&prescription.PrescribedTime)
		if err != nil {
			return prescriptions, fmt.Errorf("error scanning arv prescription from acsis: %+v", err)
		}
		prescription.PatientId = patientId
		if totalDoses.Valid {
			prescription.TotalDoses = int(totalDoses.Int64)
		}

		prescriptions = append(prescriptions, prescription)
	}

	return prescriptions, nil
}

// FindArvsByPatient finds all ARVs prescribed to a patient during pregnancy.
// It finds the current antenatal encounter, so that it can filter all prescriptions for only that encounter.
func (d *AcsisDb) FindArvsByPatient(patientId int) ([]models.Prescription, error) {
	v, err := d.FindObstetricDetails(patientId)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving arvs treatment from acsis: %+v", err)
	}
	anc, err := d.FindLatestAntenatalEncounter(patientId, v.Lmp)
	if err != nil {
		return nil, fmt.Errorf("could not find an antenatal encounter while retrieving arvs from acsis: %+v", err)
	}
	stmt := `
		SELECT
		    adep.encounter_pharmaceutical_id,
			adep.total_doses,
		   	aap.name as prescription,
		   	acfu.name as frequency, 
			aap.strength || ' ' || aapu.name as strength,
		   	adep.prescribing_physician_special_instructions || ' ' || adep.notes AS comments,
		   	adep.prescribed_time
		FROM acsis_hc_patients p
			INNER JOIN acsis_adt_encounters e ON p.patient_id = e.patient_id
			INNER JOIN acsis_adt_encounter_pharmaceuticals adep ON adep.encounter_id=e.encounter_id
			INNER JOIN acsis_adt_pharmaceuticals aap ON adep.pharmaceutical_id=aap.pharmaceutical_id
			INNER JOIN acsis_coe_frequency_units acfu ON acfu.frequency_unit_id =adep.frequency_unit_id
			INNER JOIN acsis_adt_pharmaceutical_units aapu ON aapu.pharmaceutical_unit_id=aap.strength_unit_id
		WHERE p.patient_id=$1 AND adep.prescribed_time>=$2
			AND (aap.name ILIKE '%Lamivudine%'
			OR aap.name ILIKE '%Zidovudine%'
			OR aap.name ILIKE '%Nevirapine%')
		ORDER BY adep.prescribed_time DESC;
`
	ancDate := anc.BeginDate.Format(layoutISO)
	rows, err := d.Query(stmt, patientId, ancDate)
	if err != nil {
		return nil, fmt.Errorf("error retrieving arvs from acsis: %+v", err)
	}
	defer rows.Close()
	var arvs []models.Prescription
	var totalDoses sql.NullInt64
	for rows.Next() {
		var arv models.Prescription
		err := rows.Scan(&arv.Id,
			&totalDoses,
			&arv.Pharmaceutical,
			&arv.Frequency,
			&arv.Strength,
			&arv.Comments,
			&arv.PrescribedTime)
		if err != nil {
			return arvs, fmt.Errorf("error scanning arv prescription from acsis: %+v", err)
		}
		arv.PatientId = patientId
		if totalDoses.Valid {
			arv.TotalDoses = int(totalDoses.Int64)
		}

		arvs = append(arvs, arv)
	}

	return arvs, nil
}

func (d *AcsisDb) FindPatientBasicInfo(patientId int) (*models.PatientBasicInfo, error) {
	stmt := `
		SELECT
		    hp.patient_id,
			p.first_name,
			p.last_name,
			p.middle_name,
			hp.birth_date,
			hp.ssi_number
		FROM acsis_people p
			INNER JOIN acsis_hc_patients hp ON p.person_id=hp.person_id
		WHERE hp.patient_id=$1;
`
	row := d.QueryRow(stmt, patientId)
	var info models.PatientBasicInfo
	err := row.Scan(&info.Id,
		&info.FirstName,
		&info.LastName,
		&info.MiddleName,
		&info.Dob,
		&info.Ssn)
	if err != nil {
		return nil, fmt.Errorf("error retrieving patient basic info from acsis: %+v", err)
	}
	return &info, nil
}

func (d *AcsisDb) findBirths(motherId int) ([]models.Birth, error) {
	stmt := `
	SELECT b.patient_id, bs.name as birth_status, b.last_modified_time, ahp.birth_date
	FROM acsis_hc_births b
		INNER JOIN acsis_hc_birth_statuses bs ON b.birth_status_id=bs.birth_status_id
		INNER JOIN acsis_hc_patients ahp ON b.patient_id=ahp.patient_id
	WHERE b.mother_id=$1
	ORDER BY b.last_modified_time DESC
`
	rows, err := d.Query(stmt, motherId)
	if err != nil {
		return nil, fmt.Errorf("error fetching births from acsis: %+v", err)
	}
	defer rows.Close()
	var births []models.Birth
	for rows.Next() {
		var b models.Birth
		err := rows.Scan(&b.PatientId, &b.BirthStatus, &b.Date, &b.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("error scanning births from acsis: %+v", err)
		}
		births = append(births, b)
	}
	return births, nil
}

func (d *AcsisDb) FindLatestBirth(motherId, ancId int) (*models.Birth, error) {
	anc, err := d.FindAntenatalEncounterById(motherId, ancId)
	if err != nil {
		return nil, fmt.Errorf("could not find an antenatal encounter while retrieving infant's birth details: %+v", err)
	}
	if anc == nil {
		return nil, nil
	}
	births, err := d.findBirths(motherId)
	if err != nil {
		return nil, fmt.Errorf("could not fetch births while retrieving infant birth details: %+v", err)
	}
	if births == nil || len(births) == 0 {
		return nil, nil
	}
	birth := births[0]
	if !birth.BirthDate.After(*anc.BeginDate) {
		return nil, nil
	}
	return &birth, nil
}

func (d *AcsisDb) FindInfantDiagnoses(infantId int) ([]models.InfantDiagnoses, error) {
	stmt := `
		SELECT
			aed.disease_id,
			e.patient_id,
			aai10d.name as diagnosis,
			aed.notes,
			ap.first_name || ' ' || ap.last_name as doctor,
			aed.diagnosis_time
		FROM acsis_adt_encounters e
			INNER JOIN acsis_adt_encounter_diagnoses aed ON e.encounter_id=aed.encounter_id
			INNER JOIN acsis_adt_icd10_diseases aai10d on aed.disease_id = aai10d.disease_id
			INNER JOIN acsis_hr_staff hs ON aed.doctor_id=hs.staff_id
			INNER JOIN acsis_people ap on hs.person_id = ap.person_id
		WHERE e.patient_id=$1 
		ORDER BY aed.diagnosis_time DESC;
`
	rows, err := d.Query(stmt, infantId)
	if err != nil {
		return nil, fmt.Errorf("error querying for infant diagnoses from acsis: %+v", err)
	}
	defer rows.Close()
	var diagnoses []models.InfantDiagnoses
	for rows.Next() {
		var d models.InfantDiagnoses
		err := rows.Scan(&d.DiagnosisId,
			&d.PatientId,
			&d.Diagnosis,
			&d.Comments,
			&d.Doctor,
			&d.Date)
		if err != nil {
			return nil, fmt.Errorf("error scanning infant diagnosis: %+v", err)
		}
		diagnoses = append(diagnoses, d)
	}

	return diagnoses, nil
}

// FindPregnancyInfant returns the infant that was born after the mother's current LMP.
func (d *AcsisDb) FindPregnancyInfant(motherId int) (*models.Infant, error) {
	stmt := `
	SELECT 
	       b.patient_id,
	       ppl.first_name,
	       ppl.middle_name,
	       ppl.last_name,
	       pt.birth_date,
	       mppl.first_name as mfirst_name,
	       mppl.middle_name as mmiddle_name,
	       mppl.last_name as mlast_name,
	       mpt.birth_date as mdob
    FROM acsis_hc_births b
	INNER JOIN acsis_hc_patients pt ON pt.patient_id=b.patient_id
	INNER JOIN acsis_people ppl ON pt.person_id = ppl.person_id
	INNER JOIN acsis_hc_patients mpt ON b.mother_id=mpt.patient_id
	INNER JOIN acsis_people mppl ON mppl.person_id=mpt.person_id
	WHERE b.mother_id=$1
	ORDER BY pt.birth_date DESC
	LIMIT 1;
`
	var infant models.Infant
	row := d.QueryRow(stmt, motherId)
	err := row.Scan(
		&infant.Infant.PatientId,
		&infant.Infant.FirstName,
		&infant.Infant.MiddleName,
		&infant.Infant.LastName,
		&infant.Infant.Dob,
		&infant.Mother.FirstName,
		&infant.Mother.MiddleName,
		&infant.Mother.LastName,
		&infant.Mother.Dob)
	if err != nil {
		return nil, fmt.Errorf("error querying infant basic information from acsis: %+v", err)
	}
	infant.Mother.PatientId = motherId

	// Compare the date of birth with the mother's LMP.
	// If the date is not after LMP, then this birth does not belong to the latest pregnancy.
	obs, err := d.FindObstetricDetails(motherId)
	if err != nil {
		return nil, fmt.Errorf("error fetching infant details: %+v", err)
	}
	if obs.Lmp.After(*infant.Infant.Dob) {
		return nil, fmt.Errorf("no infant found in acsis for the current pregnancy: %+v", err)
	}

	return &infant, nil
}

func (d *AcsisDb) FindInfant(infantId int) (*models.Infant, error) {
	stmt := `
	SELECT 
	       b.patient_id,
	       ppl.first_name,
	       ppl.middle_name,
	       ppl.last_name,
	       pt.birth_date,
	       mppl.first_name as mfirst_name,
	       mppl.middle_name as mmiddle_name,
	       mppl.last_name as mlast_name,
	       mpt.birth_date as mdob,
	       mpt.patient_id as mother_id
    FROM acsis_hc_births b
	INNER JOIN acsis_hc_patients pt ON pt.patient_id=b.patient_id
	INNER JOIN acsis_people ppl ON pt.person_id = ppl.person_id
	INNER JOIN acsis_hc_patients mpt ON b.mother_id=mpt.patient_id
	INNER JOIN acsis_people mppl ON mppl.person_id=mpt.person_id
	WHERE pt.patient_id=$1
	ORDER BY pt.birth_date DESC
	LIMIT 1;
`
	var infant models.Infant
	row := d.QueryRow(stmt, infantId)
	err := row.Scan(
		&infant.Infant.PatientId,
		&infant.Infant.FirstName,
		&infant.Infant.MiddleName,
		&infant.Infant.LastName,
		&infant.Infant.Dob,
		&infant.Mother.FirstName,
		&infant.Mother.MiddleName,
		&infant.Mother.LastName,
		&infant.Mother.Dob,
		&infant.Mother.PatientId)
	if err != nil {
		return nil, fmt.Errorf("error querying infant basic information from acsis: %+v", err)
	}

	return &infant, nil
}
