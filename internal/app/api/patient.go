package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/app"
	"moh.gov.bz/mch/emtct/internal/business/data/patient"
	"moh.gov.bz/mch/emtct/internal/business/data/pregnancy"
	"moh.gov.bz/mch/emtct/internal/business/data/prescription"
)

type patientResponse struct {
	Patient          *patient.Patient              `json:"patient"`
	ObstetricHistory []pregnancy.ObstetricHistory  `json:"obstetricHistory"`
	Diagnoses        []pregnancy.Diagnosis         `json:"diagnoses"`
	AncEncounter     *pregnancy.AntenatalEncounter `json:"antenatalEncounter"`
}

func (a *pregnancyRoutes) RetrievePatientHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
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
		patient, err := a.Patient.FindByPatientId(id)
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
		hivDiagnoses, err := a.Hiv.FindHivDiagnoses(id)
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

		diagnoses, err := a.Pregnancies.FindDiagnosesBeforePregnancy(id)
		if err != nil {
			log.WithFields(
				log.Fields{"request": r}).WithError(err).Error("could not retrieve obstetric history for the patient")
		}
		if diagnoses == nil {
			diagnoses = []pregnancy.Diagnosis{}
		}

		obstetricHistory, err := a.Pregnancies.FindObstetricHistory(id)
		if err != nil {
			log.WithFields(log.Fields{"request": r}).WithError(err).Error("could not retrieve obstetric history")
		}

		v, err := a.Pregnancies.FindObstetricDetails(id)
		var lmp *time.Time
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
			}).WithError(err).Error("could not find patient's obstetric details, using nil for lmp")
		}
		if v != nil {
			lmp = v.Lmp
		}
		ancEncounter, err := a.Pregnancies.FindLatestAntenatalEncounter(id, lmp)
		if err != nil {
			log.WithFields(log.Fields{"patientId": id}).WithError(err).Error("could not retrieve anc encounter")
		}

		w.Header().Add("Content-Type", "application/json")

		if patient == nil {
			emptyResponse := patientResponse{
				Patient:          nil,
				ObstetricHistory: nil,
				Diagnoses:        nil,
			}
			_ = json.NewEncoder(w).Encode(emptyResponse)
			return
		}
		response := patientResponse{
			Patient:          patient,
			ObstetricHistory: obstetricHistory,
			Diagnoses:        diagnoses,
			AncEncounter:     ancEncounter,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{"patientId": patientId, "response": response}).WithError(err).Error("error marshalling response")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

type arvsResponse struct {
	Arvs    []prescription.Prescription `json:"arvs"`
	Patient patient.BasicInfo           `json:"patient"`
}

func (a *pregnancyRoutes) ArvsHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "ArvsHandler"
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		method := "Get"
		token := r.Context().Value("user").(app.JwtToken)
		user := token.Email
		vars := mux.Vars(r)
		id := vars["patientId"]
		patientId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).
				WithError(err).
				Error("patient id is not a valid number")
			http.Error(w, "patientId is not a valid number", http.StatusBadRequest)
			return
		}
		// Find the pregnancy and the lmp so we can get the date bounds
		pregs, err := a.Pregnancies.FindCurrentPregnancy(patientId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).
				WithError(err).
				Error("error retrieving patient's current pregnancy")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if pregs == nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).
				WithError(err).
				Error("patient does not have a current pregnancy")
			http.Error(w, "patient does not have a current pregnancy", http.StatusNotFound)
			return
		}
		lmp := pregs.Lmp
		nextDate := lmp.Add(time.Hour * 24 * 7 * 54)
		arvs, err := a.Patient.FindArvsByPatient(patientId, *lmp, nextDate)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).
				WithError(err).
				Error("patient does not have a current pregnancy")
			http.Error(w, "patient does not have a current pregnancy", http.StatusNotFound)
			return
		}
		patientInfo, err := a.Patient.FindBasicInfo(patientId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).
				WithError(err).
				Error("error retrieving patient's basic info")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		arvsResponse := arvsResponse{
			Arvs:    arvs,
			Patient: *patientInfo,
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(arvsResponse); err != nil {
			log.WithFields(log.Fields{"patientId": patientId, "user": user}).
				WithError(err).
				Error("marshalling arvs response failed")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}
}
