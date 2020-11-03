-- Obstetric History
SELECT ahbs.name, b.last_modified_time
FROM acsis_hc_births b
INNER JOIN acsis_hc_birth_statuses ahbs on b.birth_status_id = ahbs.birth_status_id
WHERE mother_id=641336;

-- Obstetric Details - PARA
SELECT opd.number_liveborn_pregnancies
FROM acsis_hc_patients p
INNER JOIN acsis_hc_obstetric_patient_details opd ON p.obstetric_patient_details_id = opd.obstetric_patient_details_id
WHERE patient_id = 591232;


-- Obstetric Details - Planned pregnancy is only captured for the
-- previous pregnancy. This is in the acsis_hc_obstetric_patient_details table.
-- The field is previous_pregnancy_planned.

-- Obstetric Details - CS -- from acsis_hc_obstetric_patient_details: number_caesarean_sections

SELECT
       p.patient_id,
       hp.pregnancy_id,
       ahopd.number_liveborn_pregnancies,
       ahopd.number_aborted_pregnancies,
       ahopd.number_caesarean_sections,
       ahopd.previous_pregnancy_planned,
       hp.last_menstrual_period_date,
       hp.estimated_delivery_date,
       hp.last_modified_time
FROM acsis_hc_patients p
INNER JOIN acsis_hc_pregnancies hp ON p.patient_id = hp.patient_id
INNER JOIN acsis_hc_obstetric_patient_details ahopd on p.obstetric_patient_details_id = ahopd.obstetric_patient_details_id
WHERE p.patient_id=591232
ORDER BY hp.last_menstrual_period_date DESC
LIMIT 1;


-- Obstetric Details - Find the applicable mch encounter
-- This will be an encounter where the lmp is less than the mch begin date
SELECT e.encounter_id, amed.mch_encounter_details_id, amed.estimated_delivery_date,
       e.begin_time
FROM acsis_hc_patients p
INNER JOIN acsis_adt_encounters e ON p.patient_id=e.patient_id
INNER JOIN acsis_adt_mch_encounter_details amed ON e.encounter_details_id=amed.mch_encounter_details_id
WHERE p.patient_id=591232
ORDER BY e.begin_time DESC
LIMIT 1;

-- Obstetric Details - Gestational Age
-- This is in the acsis_adt_mch_encounter_details. This means we need to search for an existing mch encounter.
-- Two fields exist for this: by_ultrasound and by_calculation.

-- Obstetric Details - Number of details from acsis_adt_mch_encounter_details (number_of_antenatal_visits)
-- Obstetric Details - Prenatal Care Provider can be extracted from the facility type in the encounter information.
-- Private facilities have a facility_type_id of 14. Everything else is Public.

SELECT
       e.encounter_id,
       p.patient_id,
       e2.begin_time,
       pre.last_menstrual_period_date,
       ed.gestational_age_by_calculation,
       ed.gestational_age_by_ultrasound,
       CASE
           WHEN ft.facility_type_id = 14 THEN 'Private'
           ELSE 'Public'
       END AS care_provider,
       ed.number_of_antenatal_visits,
       e.begin_time as date_of_booking,
--        bs.name AS birth_status,
       ahipd.apgar_fifth_minute,
       ahipd.apgar_first_minute,
       ahdt.name AS delivery_type
FROM acsis_hc_patients p
INNER JOIN acsis_adt_encounters e ON p.patient_id=e.patient_id AND e.encounter_type='M' AND e.active IS TRUE
INNER JOIN acsis_adt_mch_encounter_details ed ON e.encounter_details_id = ed.mch_encounter_details_id
INNER JOIN acsis_hc_facilities f ON e.facility_id = f.facility_id
INNER JOIN acsis_hc_facility_types ft ON f.facility_type_id = ft.facility_type_id
INNER JOIN acsis_hc_pregnancies AS pre ON pre.patient_id=p.patient_id AND pre.active IS TRUE
LEFT JOIN acsis_adt_encounters e2 ON p.patient_id=e2.patient_id AND e.encounter_type='L'
-- LEFT JOIN acsis_adt_labour_encounter_details aaled ON e2.encounter_id = aaled.mch_encounter_id
-- LEFT JOIN acsis_hc_birth_statuses bs ON aaled.birth_status_id=bs.birth_status_id
LEFT JOIN acsis_hc_infant_patient_details ahipd ON p.patient_id = ahipd.mothers_patient_id AND ahipd.labour_encounter_id =e2.encounter_id
LEFT JOIN acsis_hc_delivery_types ahdt ON ahipd.delivery_type_id=ahdt.delivery_type_id
WHERE p.patient_id=591232
AND e.end_time IS NULL -- this won't work when we want to follow a pregnancy after birth.
ORDER BY e.begin_time DESC;
-- LIMIT 1;


-- The child's apgar can only be retrieved after birth on a separate query
-- We need to search for a labour encounter and get the most recent one.
-- Then check that the date is after the LMP and retrieve the apgar information
-- along with the gestational age.
SELECT
       p.patient_id,
       e2.begin_time,
       icd.name as icd10,
       bs.name AS birth_status,
       ahipd.apgar_fifth_minute,
       ahipd.gestational_age,
       ahipd.apgar_first_minute
FROM acsis_hc_patients p
INNER JOIN acsis_adt_encounters e2 ON p.patient_id=e2.patient_id AND e2.encounter_type='L'
INNER JOIN acsis_adt_labour_encounter_details aaled ON aaled.labour_encounter_details_id = e2.encounter_details_id
INNER JOIN acsis_hc_birth_statuses bs ON aaled.birth_status_id=bs.birth_status_id
INNER JOIN acsis_hc_infant_patient_details ahipd ON p.patient_id = ahipd.mothers_patient_id AND ahipd.labour_encounter_id =e2.encounter_id
INNER JOIN acsis_adt_encounter_encounter_diagnoses as eed ON eed.encounter_id = e2.encounter_id
LEFT JOIN acsis_adt_icd10_diseases as icd ON eed.encounter_diagnosis_id = icd.disease_id
WHERE p.patient_id=591232
ORDER BY e2.begin_time DESC;


-- When the outcome is not a live birth nor a still birth, then we need to search
-- for all diagnoses where there is an abortive outcome.
SELECT e.encounter_id, aai10d.name, aaed.diagnosis_time, e.begin_time
FROM acsis_adt_encounters AS e
INNER JOIN acsis_adt_encounter_diagnoses aaed on e.encounter_id = aaed.encounter_id
INNER JOIN acsis_adt_icd10_diseases aai10d on aaed.disease_id = aai10d.disease_id
INNER JOIN acsis_hc_pregnancies ahp on e.patient_id = ahp.patient_id
-- WHERE e.patient_id=589780005
WHERE e.begin_time > '2020-01-01' AND e.begin_time < '2020-11-03'
AND aaed.diagnosis_time < ahp.last_menstrual_period_date
-- Confirm with Dr. Beer that these are correct codes. Best to do by showing her
-- in a meeting what those codes look like.
AND (aai10d.code ILIKE 'O06%' OR aai10d.code ILIKE 'O03%' OR aai10d.code ILIKE 'O05%'
OR aai10d.code ILIKE 'O04%')
ORDER BY aaed.diagnosis_time DESC;




SELECT
       CASE
           WHEN ft.facility_type_id = 14 THEN 'Private'
           ELSE 'Public'
       END AS care_provider,
       e.begin_time as date_of_booking,
       bs.name AS birth_status,
       ahipd.apgar_fifth_minute,
       ahipd.apgar_first_minute
FROM acsis_hc_patients p
INNER JOIN acsis_adt_encounters e ON p.patient_id=e.patient_id AND e.encounter_details_id=8450564
INNER JOIN acsis_hc_facilities f ON e.facility_id = f.facility_id
INNER JOIN acsis_hc_facility_types ft ON f.facility_type_id = ft.facility_type_id
LEFT JOIN acsis_adt_labour_encounter_details aaled ON e.encounter_id = aaled.mch_encounter_id
LEFT JOIN acsis_hc_birth_statuses bs ON aaled.birth_status_id=bs.birth_status_id
LEFT JOIN acsis_hc_infant_patient_details ahipd ON p.patient_id = ahipd.mothers_patient_id
WHERE p.patient_id=591232
ORDER BY e.begin_time DESC
LIMIT 1;

SELECT * FROM acsis_adt_encounters WHERE patient_id=591232 AND encounter_type='M' ORDER BY begin_time DESC ;