CREATE TABLE contact_tracing(
    id TEXT PRIMARY KEY,
    patient_id INT,
    comments TEXT,
    date DATE NOT NULL,
    test TEXT NOT NULL,
    test_result TEXT,
    created_by TEXT NOT NULL,
    created_at DATE NOT NULL,
    updated_by TEXT,
    updated_at TEXT
);
