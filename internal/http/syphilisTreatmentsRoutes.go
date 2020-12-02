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

type infantTreatmentResponse struct {
	Prescriptions []models.Prescription `json:"prescriptions"`
	Infant        models.Infant         `json:"infant"`
}

func (a *App) InfantSyphilisTreatmentHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		id := vars["infantId"]
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		infantId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{
				"infantId": id,
				"user":     user,
			}).WithError(err).Error("patient id is not a valid number")
			http.Error(w, "patient id is not a valid number", http.StatusBadRequest)
			return
		}
		treatments, err := a.AcsisDb.FindInfantSyphilisTreatment(infantId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": infantId,
				"user":      user,
			}).WithError(err).Error("error retrieving syphilis treatment")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		infant, err := a.AcsisDb.FindInfant(infantId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": infantId,
				"user":      user,
				"handler":   "SyphilisTreatmentHandler",
			}).WithError(err).Error("error retrieving patient information")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response := infantTreatmentResponse{
			Prescriptions: treatments,
			Infant:        *infant,
		}
		result, err := json.Marshal(response)
		if err != nil {
			log.WithFields(log.Fields{
				"response":  response,
				"user":      user,
				"patientId": infantId,
				"handler":   "SyphilisTreatmentHandler",
			}).WithError(err).Error("error marshaling syphilis treatment response")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, string(result))
	}
}

type newSyphilisTreatmentRequest struct {
	PatientId  int       `json:"patientId"`
	Medication string    `json:"medication"`
	Dosage     string    `json:"dosage"`
	Comments   string    `json:"comments"`
	Date       time.Time `json:"date"`
}

func (a *App) PartnerSyphilisTreatmentHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "PartnerSyphilisTreatmentHandler"
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		vars := mux.Vars(r)
		id := vars["patientId"]
		patientId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"patientId": id,
			}).WithError(err).Error("patient id must be a number")
			http.Error(w, "patient id must be a valid number", http.StatusBadRequest)
			return
		}
		treatments, err := a.Db.FindPartnerSyphilisTreatments(patientId)
		if err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"patientId": patientId,
				"handler":   handlerName,
			}).WithError(err).Error("error while finding partner's syphilis treatment")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		patient, err := a.AcsisDb.FindPatientBasicInfo(patientId)
		if err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"patientId": patientId,
				"handler":   handlerName,
			}).WithError(err).Error("error querying patient's basic info")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response := map[string]interface{}{
			"patient":    patient,
			"treatments": treatments,
		}
		w.Header().Add("Content-type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{
				"user":     user,
				"response": response,
				"handler":  handlerName,
			}).WithError(err).Error("error encoding response")
		}
	case http.MethodPost:
		defer r.Body.Close()
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		vars := mux.Vars(r)
		id := vars["patientId"]
		patientId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"handler":   handlerName,
				"patientId": id,
			}).WithError(err).Error("patient id must be a valid number")
			http.Error(w, "patient id must be a valid number", http.StatusBadRequest)
			return
		}
		var treatmentReq newSyphilisTreatmentRequest
		if err := json.NewDecoder(r.Body).Decode(&treatmentReq); err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"patientId": patientId,
				"body":      r.Body,
				"handler":   handlerName,
			}).WithError(err).Error("error decoding request")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		location, _ := time.LoadLocation("Local")
		treatment := models.SyphilisTreatment{
			Id:         uuid.New().String(),
			PatientId:  treatmentReq.PatientId,
			Medication: treatmentReq.Medication,
			Dosage:     treatmentReq.Dosage,
			Comments:   treatmentReq.Comments,
			Date:       treatmentReq.Date.In(location),
			CreatedBy:  user,
			CreatedAt:  time.Now(),
		}
		if err := a.Db.AddPartnerSyphilisTreatment(treatment); err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"request":   treatmentReq,
				"treatment": treatment,
				"handler":   handlerName,
			}).WithError(err).Error("error adding a partner's syphilis treatment")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(treatment); err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"treatment": treatment,
				"handler":   handlerName,
			}).WithError(err).Error("error encoding response")
		}
	case http.MethodPut:
		defer r.Body.Close()
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		var treatment models.SyphilisTreatment
		if err := json.NewDecoder(r.Body).Decode(&treatment); err != nil {
			log.WithFields(log.Fields{
				"user":    user,
				"body":    r.Body,
				"handler": handlerName,
			}).WithError(err).Error("error decoding the request")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		treatment.UpdatedBy = user
		today := time.Now()
		treatment.UpdatedAt = &today
		location, _ := time.LoadLocation("Local")
		treatment.Date = treatment.Date.In(location)
		if err := a.Db.UpdatePartnerSyphilisTreatment(treatment); err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"treatment": treatment,
				"handler":   handlerName,
			}).WithError(err).Error("failed to update treatment")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(treatment); err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"treatment": treatment,
				"handler":   handlerName,
			}).WithError(err).Error("")
		}
	}
}
