package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/models"
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

type PatientResponse struct {
	Patient          models.Patient            `json:"patient"`
	ObstetricHistory []models.ObstetricHistory `json:"obstetricHistory"`
	Diagnoses        []models.Diagnosis        `json:"diagnoses"`
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

	diagnoses, err := a.Db.FindDiagnoses(patientId)
	if err != nil {
		log.WithFields(
			log.Fields{"request": r}).WithError(err).Error("could not retrieve obstetric history for the patient")
	}

	obstetricHistory, err := a.Db.FindObstetricHistory(patientId)
	if err != nil {
		log.WithFields(log.Fields{"request": r}).WithError(err).Error("could not retrieve obstetric history")
	}
	response := PatientResponse{
		Patient:          *patient,
		ObstetricHistory: obstetricHistory,
		Diagnoses:        diagnoses,
	}
	resp, err := json.Marshal(response)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId, "response": response}).WithError(err).Error("error marshalling response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}
