ALTER TABLE contact_tracing ALTER COLUMN updated_at TYPE DATE USING updated_at::date;
ALTER TABLE contact_tracing ALTER COLUMN created_at TYPE DATE USING created_at::date;
