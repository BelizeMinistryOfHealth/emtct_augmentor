ALTER TABLE contact_tracing ALTER COLUMN updated_at TYPE TEXT;
ALTER TABLE contact_tracing ALTER COLUMN patient_id DROP NOT NULL;
