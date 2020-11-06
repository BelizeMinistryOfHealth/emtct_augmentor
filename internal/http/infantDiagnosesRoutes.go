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
		result, err := json.Marshal(diagnoses)
		if err != nil {
			log.WithFields(log.Fields{
				"motherId":  id,
				"user":      user,
				"diagnoses": diagnoses,
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
