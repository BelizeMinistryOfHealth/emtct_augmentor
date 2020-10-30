package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/models"
)

type NewHospitalAdmissionRequest struct {
	PatientId    int       `json:"patientId"`
	DateAdmitted time.Time `json:"dateAdmitted"`
	Facility     string    `json:"facility"`
}

func (a *App) CreateHospitalAdmissionHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodPost:
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		var req NewHospitalAdmissionRequest
		err := parseBody(r.Body, &req)
		if err != nil {
			log.WithFields(log.Fields{
				"user": user,
			}).WithError(err).Error("error parsing request body when creating a new hospital admission")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		admission := models.HospitalAdmission{
			Id:           uuid.New().String(),
			PatientId:    req.PatientId,
			DateAdmitted: req.DateAdmitted,
			Facility:     req.Facility,
			CreatedAt:    time.Now(),
			UpdatedAt:    nil,
			CreatedBy:    user,
			UpdatedBy:    nil,
		}
		err = a.Db.CreateHospitalAdmission(admission)
		if err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"admission": admission,
			}).WithError(err).Error("error when posting a request to create a hospital admission")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(admission)
		if err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"admission": admission,
			}).WithError(err).Error("error marshalling response of newly created hospital admission")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, string(resp))
	}
}

func (a *App) HospitalAdmissionsByPatientHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
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
		admissions, err := a.Db.HospitalAdmissionsByPatientId(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
			}).WithError(err).Error("error retrieving patient's hospital admissions: %+v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		results, err := json.Marshal(admissions)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId":  patientId,
				"admissions": admissions,
			}).WithError(err).Error("failed to marshal admissions response")
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, string(results))
	default:
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func (a *App) HospitalAdmissionsApiHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodPut:
		vars := mux.Vars(r)
		admissionId := vars["admissionId"]
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		h, err := a.Db.FindHospitalAdmissionById(admissionId)
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
		var req models.HospitalAdmission
		err = parseBody(r.Body, &req)
		if err != nil {
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
		err = a.Db.EditHospitalAdmission(req)
		if err != nil {
			log.WithFields(log.Fields{
				"user":        user,
				"requestBody": req,
				"admissionId": admissionId,
			}).WithError(err).Error("failed to edit hospital admission")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		result, err := json.Marshal(req)
		if err != nil {
			log.WithFields(log.Fields{
				"user":        user,
				"admissionId": admissionId,
				"requestBody": req,
			}).WithError(err).Error("error marshalling response for edited admission")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, string(result))
	default:
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}
