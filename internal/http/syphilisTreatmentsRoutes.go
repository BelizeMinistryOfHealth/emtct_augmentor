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
