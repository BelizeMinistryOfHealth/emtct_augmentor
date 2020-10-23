package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"moh.gov.bz/mch/emtct/internal/config"
	"moh.gov.bz/mch/emtct/internal/db"
	"moh.gov.bz/mch/emtct/internal/models"
)

func TestApp_RetrievePatient(t *testing.T) {
	patientId := "1111120"

	cnf := config.DbConf{
		DbUsername: "postgres",
		DbPassword: "password",
		DbDatabase: "emtct",
		DbHost:     "localhost",
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
		DbUsername: "postgres",
		DbPassword: "password",
		DbDatabase: "emtct",
		DbHost:     "localhost",
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
		DbUsername: "postgres",
		DbPassword: "password",
		DbDatabase: "emtct",
		DbHost:     "localhost",
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
		DbUsername: "postgres",
		DbPassword: "password",
		DbDatabase: "emtct",
		DbHost:     "localhost",
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
		DbUsername: "postgres",
		DbPassword: "password",
		DbDatabase: "emtct",
		DbHost:     "localhost",
	}
	db, err := db.NewConnection(&cnf)
	if err != nil {
		t.Fatalf("error creating database connection: %+v", err)
	}

	app := App{Db: db}

	r := mux.NewRouter()
	r.HandleFunc("/patient/{id}/homeVisits", app.FindPregnancyLabResults)

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
		DbUsername: "postgres",
		DbPassword: "password",
		DbDatabase: "emtct",
		DbHost:     "localhost",
	}
	db, err := db.NewConnection(&cnf)
	if err != nil {
		t.Fatalf("error creating database connection: %+v", err)
	}

	app := App{Db: db}

	r := mux.NewRouter()
	r.HandleFunc("/patient/homeVisit/{homeVisitId}", app.FindHomeVisitById)

	homeVisitId := "1"
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

}
