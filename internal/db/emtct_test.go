package db

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"

	"moh.gov.bz/mch/emtct/internal/config"
	"moh.gov.bz/mch/emtct/internal/models"
)

var patientIds = []string{"1111120", "1111121"}

const (
	layoutISO = "2006-01-02"
)

var cnf = config.DbConf{
	Username: "postgres",
	Password: "password",
	Database: "emtct",
	Host:     "localhost",
}

func ClearTable(table string, db EmtctDb) error {
	stmt := fmt.Sprintf("DELETE FROM %s", table)
	_, err := db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("error deleting table(%s) content: %+v", table, err)
	}
	return nil
}

func SamplePatients(db EmtctDb) error {

	for i, id := range patientIds {
		patient := models.Patient{
			Id:               id,
			FirstName:        fmt.Sprintf("First Name - %d", i),
			MiddleName:       "",
			LastName:         fmt.Sprintf("Last Name - %d", i),
			Dob:              time.Date(1992, time.October, 1, 0, 0, 0, 0, time.UTC),
			Ssn:              uuid.New().String(),
			CountryOfBirth:   "Belize",
			DistrictAddress:  "Cayo",
			CommunityAddress: "Las Flores",
			Education:        "High School",
			Ethnicity:        "Mestizo",
			Hiv:              false,
			NextOfKin:        "Jon Doe",
			NextOfKinPhone:   "6539333",
		}
		err := db.CreatePatient(patient)
		if err != nil {
			return fmt.Errorf("failure creating sample patient: %+v", err)
		}
	}
	return nil
}

func SampleDiagnoses(db EmtctDb) error {
	patientId := patientIds[0]
	ds := []string{"2009-02-03", "2010-11-15", "2011-04-12", "2020-10-23"}
	var dates []time.Time
	for _, d := range ds {
		date, _ := time.Parse(layoutISO, d)
		dates = append(dates, date)
	}

	diagnoses := []models.Diagnosis{
		{
			Id:        "1",
			PatientId: patientId,
			Date:      dates[0],
			Name:      "common cold",
		},
		{
			Id:        "2",
			PatientId: patientId,
			Date:      dates[1],
			Name:      "seasonal flu",
		},
		{
			Id:        "3",
			PatientId: patientId,
			Date:      dates[2],
			Name:      "rash",
		},
		{
			Id:        "4",
			PatientId: patientId,
			Date:      dates[3],
			Name:      "common cold",
		},
	}
	for _, d := range diagnoses {
		err := db.CreateDiagnosis(d)
		if err != nil {
			return fmt.Errorf("error creating sample diagnosis: %+v", d)
		}
	}
	return nil
}

func SampleObstetricHistory(db EmtctDb) error {
	history := []models.ObstetricHistory{
		{
			Id:             "1",
			PatientId:      patientIds[0],
			Date:           time.Date(2012, time.May, 1, 0, 0, 0, 0, time.UTC),
			ObstetricEvent: "Miscarriage",
		},
		{
			Id:             "2",
			PatientId:      patientIds[0],
			Date:           time.Date(2015, time.August, 27, 0, 0, 0, 0, time.UTC),
			ObstetricEvent: "Live Born",
		},
	}
	for _, h := range history {
		err := db.CreateObstetricHistory(h)
		if err != nil {
			return fmt.Errorf("error creating sample obstetric history: %+v", err)
		}
	}
	return nil
}

func SamplePregnancies(db EmtctDb) error {
	patientId, _ := strconv.Atoi(patientIds[0])
	pregnancy := models.PregnancyVitals{
		Id:                   1,
		PatientId:            patientId,
		GestationalAge:       4,
		Para:                 10,
		Cs:                   false,
		AbortiveOutcome:      "",
		DiagnosisDate:        time.Date(2020, time.August, 3, 0, 0, 0, 0, time.UTC),
		Planned:              false,
		AgeAtLmp:             28,
		Lmp:                  time.Date(2020, time.July, 6, 0, 0, 0, 0, time.UTC),
		Edd:                  time.Date(2021, time.April, 4, 0, 0, 0, 0, time.UTC),
		DateOfBooking:        time.Date(2020, time.August, 3, 0, 0, 0, 0, time.UTC),
		PrenatalCareProvider: "Public",
		TotalChecks:          2,
	}
	err := db.CreatePregnancy(pregnancy)
	if err != nil {
		return fmt.Errorf("error inserting pregnancy: %+v", err)
	}
	return nil
}

func SampleLabResults(db EmtctDb) error {
	patientId, _ := strconv.Atoi(patientIds[0])
	labResults := []models.LabResult{
		{
			Id:              1,
			PatientId:       patientId,
			TestResult:      "Negative",
			TestName:        "Hb",
			DateSampleTaken: time.Date(2020, time.September, 10, 0, 0, 0, 0, time.UTC),
			ResultDate:      time.Date(2020, time.September, 14, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:              2,
			PatientId:       patientId,
			TestResult:      "Negative",
			TestName:        "Urinalysis",
			DateSampleTaken: time.Date(2020, time.September, 10, 0, 0, 0, 0, time.UTC),
			ResultDate:      time.Date(2020, time.September, 14, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:              3,
			PatientId:       patientId,
			TestResult:      "Negative",
			TestName:        "Hepatitis B",
			DateSampleTaken: time.Date(2020, time.June, 30, 0, 0, 0, 0, time.UTC),
			ResultDate:      time.Date(2020, time.July, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:              4,
			PatientId:       patientId,
			TestResult:      "Negative",
			TestName:        "HIV",
			DateSampleTaken: time.Date(2020, time.June, 30, 0, 0, 0, 0, time.UTC),
			ResultDate:      time.Date(2020, time.July, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:              5,
			PatientId:       patientId,
			TestResult:      "120",
			TestName:        "CD4 Count",
			DateSampleTaken: time.Date(2020, time.June, 30, 0, 0, 0, 0, time.UTC),
			ResultDate:      time.Date(2020, time.July, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:              6,
			PatientId:       patientId,
			TestResult:      "0",
			TestName:        "Viral Load",
			DateSampleTaken: time.Date(2020, time.June, 30, 0, 0, 0, 0, time.UTC),
			ResultDate:      time.Date(2020, time.July, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:              7,
			PatientId:       patientId,
			TestResult:      "Negative",
			TestName:        "Syphilis",
			DateSampleTaken: time.Date(2020, time.June, 30, 0, 0, 0, 0, time.UTC),
			ResultDate:      time.Date(2020, time.July, 3, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, l := range labResults {
		err := db.CreateLabResult(l)
		if err != nil {
			return fmt.Errorf("error inserting lab result: %+v", err)
		}
	}
	return nil
}

func SampleHomeVisits(db EmtctDb) error {
	patientId, _ := strconv.Atoi(patientIds[0])
	homeVisits := []models.HomeVisit{
		{
			Id:          uuid.New().String(),
			PatientId:   patientId,
			Reason:      "Random",
			Comments:    "Patient's vitals are normal",
			CreatedAt:   time.Date(2020, time.September, 29, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   nil,
			CreatedBy:   "nurse@health.gov.bz",
			UpdatedBy:   nil,
			DateOfVisit: time.Date(2020, time.September, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:          uuid.New().String(),
			PatientId:   patientId,
			Reason:      "Periodic",
			Comments:    "All vitals were normal. Patient was given information on breast feeding.",
			CreatedAt:   time.Date(2020, time.October, 16, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   nil,
			CreatedBy:   "nurse@health.gov.bz",
			UpdatedBy:   nil,
			DateOfVisit: time.Date(2020, time.October, 16, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, h := range homeVisits {
		err := db.CreateHomeVisit(h)
		if err != nil {
			return fmt.Errorf("error inserting sample home visit: %+v", err)
		}
	}
	return nil
}

func TestMain(m *testing.M) {
	db, err := NewConnection(&cnf)
	if err != nil {
		log.Errorf("could not connect to database; %+v", err)
		os.Exit(1)
	}
	err = SamplePatients(*db)
	if err != nil {
		log.Errorf("error creating sample patients: %+v", err)
		os.Exit(1)
	}
	err = SampleObstetricHistory(*db)
	if err != nil {
		log.Errorf("error creating sample obstetric histories: %+v", err)
		os.Exit(1)
	}
	err = SampleDiagnoses(*db)
	if err != nil {
		log.Errorf("error creating sample diagnoses %+v", err)
		os.Exit(1)
	}
	err = SamplePregnancies(*db)
	if err != nil {
		log.Errorf("error creating sample pregnancies: %+v", err)
		os.Exit(1)
	}
	err = SampleLabResults(*db)
	if err != nil {
		log.Errorf("error creating sample lab results: %+v", err)
		os.Exit(1)
	}
	err = SampleHomeVisits(*db)
	if err != nil {
		log.Errorf("error creating sample home visits: %+v", err)
		os.Exit(1)
	}

	exitVal := m.Run()
	err = ClearTable("patients", *db)
	if err != nil {
		log.Errorf("failure deleting test table contents: %+v", err)
	}
	os.Exit(exitVal)
}

func TestEmtctDb_FindPatientById(t *testing.T) {
	db, err := NewConnection(&cnf)
	if err != nil {
		t.Fatalf("failed to open db connection %+v", err)
	}
	patientId := patientIds[0]
	patient, err := db.FindPatientById(patientId)
	if err != nil {
		t.Fatalf("expected to find a patient: %+v", err)
	}
	t.Logf("patient: %+v", patient)
}

func TestEmtctDb_FindObstetricHistory(t *testing.T) {
	db, err := NewConnection(&cnf)
	if err != nil {
		t.Fatalf("failed to open db connection %+v", err)
	}
	patientId := patientIds[0]
	history, err := db.FindObstetricHistory(patientId)
	if err != nil {
		t.Fatalf("error retrieving obstetric history: %+v", err)
	}
	t.Logf("history: %+v", history)
}

func TestEmtctDb_FindDiagnoses(t *testing.T) {
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
	db, err := NewConnection(&cnf)
	if err != nil {
		t.Fatalf("failed to open db connection %+v", err)
	}
	patientId := patientIds[0]
	diagnoses, err := db.FindPregnancyDiagnoses(patientId)
	if err != nil {
		t.Fatalf("error fetching pregnancy diagnoses from the database: %+v", err)
	}

	t.Logf("diagnoses: %+v", diagnoses)
}

func TestEmtctDb_FindPregnancyLabResults(t *testing.T) {
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

func TestEmtctDb_CreateHomeVisit(t *testing.T) {
	db, err := NewConnection(&cnf)
	if err != nil {
		t.Fatalf("failed to open db connection %+v", err)
	}

	id := uuid.New().String()
	homeVisit := models.HomeVisit{
		Id:          id,
		PatientId:   1111120,
		Reason:      "Testing",
		Comments:    "Just for Tests",
		DateOfVisit: time.Time{},
		CreatedAt:   time.Time{},
		UpdatedAt:   nil,
		CreatedBy:   "nurse@health.gov.bz",
		UpdatedBy:   nil,
	}
	err = db.CreateHomeVisit(homeVisit)
	if err != nil {
		t.Errorf("failed to create a home visit: %+v", err)
	}

}

func TestEmtctDb_FindHomeVisitById(t *testing.T) {
	db, err := NewConnection(&cnf)
	if err != nil {
		t.Fatalf("failed to open db connection %+v", err)
	}
	homeVisitId := uuid.New().String()
	patientId, _ := strconv.Atoi(patientIds[0])
	h := models.HomeVisit{
		Id:          homeVisitId,
		PatientId:   patientId,
		Reason:      "Test",
		Comments:    "Test",
		DateOfVisit: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
		CreatedBy:   "nurse@health.gov.bz",
		UpdatedBy:   nil,
	}
	_ = db.CreateHomeVisit(h)

	homeVisit, err := db.FindHomeVisitById(h.Id)
	if err != nil {
		t.Fatalf("error querying database: %+v", err)
	}
	if homeVisit == nil {
		t.Error("expected a home visit but got nil")
	}
}

func TestEmtctDb_EditHivScreening(t *testing.T) {
	db, err := NewConnection(&cnf)
	if err != nil {
		t.Fatalf("failed to open db connection %+v", err)
	}
	screeningId := uuid.New().String()
	patientId, _ := strconv.Atoi(patientIds[0])
	s := models.HivScreening{
		Id:                     screeningId,
		PatientId:              patientId,
		TestName:               "PCR 1",
		ScreeningDate:          time.Now(),
		DateSampleReceivedAtHq: nil,
		SampleCode:             "ABD",
		DateSampleShipped:      time.Now(),
		Destination:            "Honduras",
		DateResultReceived:     nil,
		Result:                 "",
		DateResultShared:       nil,
		CreatedAt:              time.Now(),
		UpdatedAt:              nil,
		CreatedBy:              "nurse@health.gov.bz",
		UpdatedBy:              nil,
	}
	_ = db.CreateHivScreening(s)

	updatedBy := "nurse2@health.gov.bz"
	dateResultReceived := time.Now()
	s.Result = "120"
	s.UpdatedBy = &updatedBy
	s.DateResultReceived = &dateResultReceived

	u, err := db.EditHivScreening(s)
	if err != nil {
		t.Fatalf("error while editing hiv screening: %+v", err)
	}
	if u.Result != "120" {
		t.Errorf("got %s, want %s", u.Result, "120")
	}
	if *u.UpdatedBy != updatedBy {
		t.Errorf("got %s, want %s", *u.UpdatedBy, updatedBy)
	}
	if *u.DateResultReceived != dateResultReceived {
		t.Errorf("got %v, want %v", u.DateResultReceived, dateResultReceived)
	}
}
