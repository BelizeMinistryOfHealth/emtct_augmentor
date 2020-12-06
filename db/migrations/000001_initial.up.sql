CREATE TABLE patients(
  id BIGINT PRIMARY KEY,
  first_name TEXT NOT NULL,
  middle_name TEXT,
  last_name TEXT NOT NULL,
  dob DATE NOT NULL,
  ssn TEXT,
  country_of_birth TEXT,
  district_address TEXT,
  community_address TEXT,
  education TEXT,
  ethnicity TEXT,
  hiv boolean DEFAULT false,
  next_of_kin TEXT,
  next_of_kin_phone TEXT
);
