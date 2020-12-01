package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/app"
	"moh.gov.bz/mch/emtct/internal/business/data/admissions"
	"moh.gov.bz/mch/emtct/internal/business/data/patient"
)

type AdmissionRoutes struct {
	Admissions admissions.Admissions
	Patients   patient.Patients
}

type admissionsResponse struct {
	HospitalAdmissions []admissions.HospitalAdmission `json:"hospitalAdmissions"`
	Patient            patient.BasicInfo              `json:"patient"`
}

func (a *AdmissionRoutes) AdmissionsByPatientHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		patientId := vars["patientId"]
		id, err := strconv.Atoi(patientId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
			}).Error("patientId provided is not a valid number")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		admissions, err := a.Admissions.FindByPatientId(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
			}).WithError(err).Error("error retrieving patient's hospital admissions: %+v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		patient, err := a.Patients.FindBasicInfo(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId":  id,
				"admissions": admissions,
				"handler":    "HospitalAdmissionsByPatientHandler",
			}).WithError(err).Error("error retrieving patient information")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response := admissionsResponse{
			HospitalAdmissions: admissions,
			Patient:            *patient,
		}
		w.Header().Add("Content-Type", "application/json")
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
	PatientId      int       `json:"patientId"`
	DateAdmitted   time.Time `json:"dateAdmitted"`
	Facility       string    `json:"facility"`
	Reason         string    `json:"reason"`
	MchEncounterId int       `json:"mchEncounterId"`
}

func (a *AdmissionRoutes) AdmissionsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodPost:
		token := r.Context().Value("user").(app.JwtToken)
		user := token.Email
		var req newAdmissionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.WithFields(log.Fields{
				"user": user,
			}).WithError(err).Error("error parsing request body when creating a new hospital admission")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		admission := admissions.HospitalAdmission{
			Id:             uuid.New().String(),
			PatientId:      req.PatientId,
			MchEncounterId: req.MchEncounterId,
			DateAdmitted:   req.DateAdmitted,
			Facility:       req.Facility,
			Reason:         req.Reason,
			CreatedAt:      time.Now(),
			UpdatedAt:      nil,
			CreatedBy:      user,
			UpdatedBy:      nil,
		}
		err := a.Admissions.Create(admission)
		if err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"admission": admission,
			}).WithError(err).Error("error when posting a request to create a hospital admission")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(admission); err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"admission": admission,
			}).WithError(err).Error("error marshalling response of newly created hospital admission")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPut:
		vars := mux.Vars(r)
		admissionId := vars["admissionId"]
		token := r.Context().Value("user").(app.JwtToken)
		user := token.Email
		h, err := a.Admissions.FindById(admissionId)
		if err != nil {
			log.WithFields(log.Fields{
				"user":        user,
				"admissionId": admissionId,
			}).WithError(err).Error("error while editing hospital admission")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if h == nil {
			log.WithFields(log.Fields{
				"user":        user,
				"admissionId": admissionId,
			}).WithError(err).Error("tried to edit a non-existent hospital admission")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		var req admissions.HospitalAdmission
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.WithFields(log.Fields{
				"user":        user,
				"admissionId": admissionId,
			}).WithError(err).Error("error while parsing body for editing a hospital admission")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if req.Id != admissionId {
			log.WithFields(log.Fields{
				"user":        user,
				"admissionId": admissionId,
				"request":     req,
			}).WithError(err).Error("the request's id must match the admission id in the resource url")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		now := time.Now()
		req.UpdatedBy = &user
		req.UpdatedAt = &now
		err = a.Admissions.Edit(req)
		if err != nil {
			log.WithFields(log.Fields{
				"user":        user,
				"requestBody": req,
				"admissionId": admissionId,
			}).WithError(err).Error("failed to edit hospital admission")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(req); err != nil {
			log.WithFields(log.Fields{
				"user":        user,
				"admissionId": admissionId,
				"requestBody": req,
			}).WithError(err).Error("error marshalling response for edited admission")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
