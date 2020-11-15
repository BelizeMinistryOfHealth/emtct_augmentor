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
	Infant    models.Infant            `json:"infant"`
}

func (a *App) InfantDiagnosesHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		id := vars["infantId"]
		infantId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{"infantId": id}).WithError(err).
				Error("the infant id is not a valid number")
			http.Error(w, "infant id must be a numeric value", http.StatusBadRequest)
			return
		}
		mId := vars["motherId"]
		motherId, err := strconv.Atoi(mId)
		if err != nil {
			log.WithFields(log.Fields{"motherId": id}).WithError(err).
				Error("the mother id is not a valid number")
			http.Error(w, "mother id must be a numeric value", http.StatusBadRequest)
			return
		}
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		diagnoses, err := a.AcsisDb.FindInfantDiagnoses(infantId)
		if err != nil {
			log.WithFields(log.Fields{
				"infantId": infantId,
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

		infantInfo, err := a.AcsisDb.FindPregnancyInfant(motherId)
		if err != nil {
			log.WithFields(log.Fields{
				"infantId": infantId,
				"motherId": motherId,
				"user":     user,
			}).WithError(err).Error("could not find infant info")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		response := infantDiagnosesResponse{
			Diagnoses: diagnoses,
			Infant:    *infantInfo,
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
