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

func TestEmtctDb_FindPregnancyLabResults(t *testing.T) {
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

	labResults, err := db.FindPregnancyLabResults(patientId)
	if err != nil {
		t.Fatalf("error retrieving lab results from the database: %+v", err)
	}
	if len(labResults) == 0 {
		t.Errorf("expected to have at least one lab result")
	}
}

func TestEmtctDb_FindHomeVisitsByPatientId(t *testing.T) {
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

	homeVisits, err := db.FindHomeVisitsByPatientId(patientId)
	if err != nil {
		t.Fatalf("error querying database for home visits: %+v", err)
	}

	if len(homeVisits) == 0 {
		t.Errorf("expected at least one home visit, but got 0")
	}
	t.Logf("Home visits: %+v", homeVisits)
}

func TestEmtctDb_FindHomeVisitById(t *testing.T) {
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
	homeVisitId := "1"

	homeVisit, err := db.FindHomeVisitById(homeVisitId)
	if err != nil {
		t.Fatalf("error querying database: %+v", err)
	}
	if homeVisit == nil {
		t.Error("expected a home visit but got nil")
	}

}
