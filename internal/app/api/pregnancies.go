package api

import (
	"encoding/json"
	"fmt"
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

type pregnancyResponse struct {
	Vitals    *pregnancy.Vitals     `json:"vitals"`
	Diagnoses []pregnancy.Diagnosis `json:"diagnoses"`
}

func (a *pregnancyRoutes) FindCurrentPregnancy(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		patientId := vars["patientId"]
		id, err := strconv.Atoi(patientId)
		if err != nil {
			log.WithFields(log.Fields{"patientId": patientId}).WithError(err).Error("patient id is not a number")
			http.Error(w, "the patient id provided is invalid", http.StatusBadRequest)
			return
		}
		preg, err := a.Pregnancies.FindCurrentPregnancy(id)
		if err != nil {
			log.WithFields(log.Fields{"patientId": patientId}).
				WithError(err).
				Error("error retrieving current pregnancy from database")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		diagnoses, err := a.Pregnancies.FindDiagnosesDuringPregnancy(id)
		if err != nil {
			log.WithFields(log.Fields{"patientId": patientId, "pregnancy": preg}).
				WithError(err).
				Error("error fetching diagnoses for a pregnancy")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if preg == nil {
			w.Header().Add("Content-Type", "application/json")
			r, _ := json.Marshal(nil)
			fmt.Fprintf(w, string(r))
		}

		response := pregnancyResponse{
			Vitals:    preg,
			Diagnoses: diagnoses,
		}

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{"patientId": patientId, "pregnancy": preg, "diagnoses": diagnoses}).
				WithError(err).
				Error("error marshalling pregnancy")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

}

type pregnancyLabResultsResponse struct {
	LabResults []labs.LabResult  `json:"labResults"`
	Patient    patient.BasicInfo `json:"patient"`
}

func (a *pregnancyRoutes) FindPregnancyLabResults(w http.ResponseWriter, r *http.Request) {
	handlerName := "FindPregnancyLabResults"
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		method := "GET"
		vars := mux.Vars(r)
		patientId := vars["id"]
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
		patient, err := a.Patient.FindBasicInfo(id)
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
	Patient          patient.BasicInfo            `json:"patient"`
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
		patientInfo, err := a.Patient.FindBasicInfo(patientId)
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
