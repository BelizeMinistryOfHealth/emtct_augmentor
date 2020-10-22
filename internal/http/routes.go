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
	Patient          *models.Patient           `json:"patient"`
	ObstetricHistory []models.ObstetricHistory `json:"obstetricHistory"`
	Diagnoses        []models.Diagnosis        `json:"diagnoses"`
}

func (a *App) RetrievePatient(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		return
	}

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

	if patient == nil {
		w.Header().Add("Content-Type", "application/json")
		emptyResponse := PatientResponse{
			Patient:          nil,
			ObstetricHistory: nil,
			Diagnoses:        nil,
		}
		resp, _ := json.Marshal(emptyResponse)
		fmt.Fprint(w, string(resp))
		return
	}
	response := PatientResponse{
		Patient:          patient,
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

type PregnancyResponse struct {
	Vitals    *models.PregnancyVitals `json:"vitals"`
	Diagnoses []models.Diagnosis      `json:"diagnoses"`
}

func (a *App) FindCurrentPregnancy(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}
	vars := mux.Vars(r)
	patientId := vars["id"]
	pregnancy, err := a.Db.FindCurrentPregnancy(patientId)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId}).
			WithError(err).
			Error("error retrieving current pregnancy from database")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	diagnoses, err := a.Db.FindPregnancyDiagnoses(patientId)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId, "pregnancy": pregnancy}).
			WithError(err).
			Error("error fetching diagnoses for a pregnancy")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if pregnancy == nil {
		w.Header().Add("Content-Type", "application/json")
		r, _ := json.Marshal(nil)
		fmt.Fprintf(w, string(r))
	}

	response := PregnancyResponse{
		Vitals:    pregnancy,
		Diagnoses: diagnoses,
	}

	resp, err := json.Marshal(response)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId, "pregnancy": pregnancy, "diagnoses": diagnoses}).
			WithError(err).
			Error("error marshalling pregnancy")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}

func (a *App) FindPregnancyLabResults(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	patientId := vars["id"]
	labResults, err := a.Db.FindPregnancyLabResults(patientId)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId}).
			WithError(err).
			Error("error while retrieving lab results")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	results, err := json.Marshal(labResults)
	if err != nil {
		log.WithFields(log.Fields{"labResults": labResults, "patientId": patientId}).
			WithError(err).
			Error("error marshalling lab results for a pregnancy")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(results))
}
