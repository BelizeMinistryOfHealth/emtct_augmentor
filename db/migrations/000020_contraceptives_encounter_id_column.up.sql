ALTER TABLE contraceptive_used ADD COLUMN mch_encounter_id int NOT NULL;
ALTER TABLE contraceptive_used DROP CONSTRAINT  fk_home_visit_patient;
