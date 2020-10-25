CREATE TABLE hiv_screening(
    id TEXT PRIMARY KEY,
    patient_id bigint,
    test_name TEXT NOT NULL,
    screening_date TIMESTAMP NOT NULL,
    date_sample_received_at_hq DATE,
    sample_code TEXT NOT NULL,
    date_sample_shipped DATE NOT NULL,
    destination TEXT NOT NULL,
    date_result_received DATE,
    result TEXT,
    date_result_shared DATE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by TEXT NOT NULL,
    updated_by TEXT,
    CONSTRAINT fk_home_visit_patient
        FOREIGN KEY(patient_id)
            REFERENCES patients(id)
            ON DELETE CASCADE
);
