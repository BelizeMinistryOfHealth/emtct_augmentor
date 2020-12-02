package db

import (
	"database/sql"
	"fmt"
	"time"

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

func (d *AcsisDb) FindInfantSyphilisTreatment(patientId int) ([]models.Prescription, error) {
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
		WHERE p.patient_id=$1 
		  AND aap.pharmaceutical_id=510
		ORDER BY adep.prescribed_time DESC;
`
	rows, err := d.Query(stmt, patientId)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("error retrieving syphilis treatment for infant from acsis: %+v", err)
	}
	var prescriptions []models.Prescription
	for rows.Next() {
		var p models.Prescription
		var totalDoses sql.NullInt64
		err := rows.Scan(&p.Id, &totalDoses, &p.Pharmaceutical, &p.Frequency, &p.Strength, &p.Comments, &p.PrescribedTime)
		if err != nil {
			return nil, fmt.Errorf("error scanning syphilis prescriptions for infant from acsis: %+v", err)
		}
		p.PatientId = patientId
		if totalDoses.Valid {
			p.TotalDoses = int(totalDoses.Int64)
		}
		prescriptions = append(prescriptions, p)
	}
	return prescriptions, nil
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
			return prescriptions, fmt.Errorf("error scanning syphillis prescription from acsis: %+v", err)
		}
		prescription.PatientId = patientId
		if totalDoses.Valid {
			prescription.TotalDoses = int(totalDoses.Int64)
		}

		prescriptions = append(prescriptions, prescription)
	}

	return prescriptions, nil
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
		return nil, nil
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
