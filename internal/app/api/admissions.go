package api

import (
	"encoding/json"
	"moh.gov.bz/mch/emtct/internal/auth"
	"moh.gov.bz/mch/emtct/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/business/data/admissions"
	"moh.gov.bz/mch/emtct/internal/business/data/patient"
)

type AdmissionRoutes struct {
	Admissions admissions.Admissions
	Patients   patient.Patients
}

type admissionsResponse struct {
	HospitalAdmissions []admissions.HospitalAdmission `json:"hospitalAdmissions"`
	Patient            models.Patient                 `json:"patient"`
}

func (a *AdmissionRoutes) AdmissionsByPregnancyHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "HospitalAdmissionsByPatientHandler"
	w.Header().Add("Content-Type", "application/json")
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		patientId := vars["patientId"]
		pregnancyId, err := strconv.Atoi(vars["pregnancyId"])
		if err != nil {
			log.WithFields(log.Fields{
				"handler":     handlerName,
				"method":      r.Method,
				"pregnancyId": vars["pregnancyId"],
			}).WithError(err).Error("pregnancy id is not a valid number")
			http.Error(w, "pregnancy id is not a valid number", http.StatusBadRequest)
			return
		}
		admissions, err := a.Admissions.FindByPregnancyId(pregnancyId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
			}).WithError(err).Error("error retrieving patient's hospital admissions: %+v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		patient, err := a.Patients.FindByPatientId(patientId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId":  patientId,
				"admissions": admissions,
				"handler":    handlerName,
			}).WithError(err).Error("error retrieving patient information")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response := admissionsResponse{
			HospitalAdmissions: admissions,
			Patient:            *patient,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{
				"patientId":  patientId,
				"admissions": admissions,
			}).WithError(err).Error("failed to marshal admissions response")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

type newAdmissionRequest struct {
	PatientId    int       `json:"patientId"`
	DateAdmitted time.Time `json:"dateAdmitted"`
	Facility     string    `json:"facility"`
	Reason       string    `json:"reason"`
}

func (a *AdmissionRoutes) AdmissionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	handlerName := "AdmissionsHandler"
	defer r.Body.Close()
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodPost:
		token := r.Context().Value("user").(auth.JwtToken)
		user := token.Email
		vars := mux.Vars(r)
		pregnancyId, err := strconv.Atoi(vars["pregnancyId"])
		if err != nil {
			log.WithFields(log.Fields{
				"handler":     handlerName,
				"method":      r.Method,
				"pregnancyId": vars["pregnancyId"],
			}).WithError(err).Error("pregnancy id must be a valid number")
			http.Error(w, "pregnancy id must be a valid number", http.StatusBadRequest)
			return
		}
		var req newAdmissionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.WithFields(log.Fields{
				"user": user,
			}).WithError(err).Error("error parsing request body when creating a new hospital admission")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		admission := admissions.HospitalAdmission{
			ID:           uuid.New().String(),
			PatientId:    req.PatientId,
			PregnancyId:  pregnancyId,
			DateAdmitted: req.DateAdmitted,
			Facility:     req.Facility,
			Reason:       req.Reason,
			CreatedAt:    time.Now(),
			UpdatedAt:    nil,
			CreatedBy:    user,
			UpdatedBy:    nil,
		}
		err = a.Admissions.Save(admission)
		if err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"admission": admission,
			}).WithError(err).Error("error when posting a request to create a hospital admission")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(admission); err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"admission": admission,
			}).WithError(err).Error("error marshalling response of newly created hospital admission")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPut:
		token := r.Context().Value("user").(auth.JwtToken)
		user := token.Email
		vars := mux.Vars(r)
		pregnancyId, err := strconv.Atoi(vars["pregnancyId"])
		if err != nil {
			log.WithFields(log.Fields{
				"handler":     handlerName,
				"method":      r.Method,
				"pregnancyId": vars["pregnancyId"],
			}).WithError(err).Error("pregnancy id is not a valid number")
			http.Error(w, "pregnancy id is not a valid number", http.StatusBadRequest)
			return
		}
		patientId, err := strconv.Atoi(vars["patientId"])
		if err != nil {
			log.WithFields(log.Fields{
				"handler":   handlerName,
				"method":    r.Method,
				"patientId": vars["patientId"],
			}).WithError(err).Error("patient is not a valid number", http.StatusBadRequest)
			return
		}
		var req admissions.HospitalAdmission
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.WithFields(log.Fields{
				"user": user,
			}).WithError(err).Error("error while parsing body for editing a hospital admission")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		now := time.Now()
		req.UpdatedBy = &user
		req.UpdatedAt = &now
		req.PregnancyId = pregnancyId
		req.PatientId = patientId
		err = a.Admissions.Update(req)
		if err != nil {
			log.WithFields(log.Fields{
				"user":        user,
				"requestBody": req,
			}).WithError(err).Error("failed to edit hospital admission")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(req); err != nil {
			log.WithFields(log.Fields{
				"user":        user,
				"requestBody": req,
			}).WithError(err).Error("error marshalling response for edited admission")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
