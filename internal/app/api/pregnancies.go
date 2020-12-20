package api

import (
	"encoding/json"
	"moh.gov.bz/mch/emtct/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/app"
	"moh.gov.bz/mch/emtct/internal/business/data/hiv"
	"moh.gov.bz/mch/emtct/internal/business/data/labs"
	"moh.gov.bz/mch/emtct/internal/business/data/patient"
	"moh.gov.bz/mch/emtct/internal/business/data/pregnancy"
)

type pregnancyRoutes struct {
	Pregnancies pregnancy.Pregnancies
	Patient     patient.Patients
	Hiv         hiv.HIV
	Lab         labs.Labs
}

func (a *pregnancyRoutes) GetPregnancy(w http.ResponseWriter, r *http.Request) {
	handlerName := "GetPregnancy"
	w.Header().Add("Content-Type", "application/json")
	type response struct {
		Pregnancy models.Pregnancy `json:"pregnancy"`
		Interval  int              `json:"interval"`
		Patient   models.Patient   `json:"patient"`
	}

	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		pregId := vars["pregnancyId"]
		pregnancyId, err := strconv.Atoi(pregId)
		if err != nil {
			log.WithFields(log.Fields{
				"handler":     handlerName,
				"method":      r.Method,
				"pregnancyId": pregId,
			}).WithError(err).Error("invalid pregnancyId")
			http.Error(w, "invalid pregnancy id", http.StatusBadRequest)
			return
		}
		patientId := vars["patientId"]

		pregs, err := a.Patient.GetPregnancies(patientId)
		idx := models.PregnancyIndex(pregs, pregnancyId)
		if idx == -1 {
			// not available
			return
		}
		preg := pregs[idx]
		var interval int
		if len(pregs) > 1 {
			// increment the index to get the previous pregnancy
			prevPreg := pregs[idx+1]
			interval = preg.Lmp.Year() - prevPreg.Lmp.Year()
		}

		patient, err := a.Patient.FindByPatientId(patientId)
		if err != nil {
			http.Error(w, "patient does not exist", http.StatusNoContent)
			return
		}
		resp := response{
			Pregnancy: preg,
			Interval:  interval,
			Patient:   *patient,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.WithFields(log.Fields{
				"handler":  handlerName,
				"method":   r.Method,
				"response": resp,
			}).WithError(err).Error("failed to encode pregnancy")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

}

type pregnancyLabResultsResponse struct {
	LabResults []labs.LabResult `json:"labResults"`
	Patient    models.Patient   `json:"patient"`
}

func (a *pregnancyRoutes) FindPregnancyLabResults(w http.ResponseWriter, r *http.Request) {
	handlerName := "FindPregnancyLabResults"
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		method := "GET"
		vars := mux.Vars(r)
		patientId := vars["patientId"]
		id, err := strconv.Atoi(patientId)
		if err != nil {
			log.WithFields(log.Fields{"patientId": patientId}).
				WithError(err).Error("patient id is not a numeric value")
			http.Error(w, "patient id must be a valid number", http.StatusBadRequest)
			return
		}
		preg, err := a.Pregnancies.FindCurrentPregnancy(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
				"handler":   handlerName,
				"method":    method,
			}).WithError(err).Error("no current pregnancy found")
			http.Error(w, "did not find any current pregnancy", http.StatusNotFound)
			return
		}
		labResults, err := a.Lab.FindLabTestsDuringPregnancy(id, preg.Lmp)
		if err != nil {
			log.WithFields(log.Fields{"patientId": patientId}).
				WithError(err).
				Error("error while retrieving lab results")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		patient, err := a.Patient.FindByPatientId(patientId)
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
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{"labResults": labResults, "patientId": patientId}).
				WithError(err).
				Error("error marshalling lab results for a pregnancy")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

type obstetricHistoryResponse struct {
	ObstetricHistory []pregnancy.ObstetricHistory `json:"obstetricHistory"`
	Patient          models.Patient               `json:"patient"`
}

func (a *pregnancyRoutes) ObstetricHistoryHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		id := vars["patientId"]
		patientId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{"patientId": id}).Error("patientId is not a number")
			http.Error(w, "patientId must be a number", http.StatusBadRequest)
			return
		}
		token := r.Context().Value("user").(app.JwtToken)
		user := token.Email
		obstetricHistory, err := a.Pregnancies.FindObstetricHistory(patientId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
				"user":      user,
				"handler":   "ObstetricHistoryHandler",
			}).WithError(err).Error("error retrieving obstetric history")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		patientInfo, err := a.Patient.FindByPatientId(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
				"user":      user,
				"handler":   "ObstetricHistoryHandler",
			}).
				WithError(err).
				Error("error retrieving patient basic info")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// Return an empty array if no results are found
		if obstetricHistory == nil {
			obstetricHistory = []pregnancy.ObstetricHistory{}
		}
		response := obstetricHistoryResponse{
			ObstetricHistory: obstetricHistory,
			Patient:          *patientInfo,
		}
		w.Header().Add("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
				"response":  response,
				"user":      user,
				"handler":   "ObstetricHistoryHandler",
			}).
				WithError(err).
				Error("error marshalling obstetric history response")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
