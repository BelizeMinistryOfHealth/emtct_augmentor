CREATE TABLE lab_results(
    id BIGINT PRIMARY KEY,
    patient_id BIGINT,
    test_result TEXT NOT NULL,
    test_name TEXT NOT NULL,
    date_sample_taken DATE,
    result_date DATE,
    CONSTRAINT fk_lab_results_patient
        FOREIGN KEY(patient_id)
        REFERENCES patients(id)
        ON DELETE CASCADE
);
