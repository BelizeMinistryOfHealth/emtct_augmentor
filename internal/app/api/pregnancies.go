package api

import (
	"encoding/json"
	"moh.gov.bz/mch/emtct/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/business/data/hiv"
	"moh.gov.bz/mch/emtct/internal/business/data/labs"
	"moh.gov.bz/mch/emtct/internal/business/data/patient"
	"moh.gov.bz/mch/emtct/internal/business/data/pregnancy"
)

type pregnancyRoutes struct {
	Patient patient.Patients
	Hiv     hiv.HIV
	Lab     labs.Labs
}

func (a *pregnancyRoutes) GetPregnancy(w http.ResponseWriter, r *http.Request) {
	handlerName := "GetPregnancy"
	w.Header().Add("Content-Type", "application/json")
	type response struct {
		Pregnancy                models.Pregnancy   `json:"pregnancy"`
		Interval                 int                `json:"interval"`
		Patient                  models.Patient     `json:"patient"`
		DiagnosesBeforePregnancy []models.Diagnosis `json:"diagnosesBeforePregnancy"`
		DiagnosesDuringPregnancy []models.Diagnosis `json:"diagnosesDuringPregnancy"`
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
		diagnoses, err := a.Patient.GetDiagnoses(patientId)
		if err != nil {
			log.WithFields(log.Fields{
				"handler":   handlerName,
				"method":    r.Method,
				"patientId": patientId,
			}).WithError(err).Error("did not find any diagnoses")
		}

		endDate := preg.Edd
		if preg.Edd != nil {
			l := *preg.Lmp
			e := l.Add(time.Hour * 24 * 7 * 54)
			endDate = &e
		}
		diagnosesDuringPregnancy := DiagnosesBetweenDates(diagnoses, *preg.Lmp, *endDate)
		diagnosesBeforePregnancy := DiagnosesBeforeDate(diagnoses, *preg.Lmp)

		resp := response{
			Pregnancy:                preg,
			Interval:                 interval,
			Patient:                  *patient,
			DiagnosesBeforePregnancy: diagnosesBeforePregnancy,
			DiagnosesDuringPregnancy: diagnosesDuringPregnancy,
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

func DiagnosesBeforeDate(diagnoses []models.Diagnosis, d time.Time) []models.Diagnosis {
	var ds []models.Diagnosis
	for _, i := range diagnoses {
		if i.Date.Before(d) {
			ds = append(ds, i)
		}
	}
	return ds
}

func DiagnosesBetweenDates(diagnoses []models.Diagnosis, start, end time.Time) []models.Diagnosis {
	var ds []models.Diagnosis
	for _, i := range diagnoses {
		if i.Date.After(start) && i.Date.Before(end) {
			ds = append(ds, i)
		}
	}
	return ds
}

type pregnancyLabResultsResponse struct {
	LabResults []models.LabResult `json:"labResults"`
	Patient    models.Patient     `json:"patient"`
}

func (a *pregnancyRoutes) FindPregnancyLabResults(w http.ResponseWriter, r *http.Request) {
	handlerName := "FindPregnancyLabResults"
	w.Header().Add("Content-Type", "application/json")

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
		pregId := vars["pregnancyId"]
		pregnancyId, err := strconv.Atoi(pregId)
		if err != nil {
			log.WithFields(log.Fields{"patientId": patientId, "pregnancyId": pregId}).WithError(err).
				Error("pregnancy id is not a valid numeric value")
			http.Error(w, "pregnancy id is not a valid numeric value", http.StatusBadRequest)
			return
		}
		preg, err := a.Patient.GetPregnancy(pregnancyId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
				"handler":   handlerName,
				"method":    method,
			}).WithError(err).Error("no pregnancy found")
			http.Error(w, "did not find any pregnancy with given id", http.StatusNotFound)
			return
		}
		endDate := preg.Lmp.Add(time.Hour * 24 * 7 * 54)
		labResults, err := a.Patient.FindLabTestsInPeriod(patientId, *preg.Lmp, endDate)
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
