package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (a *App) InfantHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		motherId := vars["patientId"]
		id, err := strconv.Atoi(motherId)
		if err != nil {
			log.WithFields(log.Fields{
				"motherId": motherId,
			}).WithError(err).Error("motherId is not a valid number")
			http.Error(w, "motherId is not a valid number", http.StatusBadRequest)
			return
		}
		infant, err := a.AcsisDb.FindPregnancyInfant(id)
		if err != nil {
			log.WithFields(log.Fields{
				"motherId": motherId,
			}).WithError(err).Error("error retrieving pregnancy infant")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		result, err := json.Marshal(infant)
		if err != nil {
			log.WithFields(log.Fields{
				"infant": infant,
			}).WithError(err).Error("error marshalling infant data")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(result))
	}
}
