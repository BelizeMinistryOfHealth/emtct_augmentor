ALTER TABLE syphilis_treatment_partner ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at::timestamp;
ALTER TABLE syphilis_treatment_partner ALTER COLUMN created_at TYPE TIMESTAMP USING created_at::timestamp;
ALTER TABLE syphilis_treatment_partner ALTER COLUMN date TYPE TIMESTAMP USING date::timestamp;
