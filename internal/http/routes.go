package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

// TestAuth tests that authentication is working
func TestAuth(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	log.Printf("user: %+v", user)
	fmt.Fprintf(w, "TEST")
}

func (a *App) RetrievePatient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patientId := vars["id"]
	patient, err := a.Db.FindPatientById(patientId)
	if err != nil {
		log.WithFields(
			log.Fields{"request": r}).WithError(err).Error("could not find patient with specified id")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(patient)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId, "patient": patient}).WithError(err).Error("error marshalling response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}
