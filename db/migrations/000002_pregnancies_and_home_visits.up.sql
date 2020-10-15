CREATE TABLE obstetric_history(
  patient_id BIGINT NOT NULL,
  event_date DATE NOT NULL,
  event_name TEXT NOT NULL,
  CONSTRAINT fk_obstetric_history_patient
    FOREIGN KEY(patient_id)
    REFERENCES patients(id)
    ON DELETE CASCADE
);

CREATE TABLE diagnoses(
  id BIGINT PRIMARY KEY,
  patient_id BIGINT NOT NULL,
  diagnosis_date DATE NOT NULL,
  diagnosis_name TEXT NOT NULL,
  CONSTRAINT fk_diagnosis_patient
    FOREIGN KEY(patient_id)
    REFERENCES patients(id)
    ON DELETE CASCADE
);

CREATE TABLE pregnancies(
  id BIGINT PRIMARY KEY,
  patient_id BIGINT NOT NULL,
  gestational_age INT,
  para INT,
  cs BOOLEAN,
  abortive_outcome TEXT,
  diagnosis_date DATE,
  planned BOOLEAN,
  age_at_lmp INT,
  lmp DATE,
  edd DATE,
  date_of_booking DATE,
  prenatal_care_provider TEXT,
  total_checks INT,
  CONSTRAINT fk_pregnancies_patient
    FOREIGN KEY(patient_id)
    REFERENCES patients(id)
    ON DELETE CASCADE
);


CREATE TABLE home_visits(
  id BIGINT PRIMARY KEY,
  patient_id BIGINT,
  reason TEXT NOT NULL,
  visit_date DATE NOT NULL,
  comments TEXT,
  CONSTRAINT fk_home_visits_patient
    FOREIGN KEY(patient_id)
    REFERENCES patients(id)
    ON DELETE CASCADE
)