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
	var pregnancy *models.PregnancyVitals
	_ = json.Unmarshal(body, &pregnancy)
	t.Logf("resp: %+v", pregnancy)
	if pregnancy != nil {
		t.Errorf("want: nil, got: %v", pregnancy)
	}
}
