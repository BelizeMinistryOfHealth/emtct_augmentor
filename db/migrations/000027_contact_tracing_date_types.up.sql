ALTER TABLE contact_tracing ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at::timestamp;
ALTER TABLE contact_tracing ALTER COLUMN created_at TYPE TIMESTAMP USING created_at::timestamp;
