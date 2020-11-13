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
