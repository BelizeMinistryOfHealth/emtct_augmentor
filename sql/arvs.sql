SELECT adep.total_doses,
       aap.name as prescription,
       acfu.name, aap.strength || ' ' || aapu.name as strength,
       adep.prescribing_physician_special_instructions || ' ' || adep.notes AS comments,
       adep.prescribed_time

FROM acsis_hc_patients p
INNER JOIN acsis_adt_encounters e ON p.patient_id = e.patient_id
INNER JOIN acsis_adt_encounter_pharmaceuticals adep ON adep.encounter_id=e.encounter_id
INNER JOIN acsis_adt_pharmaceuticals aap ON adep.pharmaceutical_id=aap.pharmaceutical_id
INNER JOIN acsis_coe_frequency_units acfu ON acfu.frequency_unit_id =adep.frequency_unit_id
INNER JOIN acsis_adt_pharmaceutical_units aapu ON aapu.pharmaceutical_unit_id=aap.strength_unit_id
WHERE p.patient_id=591232
AND adep.prescribed_time>'2020-06-19'
ORDER BY adep.prescribed_time DESC;


SELECT p.first_name,
       p.last_name,
       p.middle_name,
       hp.birth_date,
       hp.ssi_number
FROM acsis_people p
INNER JOIN acsis_hc_patients hp ON p.person_id=hp.person_id
WHERE hp.patient_id=591232;

SELECT encounter_id, begin_time
FROM acsis_adt_encounters
WHERE patient_id=591232
AND encounter_type='M'
ORDER BY begin_time DESC;



