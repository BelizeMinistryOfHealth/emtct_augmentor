package db

import (
	"testing"

	_ "github.com/lib/pq"

	"moh.gov.bz/mch/emtct/internal/config"
)

func TestEmtctDb_FindPatientById(t *testing.T) {
	cnf := config.DbConf{
		DbUsername: "postgres",
		DbPassword: "password",
		DbDatabase: "emtct",
		DbHost:     "localhost",
	}
	db, err := NewConnection(&cnf)
	if err != nil {
		t.Fatalf("failed to open db connection %+v", err)
	}
	patient, err := db.FindPatientById("1111120")
	if err != nil {
		t.Fatalf("expected to find a patient: %+v", err)
	}
	t.Logf("patient: %+v", patient)
}

func TestEmtctDb_FindObstetricHistory(t *testing.T) {
	cnf := config.DbConf{
		DbUsername: "postgres",
		DbPassword: "password",
		DbDatabase: "emtct",
		DbHost:     "localhost",
	}
	db, err := NewConnection(&cnf)
	if err != nil {
		t.Fatalf("failed to open db connection %+v", err)
	}
	history, err := db.FindObstetricHistory("1111120")
	if err != nil {
		t.Fatalf("error retrieving obstetric history: %+v", err)
	}
	t.Logf("history: %+v", history)
}

func TestEmtctDb_FindDiagnoses(t *testing.T) {
	cnf := config.DbConf{
		DbUsername: "postgres",
		DbPassword: "password",
		DbDatabase: "emtct",
		DbHost:     "localhost",
	}
	db, err := NewConnection(&cnf)
	if err != nil {
		t.Fatalf("failed to open db connection %+v", err)
	}
	diagnoses, err := db.FindDiagnoses("1111120")
	if err != nil {
		t.Fatalf("error retreiving diagnoses: %+v", err)
	}

	t.Logf("diagnosis: %+v", diagnoses)
}

func TestEmtctDb_FindPregnancyDiagnoses(t *testing.T) {
	cnf := config.DbConf{
		DbUsername: "postgres",
		DbPassword: "password",
		DbDatabase: "emtct",
		DbHost:     "localhost",
	}
	db, err := NewConnection(&cnf)
	if err != nil {
		t.Fatalf("failed to open db connection %+v", err)
	}
	patientId := "1111120"
	diagnoses, err := db.FindPregnancyDiagnoses(patientId)
	if err != nil {
		t.Fatalf("error fetching pregnancy diagnoses from the database: %+v", err)
	}

	t.Logf("diagnoses: %+v", diagnoses)
}
