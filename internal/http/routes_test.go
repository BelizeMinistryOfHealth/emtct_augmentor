package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/config"
	"moh.gov.bz/mch/emtct/internal/db"
	"moh.gov.bz/mch/emtct/internal/fixtures"
	"moh.gov.bz/mch/emtct/internal/models"
)

func TestMain(m *testing.M) {
	cnf := config.DbConf{
		Username: "postgres",
		Password: "password",
		Database: "emtct",
		Host:     "localhost",
	}
	d, err := db.NewConnection(&cnf)
	if err != nil {
		log.Errorf("could not connect to database; %+v", err)
		os.Exit(1)
	}
	err = fixtures.SamplePatients(*d)
	if err != nil {
		log.Errorf("could not create sample patients: %+v", err)
		os.Exit(1)
	}
	err = fixtures.SampleObstetricHistory(*d)
	if err != nil {
		log.Errorf("could not create sample obstetric history: %+v", err)
		os.Exit(1)
	}
	err = fixtures.SampleDiagnoses(*d)
	if err != nil {
		log.Errorf("could not create sample diagnoses: %+v", err)
		os.Exit(1)
	}
	err = fixtures.SamplePregnancies(*d)
	if err != nil {
		log.Errorf("could not create sample pregnancies: %+v", err)
		os.Exit(1)
	}
	err = fixtures.SampleLabResults(*d)
	if err != nil {
		log.Errorf("could not create sample lab results: %+v", err)
		os.Exit(1)
	}
	err = fixtures.SampleHomeVisits(*d)
	if err != nil {
		log.Errorf("could not create sample home visit: %+v", err)
		os.Exit(1)
	}

	exitVal := m.Run()
	err = fixtures.ClearTable("patients", *d)
	if err != nil {
		log.Errorf("could not clear test data: %+v", err)
		os.Exit(1)
	}
	os.Exit(exitVal)

}

func TestApp_RetrievePatient(t *testing.T) {
	patientId := "1111120"

	cnf := config.DbConf{
		Username: "postgres",
		Password: "password",
		Database: "emtct",
		Host:     "localhost",
	}
	db, err := db.NewConnection(&cnf)
	if err != nil {
		t.Fatalf("error creating database connection: %+v", err)
	}

	app := App{Db: db}

	r := mux.NewRouter()
	r.HandleFunc("/patient/{id}", app.RetrievePatient)

	req, err := http.NewRequest("GET", fmt.Sprintf("/patient/%s", patientId), nil)
	if err != nil {
		t.Fatalf("error creating request: %+v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code error, want 200, got %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var patient PatientResponse
	_ = json.Unmarshal(body, &patient)
	t.Logf("resp: %+v", patient)
	if patient.Patient.Id != patientId {
		t.Errorf("want: %s, got: %s", patientId, patient.Patient.Id)
	}
}

func TestApp_FindCurrentPregnancy(t *testing.T) {
	patientId := "1111120"

	cnf := config.DbConf{
		Username: "postgres",
		Password: "password",
		Database: "emtct",
		Host:     "localhost",
	}
	db, err := db.NewConnection(&cnf)
	if err != nil {
		t.Fatalf("error creating database connection: %+v", err)
	}

	app := App{Db: db}

	r := mux.NewRouter()
	r.HandleFunc("/patient/{id}", app.FindCurrentPregnancy)

	req, err := http.NewRequest("GET", fmt.Sprintf("/patient/%s", patientId), nil)
	if err != nil {
		t.Fatalf("error creating request: %+v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code error, want 200, got %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var pregnancy PregnancyResponse
	_ = json.Unmarshal(body, &pregnancy)
	t.Logf("resp: %+v", pregnancy)
	if pregnancy.Vitals.PatientId != 1111120 {
		t.Errorf("want: %d, got: %d", 1111120, pregnancy.Vitals.PatientId)
	}
}

func TestApp_FindCurrentPregnancy_WhenPatientDoesNotExist(t *testing.T) {
	patientId := "1111121"

	cnf := config.DbConf{
		Username: "postgres",
		Password: "password",
		Database: "emtct",
		Host:     "localhost",
	}
	db, err := db.NewConnection(&cnf)
	if err != nil {
		t.Fatalf("error creating database connection: %+v", err)
	}

	app := App{Db: db}

	r := mux.NewRouter()
	r.HandleFunc("/patient/{id}/currentPregnancy", app.FindCurrentPregnancy)

	req, err := http.NewRequest("GET", fmt.Sprintf("/patient/%s/currentPregnancy", patientId), nil)
	if err != nil {
		t.Fatalf("error creating request: %+v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code error, want 200, got %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var pregnancy *models.PregnancyVitals
	_ = json.Unmarshal(body, &pregnancy)
	t.Logf("resp: %+v", pregnancy)
	if pregnancy != nil {
		t.Errorf("want: nil, got: %v", pregnancy)
	}
}

func TestApp_FindPregnancyLabResults(t *testing.T) {
	patientId := "1111120"

	cnf := config.DbConf{
		Username: "postgres",
		Password: "password",
		Database: "emtct",
		Host:     "localhost",
	}
	db, err := db.NewConnection(&cnf)
	if err != nil {
		t.Fatalf("error creating database connection: %+v", err)
	}

	app := App{Db: db}

	r := mux.NewRouter()
	r.HandleFunc("/patient/{id}/currentPregnancy/labResults", app.FindPregnancyLabResults)

	req, err := http.NewRequest("GET", fmt.Sprintf("/patient/%s/currentPregnancy/labResults", patientId), nil)
	if err != nil {
		t.Fatalf("error creating request: %+v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code error, want 200, got %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var labResults []models.LabResult
	err = json.Unmarshal(body, &labResults)
	if err != nil {
		t.Fatalf("failed to unmarshal lab results: %+v", err)
	}

	if len(labResults) == 0 {
		t.Error("expected more than one lab result but got 0")
	}
}

func TestApp_FindHomeVisitsByPatient(t *testing.T) {
	patientId := "1111120"

	cnf := config.DbConf{
		Username: "postgres",
		Password: "password",
		Database: "emtct",
		Host:     "localhost",
	}
	db, err := db.NewConnection(&cnf)
	if err != nil {
		t.Fatalf("error creating database connection: %+v", err)
	}

	app := App{Db: db}

	r := mux.NewRouter()
	r.HandleFunc("/patient/{id}/homeVisits", app.FindHomeVisitsByPatient)

	req, err := http.NewRequest("GET", fmt.Sprintf("/patient/%s/homeVisits", patientId), nil)
	if err != nil {
		t.Fatalf("error creating request: %+v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code error, want 200, got %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var s []interface{}
	err = json.Unmarshal(body, &s)
	var homeVisits []models.HomeVisit
	err = json.Unmarshal(body, &homeVisits)
	if err != nil {
		t.Fatalf("failed to unmarshal home visits: %+v", err)
	}

	if len(homeVisits) == 0 {
		t.Error("expected more than one home visit, but got 0")
	}
}

func TestApp_FindHomeVisitById(t *testing.T) {
	cnf := config.DbConf{
		Username: "postgres",
		Password: "password",
		Database: "emtct",
		Host:     "localhost",
	}
	db, err := db.NewConnection(&cnf)
	if err != nil {
		t.Fatalf("error creating database connection: %+v", err)
	}

	app := App{Db: db}

	r := mux.NewRouter()
	r.HandleFunc("/patient/homeVisit/{homeVisitId}", app.HomeVisitApi)

	patientId, _ := strconv.Atoi(fixtures.PatientIds[0])
	h := models.HomeVisit{
		Id:          uuid.New().String(),
		PatientId:   patientId,
		Reason:      "Test",
		Comments:    "Test",
		DateOfVisit: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
		CreatedBy:   "nurse@health.gov.bz",
		UpdatedBy:   nil,
	}
	homeVisitId := h.Id

	app.Db.CreateHomeVisit(h)

	req, err := http.NewRequest("GET", fmt.Sprintf("/patient/homeVisit/%s", homeVisitId), nil)
	if err != nil {
		t.Fatalf("error creating request: %+v", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code error, want 200, got %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var homeVisit models.HomeVisit
	err = json.Unmarshal(body, &homeVisit)
	if err != nil {
		t.Fatalf("failed to unmarshal home visit: %v", err)
	}
	t.Logf("homVisit: %+v", homeVisit)
	if len(homeVisit.Id) == 0 {
		t.Errorf("expected a homeVisit")
	}
}

func TestApp_EditHomeVisit(t *testing.T) {
	cnf := config.DbConf{
		Username: "postgres",
		Password: "password",
		Database: "emtct",
		Host:     "localhost",
	}
	db, err := db.NewConnection(&cnf)
	if err != nil {
		t.Fatalf("error creating database connection: %+v", err)
	}

	app := App{Db: db}

	// Configure the router
	r := mux.NewRouter()
	r.HandleFunc("/patient/homeVisit/{homeVisitId}", app.HomeVisitApi)
	patientId, _ := strconv.Atoi(fixtures.PatientIds[0])

	// Create Home Visit For Testing
	h := models.HomeVisit{
		Id:          uuid.New().String(),
		PatientId:   patientId,
		Reason:      "Test",
		Comments:    "Test",
		DateOfVisit: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
		CreatedBy:   "nurse@health.gov.bz",
		UpdatedBy:   nil,
	}
	homeVisitId := h.Id
	app.Db.CreateHomeVisit(h)
	homeVisitRequest := HomeVisitRequest{
		Reason:      "Edited Reason",
		Comments:    "Edited Comments",
		DateOfVisit: h.DateOfVisit,
	}
	body, err := json.Marshal(homeVisitRequest)
	if err != nil {
		t.Fatalf("error marshalling request: %+v", err)
	}

	// Create the PUT Request
	req, err := http.NewRequest("PUT", fmt.Sprintf("/patient/homeVisit/%s", homeVisitId), bytes.NewReader(body))
	if err != nil {
		t.Fatalf("error creating request: %+v", err)
	}
	w := httptest.NewRecorder()

	// Configure the jwt token for testing
	jwtToken := JwtToken{
		Email: "nurse@health.gov.bz",
	}
	ctx := context.WithValue(req.Context(), "user", jwtToken)
	req = req.WithContext(ctx)
	r.ServeHTTP(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code error, want 200, got %d", resp.StatusCode)
	}

	// Unmarshal the response into the HomeVisit model
	respBytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var homeVisit models.HomeVisit
	err = json.Unmarshal(respBytes, &homeVisit)
	if err != nil {
		t.Fatalf("failed to unmarshal home visit: %v", err)
	}

	if homeVisit.Reason != "Edited Reason" {
		t.Errorf("wanted: %s, got: %s", "Expected Reason", homeVisit.Reason)
	}
}

func TestApp_CreateHivScreening(t *testing.T) {
	cnf := config.DbConf{
		Username: "postgres",
		Password: "password",
		Database: "emtct",
		Host:     "localhost",
	}
	db, err := db.NewConnection(&cnf)
	if err != nil {
		t.Fatalf("error creating database connection: %+v", err)
	}

	app := App{Db: db}

	// Configure the router
	r := mux.NewRouter()
	r.HandleFunc("/patient/hivHomeScreening", app.CreateHivScreeningHandler)

	// Create Request
	dateShipped := time.Date(2020, time.September, 30, 0, 0, 0, 0, time.UTC)
	request := NewHivScreeningRequest{
		PatientId:              1111120,
		TestName:               "PCR 1",
		ScreeningDate:          time.Date(2020, time.September, 21, 10, 20, 00, 00, time.UTC),
		DateSampleReceivedAtHq: nil,
		SampleCode:             "CODE_PCR_1234",
		DateSampleShipped:      &dateShipped,
		Destination:            "Honduras",
		DateResultReceived:     nil,
		Result:                 "",
		DateResultShared:       nil,
	}
	body, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("error marshalling request: %+v", err)
	}
	// Create POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("/patient/hivHomeScreening"), bytes.NewReader(body))
	if err != nil {
		t.Fatalf("error creating new http request: %+v", err)
	}
	w := httptest.NewRecorder()

	// Configure the jwt token for testing
	jwtToken := JwtToken{
		Email: "nurse@health.gov.bz",
	}
	ctx := context.WithValue(req.Context(), "user", jwtToken)
	req = req.WithContext(ctx)
	r.ServeHTTP(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code error, want 200, got %d", resp.StatusCode)
	}

	// Unmarshal the response into the HomeVisit model
	respBytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var screening models.HivScreening
	err = json.Unmarshal(respBytes, &screening)
	if err != nil {
		t.Errorf("error unmarshalling result of new screening: %+v", err)
	}
	t.Logf("hiv screening: %+v", screening)
}
