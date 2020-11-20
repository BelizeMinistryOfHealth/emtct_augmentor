package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/models"
)

type treatmentResponse struct {
	Prescriptions []models.Prescription   `json:"prescriptions"`
	Patient       models.PatientBasicInfo `json:"patient"`
}

func (a *App) SyphilisTreatmentHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		id := vars["patientId"]
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		patientId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
			}).WithError(err).Error("patient id is not a valid number")
			http.Error(w, "patient id is not a valid number", http.StatusBadRequest)
			return
		}
		treatments, err := a.AcsisDb.FindPatientSyphilisTreatment(patientId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
				"user":      user,
			}).WithError(err).Error("error retrieving syphilis treatment")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		patient, err := a.AcsisDb.FindPatientBasicInfo(patientId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
				"user":      user,
				"handler":   "SyphilisTreatmentHandler",
			}).WithError(err).Error("error retrieving patient information")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response := treatmentResponse{
			Prescriptions: treatments,
			Patient:       *patient,
		}
		result, err := json.Marshal(response)
		if err != nil {
			log.WithFields(log.Fields{
				"response":  response,
				"user":      user,
				"patientId": patientId,
				"handler":   "SyphilisTreatmentHandler",
			}).WithError(err).Error("error marshaling syphilis treatment response")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, string(result))
	}
}

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
	}
}
