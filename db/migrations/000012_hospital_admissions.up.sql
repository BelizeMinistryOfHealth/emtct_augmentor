CREATE TABLE hospital_admission(
    id TEXT PRIMARY KEY,
    patient_id int NOT NULL,
    date_admitted DATE NOT NULL,
    facility TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    created_by TEXT NOT NULL,
    updated_at TIMESTAMP,
    updated_by TEXT,
    CONSTRAINT fk_home_visit_patient
        FOREIGN KEY(patient_id)
            REFERENCES patients(id)
            ON DELETE CASCADE
);
