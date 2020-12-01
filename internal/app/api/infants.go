package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/app"
	"moh.gov.bz/mch/emtct/internal/business/data/infant"
	"moh.gov.bz/mch/emtct/internal/business/data/pregnancy"
	"moh.gov.bz/mch/emtct/internal/db"
)

type InfantRoutes struct {
	Infant      infant.Infants
	AcsisDb     db.AcsisDb
	Pregnancies pregnancy.Pregnancies
}

func (i InfantRoutes) InfantHandlers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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
		// Find current pregnancy
		preg, err := i.Pregnancies.FindLatest(id)
		if err != nil {
			log.WithFields(log.Fields{
				"motherId": id,
			}).WithError(err).Error("error retrieving this patient's latest pregnancy")
			http.Error(w, "could not retrieve the mother's latest pregnancy", http.StatusInternalServerError)
			return
		}
		infant, err := i.Infant.FindPregnancyInfant(*preg)
		if err != nil {
			log.WithFields(log.Fields{
				"motherId": motherId,
			}).WithError(err).Error("error retrieving pregnancy infant")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if infant == nil {
			log.WithFields(log.Fields{
				"motherId": motherId,
			}).Error("no infant exists for current pregnancy")
			http.Error(w, "no infant exists for relevant pregnancy", http.StatusNotFound)
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

type infantDiagnosesResponse struct {
	Diagnoses []infant.Diagnoses `json:"diagnoses"`
	Infant    infant.Infant      `json:"infant"`
}

func (i InfantRoutes) InfantDiagnosesHandler(w http.ResponseWriter, r *http.Request) {
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
		token := r.Context().Value("user").(app.JwtToken)
		user := token.Email
		diagnoses, err := i.Infant.FindInfantDiagnoses(infantId)
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
			diagnoses = []infant.Diagnoses{}
		}

		// Find current pregnancy
		preg, err := i.Pregnancies.FindLatest(motherId)
		if err != nil {
			log.WithFields(log.Fields{
				"motherId": id,
			}).WithError(err).Error("error retrieving this patient's latest pregnancy")
			http.Error(w, "could not retrieve the mother's latest pregnancy", http.StatusInternalServerError)
			return
		}
		infantInfo, err := i.Infant.FindPregnancyInfant(*preg)
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
