ALTER TABLE hospital_admission ADD COLUMN mch_encounter_id int NOT NULL;
ALTER TABLE hospital_admission DROP CONSTRAINT fk_home_visit_patient;
