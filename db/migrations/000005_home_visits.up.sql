CREATE TABLE home_visit(
    id BIGINT PRIMARY KEY,
    patient_id BIGINT NOT NULL,
    reason TEXT,
    comments TEXT,
    created_at DATE NOT NULL,
    updated_at DATE,
    created_by TEXT,
    updated_by TEXT,
    CONSTRAINT fk_home_visit_patient
        FOREIGN KEY(patient_id)
        REFERENCES patients(id)
        ON DELETE CASCADE
);
