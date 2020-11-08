ALTER TABLE hiv_screening ADD COLUMN mch_encounter_id int NOT NULL;
ALTER TABLE hiv_screening DROP CONSTRAINT fk_home_visit_patient;
ALTER TABLE hiv_screening ADD COLUMN date_sample_taken DATE;
ALTER TABLE hiv_screening DROP COLUMN date_sample_shipped;
ALTER TABLE hiv_screening ADD COLUMN date_sample_shipped DATE;
