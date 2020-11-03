SELECT * FROM acsis_hc_patients WHERE currently_pregnant IS TRUE;

-- patient id for dev: 730875

-- ALL HIV+
SELECT * FROM acsis_hc_patient_chronic_diseases
WHERE chronic_disease_id IN (12, 13)
ORDER BY last_modified_time DESC;

-- FIND pregnant women who are HIV+
SELECT p.patient_id, p.birth_date, l.first_name, l.last_name
FROM acsis_hc_patients as p
INNER JOIN acsis_people as l on p.person_id = l.person_id
INNER JOIN acsis_hc_patient_chronic_diseases AS c ON c.patient_id = p.patient_id
INNER JOIN acsis_hc_chronic_diseases ahcd on c.chronic_disease_id = ahcd.chronic_disease_id
WHERE p.currently_pregnant IS TRUE
AND c.chronic_disease_id IN (12, 13);

-- FIND all HIV Confirmation tests for a patient
SELECT altr.test_request_id, e.encounter_id, alt.name,
    a.label, aludli.name as list_item, altri.released_time,
       l.first_name, l.last_name, l.middle_name, l.maiden_name,
       p.birth_date, p.ssi_number,p.birth_place, concat(l2.first_name, ' ', l2.last_name) as next_of_kin,
       ac2.phone1 as next_of_kin_phone,
       ae.name as ethnicity, ahsl.name as education ,p.currently_pregnant, p.last_menstrual_period_date,
       concat(ac.address1, ac.address2, ',', am.name, ',', aterr.name) as address
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
LEFT JOIN acsis_adt_next_of_kins aanok on p.next_of_kin_id = aanok.next_of_kin_id
LEFT JOIN acsis_people l2 on aanok.person_id = l2.person_id
LEFT JOIN acsis_contacts ac2 ON l2.contact_id = ac2.contact_id
LEFT JOIN acsis_ethnicities ae on p.ethnicity_id = ae.ethnicity_id
LEFT JOIN acsis_hc_schooling_levels ahsl on p.schooling_level_id = ahsl.schooling_level_id
INNER JOIN acsis_hc_pregnancies ahp on p.patient_id = ahp.patient_id AND ahp.active IS TRUE
WHERE altri.test_id IN (2) -- the HIV Test
  AND altrrc.user_defined_list_value IS NOT NULL
--   AND ahp.last_menstrual_period_date > '2019-10-10'
  -- user_defined_list_value == 2 filter out positive results
  -- AND altrrc.user_defined_list_value = 2
  -- test_result_id 348 is the HIV Confirmation result
  AND a.test_result_id = 348
  AND p.patient_id = 591232
ORDER BY released_time DESC;


-- Pregnancies
SELECT * FROM acsis_hc_pregnancies
WHERE last_modified_time>'2020-01-01' AND active IS TRUE AND estimated_delivery_date > '2020-11-01';


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

-- Obstetric Details - Gestational Age
-- This is in the acsis_adt_mch_encounter_details. This means we need to search for an existing mch encounter.
-- Two fields exist for this: by_ultrasound and by_calculation.

-- Obstetric Details - Planned pregnancy is only captured for the
-- previous pregnancy. This is in the acsis_hc_obstetric_patient_details table.
-- The field is previous_pregnancy_planned.

-- Obstetric Details - Number of details from acsis_adt_mch_encounter_details (number_of_antenatal_visits)
-- Obstetric Details - Prenatal Care Provider can be extracted from the facility type in the encounter information.
-- Private facilities have a facility_type_id of 14. Everything else is Public.

-- Obstetric Details - CS -- from acsis_hc_obstetric_patient_details: number_caesarean_sections

-- Illnesses before Pregnancy
SELECT aaed.encounter_diagnosis_id,
       		e.patient_id,
			aai10d.name,
			aaed.diagnosis_time
		FROM acsis_adt_encounters AS e
		INNER JOIN acsis_adt_encounter_diagnoses aaed on e.encounter_id = aaed.encounter_id
		INNER JOIN acsis_adt_icd10_diseases aai10d on aaed.disease_id = aai10d.disease_id
		WHERE e.patient_id=591232
		  AND aaed.diagnosis_time < (SELECT ahp.last_menstrual_period_date
		      FROM acsis_hc_pregnancies ahp WHERE ahp.patient_id = e.patient_id ORDER BY
		      ahp.last_menstrual_period_date DESC LIMIT 1)
		ORDER BY aaed.diagnosis_time DESC;

-- Illnesses During Pregnancy. Requires pregnancy_id
SELECT e.encounter_id, aai10d.name, aaed.diagnosis_time
FROM acsis_adt_encounters AS e
INNER JOIN acsis_adt_encounter_diagnoses aaed on e.encounter_id = aaed.encounter_id
INNER JOIN acsis_adt_icd10_diseases aai10d on aaed.disease_id = aai10d.disease_id
WHERE e.patient_id=589780005
AND aaed.diagnosis_time < (SELECT ahp.last_menstrual_period_date
		      FROM acsis_hc_pregnancies ahp WHERE ahp.patient_id = e.patient_id ORDER BY
		      ahp.estimated_delivery_date DESC LIMIT 1)
AND aaed.diagnosis_time > (SELECT ahp.last_menstrual_period_date
		      FROM acsis_hc_pregnancies ahp WHERE ahp.patient_id = e.patient_id ORDER BY
		      ahp.last_menstrual_period_date DESC LIMIT 1)
ORDER BY aaed.diagnosis_time DESC;

