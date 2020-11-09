SELECT p.patient_id,
       ppl.first_name || ' ' || ppl.last_name as patient,
       ahp.last_menstrual_period_date, ahp.estimated_delivery_date
FROM acsis_hc_patients p
INNER JOIN acsis_people ppl ON p.person_id=ppl.person_id AND ppl.gender_id =2
INNER JOIN acsis_hc_pregnancies ahp ON ahp.patient_id=p.patient_id AND ahp.estimated_delivery_date > '2020-12-01'
INNER JOIN acsis_adt_encounters e ON p.patient_id=e.patient_id
INNER JOIN acsis_adt_encounter_diagnoses ed ON e.encounter_id=ed.encounter_id
WHERE
 p.patient_id=64628
  AND ed.disease_id IN (473, 474, 475, 476, 477, 9921, 32590);


SELECT p.patient_id, ahp.pregnancy_id, ed.diagnosis_time,
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
		WHERE ed.disease_id IN (473, 474, 475, 476, 477, 9921, 32590) -- the HIV Test
		AND p.patient_id = 64628
	    ORDER BY ed.diagnosis_time DESC LIMIT 1;



SELECT aaed.encounter_diagnosis_id,
       		e.patient_id,
			aai10d.name,
			aaed.diagnosis_time
		FROM acsis_adt_encounters AS e
		INNER JOIN acsis_adt_encounter_diagnoses aaed on e.encounter_id = aaed.encounter_id
		INNER JOIN acsis_adt_icd10_diseases aai10d on aaed.disease_id = aai10d.disease_id
		WHERE e.patient_id=64628
		AND aaed.diagnosis_time < (SELECT ahp.last_menstrual_period_date
		      FROM acsis_hc_pregnancies ahp WHERE ahp.patient_id = e.patient_id ORDER BY
		      ahp.last_menstrual_period_date DESC LIMIT 1)
		ORDER BY aaed.diagnosis_time DESC;


SELECT p.patient_id,
       ppl.first_name || ' ' || ppl.last_name as patient,
       ahp.last_menstrual_period_date, ahp.estimated_delivery_date
FROM acsis_hc_patients p
INNER JOIN acsis_people ppl ON p.person_id=ppl.person_id AND ppl.gender_id =2
INNER JOIN acsis_hc_pregnancies ahp ON ahp.patient_id=p.patient_id AND ahp.estimated_delivery_date > '2020-07-01'
INNER JOIN acsis_adt_encounters e ON p.patient_id=e.patient_id
INNER JOIN acsis_adt_encounter_diagnoses ed ON e.encounter_id=ed.encounter_id
WHERE
      e.active IS TRUE
      AND p.active IS TRUE
      AND ed.active IS TRUE
  AND ed.disease_id IN (473, 474, 475, 476, 477, 9921, 32590);

SELECT disease_id, code, name
FROM acsis_adt_icd10_diseases
WHERE name ILIKE '%hiv%';
