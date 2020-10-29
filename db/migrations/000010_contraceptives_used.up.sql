CREATE TABLE contraceptive_used(
    id TEXT PRIMARY KEY,
    patient_id bigint,
    comments TEXT,
    date_used TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    created_by TEXT NOT NULL,
    updated_at TIMESTAMP,
    updated_by TEXT,
    CONSTRAINT fk_home_visit_patient
        FOREIGN KEY(patient_id)
            REFERENCES patients(id)
            ON DELETE CASCADE
)
