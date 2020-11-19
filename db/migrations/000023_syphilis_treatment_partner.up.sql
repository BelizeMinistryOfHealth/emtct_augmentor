CREATE TABLE syphilis_treatment_partner(
    id Text PRIMARY KEY,
    medication_name TEXT NOT NULL,
    dosage TEXT,
    comments TEXT,
    date DATE NOT NULL,
    created_by TEXT NOT NULL,
    created_at DATE NOT NULL,
    updated_by TEXT,
    updated_at DATE
);
