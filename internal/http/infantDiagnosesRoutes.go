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

type infantDiagnosesResponse struct {
	Diagnoses []models.InfantDiagnoses `json:"diagnoses"`
	Patient   models.PatientBasicInfo  `json:"patient"`
}

func (a *App) InfantDiagnosesHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		id := vars["motherId"]
		motherId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{"motherId": id}).WithError(err).
				Error("the mother id is not a valid number")
			http.Error(w, "motherId must be a numeric value", http.StatusBadRequest)
			return
		}
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		diagnoses, err := a.AcsisDb.InfantDiagnoses(motherId)
		if err != nil {
			log.WithFields(log.Fields{
				"motherId": motherId,
				"user":     user,
				"handler":  "InfantDiagnosesHandler",
			}).
				WithError(err).
				Error("error while fetching infant diagnoses")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// Return an empty array if no results are found
		if diagnoses == nil {
			diagnoses = []models.InfantDiagnoses{}
		}
		patientInfo, err := a.AcsisDb.FindPatientBasicInfo(motherId)
		if err != nil {
			log.WithFields(log.Fields{
				"motherId": motherId,
				"user":     user,
			}).
				WithError(err).
				Error("error retrieving patient basic info")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		patient := models.PatientBasicInfo{}
		if patientInfo != nil {
			patient = *patientInfo
		}

		response := infantDiagnosesResponse{
			Diagnoses: diagnoses,
			Patient:   patient,
		}
		result, err := json.Marshal(response)
		if err != nil {
			log.WithFields(log.Fields{
				"motherId": id,
				"user":     user,
				"response": response,
			}).
				WithError(err).
				Error("error while marshalling the infant diagnoses")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(result))
	}
}
