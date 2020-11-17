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

type infantSyphilisScreeningResponse struct {
	Infant     models.Infant              `json:"infant"`
	Screenings []models.SyphilisScreening `json:"screenings"`
}

func (a *App) InfantSyphilisScreeningHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		infantId := vars["infantId"]
		id, err := strconv.Atoi(infantId)
		if err != nil {
			log.WithFields(log.Fields{
				"infantId": infantId,
				"handler":  "InfantSyphilisScreeningHandler",
			}).WithError(err).Error("infantId is not a valid number: %+v", err)
			http.Error(w, "infantId is not a valid number", http.StatusBadRequest)
			return
		}
		infantInfo, err := a.AcsisDb.FindInfant(id)
		if err != nil {
			log.WithFields(log.Fields{
				"infantId": id,
				"handler":  "InfantSyphilisScreeningHandler",
			}).WithError(err).Error("error retrieving infant information")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if infantInfo == nil {
			log.WithFields(log.Fields{
				"infantId": id,
				"handler":  "InfantSyphilisScreeningHandler",
			}).Error("no infant exists with given id")
			http.Error(w, fmt.Sprintf("infant with id %d does not exist", id), http.StatusNotFound)
			return
		}
		screenings, err := a.AcsisDb.FindInfantSyphilisScreenings(id, *infantInfo.Infant.Dob)
		if err != nil {
			log.WithFields(log.Fields{
				"infantId":   id,
				"handler":    "InfantSyphilisScreeningHandler",
				"infantInfo": infantInfo,
			}).WithError(err).Error("error retrieving syphilis screenings for infant")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response := infantSyphilisScreeningResponse{
			Infant:     *infantInfo,
			Screenings: screenings,
		}
		result, err := json.Marshal(response)
		if err != nil {
			log.WithFields(log.Fields{
				"infantId":   id,
				"screenings": screenings,
				"handler":    "InfantSyphilisScreeningHandler",
			}).WithError(err).Error("error marshalling screening data")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(result))
	}
}
