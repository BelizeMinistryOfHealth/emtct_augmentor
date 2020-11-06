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

type obstetricHistoryResponse struct {
	ObstetricHistory []models.ObstetricHistory `json:"obstetricHistory"`
	Patient          models.PatientBasicInfo   `json:"patient"`
}

func (a *App) ObstetricHistoryHandler(w http.ResponseWriter, r *http.Request) {
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
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		obstetricHistory, err := a.AcsisDb.FindObstetricHistory(patientId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
				"user":      user,
				"handler":   "ObstetricHistoryHandler",
			}).WithError(err).Error("error retrieving obstetric history")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		patientInfo, err := a.AcsisDb.FindPatientBasicInfo(patientId)
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
			obstetricHistory = []models.ObstetricHistory{}
		}
		response := obstetricHistoryResponse{
			ObstetricHistory: obstetricHistory,
			Patient:          *patientInfo,
		}
		result, err := json.Marshal(response)
		if err != nil {
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
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(result))
	}
}
