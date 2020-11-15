package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/models"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

// TestAuth tests that authentication is working
func TestAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

type PatientResponse struct {
	Patient          *models.Patient            `json:"patient"`
	ObstetricHistory []models.ObstetricHistory  `json:"obstetricHistory"`
	Diagnoses        []models.Diagnosis         `json:"diagnoses"`
	AncEncounter     *models.AntenatalEncounter `json:"antenatalEncounter"`
}

func (a *App) RetrievePatient(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	patientId := vars["id"]
	id, err := strconv.Atoi(patientId)
	if err != nil {
		log.WithFields(log.Fields{
			"patientId": patientId,
		}).WithError(err).Error("patient id is not a valid numeric value")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	patient, err := a.AcsisDb.FindByPatientId(id)
	if err != nil {
		log.WithFields(
			log.Fields{"request": r}).WithError(err).Error("could not find patient with specified id")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// We assume that the patient's HIV status is negative.
	patient.Hiv = false

	// Retrieve all the hiv diagnoses so we can get the
	// patient's HIV status and first date of diagnoses.
	hivDiagnoses, err := a.AcsisDb.FindHivDiagnoses(id)
	if err != nil {
		log.WithFields(log.Fields{
			"patientId": id,
			"patient":   patient,
			"handler":   "RetrievePatient",
		}).WithError(err).Error("error retrieving patient hiv diagnoses: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// If there are any diagnoses, we retrieve the first diagnosis and use its date as the diagnosis date.
	if len(hivDiagnoses) > 0 {
		patient.Hiv = true
		patient.HivDiagnosisDate = &hivDiagnoses[0].Date
	}

	diagnoses, err := a.AcsisDb.FindDiagnosesBeforePregnancy(id)
	if err != nil {
		log.WithFields(
			log.Fields{"request": r}).WithError(err).Error("could not retrieve obstetric history for the patient")
	}
	if diagnoses == nil {
		diagnoses = []models.Diagnosis{}
	}

	obstetricHistory, err := a.AcsisDb.FindObstetricHistory(id)
	if err != nil {
		log.WithFields(log.Fields{"request": r}).WithError(err).Error("could not retrieve obstetric history")
	}

	v, err := a.AcsisDb.FindObstetricDetails(id)
	var lmp *time.Time
	if err != nil {
		log.WithFields(log.Fields{
			"patientId": id,
		}).WithError(err).Error("could not find patient's obstetric details, using nil for lmp")
	}
	if v != nil {
		lmp = v.Lmp
	}
	ancEncounter, err := a.AcsisDb.FindLatestAntenatalEncounter(id, lmp)
	if err != nil {
		log.WithFields(log.Fields{"patientId": id}).WithError(err).Error("could not retrieve anc encounter")
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
		AncEncounter:     ancEncounter,
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
	id, err := strconv.Atoi(patientId)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId}).WithError(err).Error("patient id is not a number")
		http.Error(w, "the patient id provided is invalid", http.StatusBadRequest)
		return
	}
	pregnancy, err := a.AcsisDb.FindCurrentPregnancy(id)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId}).
			WithError(err).
			Error("error retrieving current pregnancy from database")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	diagnoses, err := a.AcsisDb.FindDiagnosesDuringPregnancy(id)
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

type pregnancyLabResultsResponse struct {
	LabResults []models.LabResult      `json:"labResults"`
	Patient    models.PatientBasicInfo `json:"patient"`
}

func (a *App) FindPregnancyLabResults(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	patientId := vars["id"]
	id, err := strconv.Atoi(patientId)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId}).
			WithError(err).Error("patient id is not a numeric value")
		http.Error(w, "patient id must be a valid number", http.StatusBadRequest)
		return
	}
	labResults, err := a.AcsisDb.FindLabTestsDuringPregnancy(id)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId}).
			WithError(err).
			Error("error while retrieving lab results")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	patient, err := a.AcsisDb.FindPatientBasicInfo(id)
	if err != nil {
		log.WithFields(log.Fields{
			"patientId": id,
			"handler":   "FindPregnancyLabResults",
		}).
			WithError(err).
			Error("error fetching patient information")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := pregnancyLabResultsResponse{Patient: *patient, LabResults: labResults}
	results, err := json.Marshal(response)
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
