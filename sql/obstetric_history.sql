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
AND hp.end_time IS NULL;


-- Obstetric Details - Gestational Age
-- This is in the acsis_adt_mch_encounter_details. This means we need to search for an existing mch encounter.
-- Two fields exist for this: by_ultrasound and by_calculation.

-- Obstetric Details - Number of details from acsis_adt_mch_encounter_details (number_of_antenatal_visits)
-- Obstetric Details - Prenatal Care Provider can be extracted from the facility type in the encounter information.
-- Private facilities have a facility_type_id of 14. Everything else is Public.

SELECT
       p.patient_id,
       ed.gestational_age_by_calculation,
       ed.gestational_age_by_ultrasound,
       CASE
           WHEN ft.facility_type_id = 14 THEN 'Private'
           ELSE 'Public'
       END AS care_provider,
       ed.number_of_antenatal_visits,
       e.begin_time as date_of_booking
FROM acsis_hc_patients p
INNER JOIN acsis_adt_encounters e ON p.patient_id=e.patient_id AND e.encounter_type='M'
INNER JOIN acsis_adt_mch_encounter_details ed ON e.encounter_details_id = ed.mch_encounter_details_id
INNER JOIN acsis_hc_facilities f ON e.facility_id = f.facility_id
INNER JOIN acsis_hc_facility_types ft ON f.facility_type_id = ft.facility_type_id
WHERE p.patient_id=591232
AND e.end_time IS NULL
ORDER BY e.begin_time DESC
LIMIT 1;
