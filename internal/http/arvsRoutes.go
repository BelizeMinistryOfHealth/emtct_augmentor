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

type ArvsResponse struct {
	Arvs    []models.Prescription   `json:"arvs"`
	Patient models.PatientBasicInfo `json:"patient"`
}

func (a *App) ArvsHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		id := vars["patientId"]
		patientId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{"patientId": id}).
				WithError(err).
				Error("patient id is not a valid number")
			http.Error(w, "patientId is not a valid number", http.StatusBadRequest)
			return
		}
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		arvs, err := a.AcsisDb.FindArvsByPatient(patientId)
		if err != nil {
			log.WithFields(log.Fields{"patientId": id, "user": user}).
				WithError(err).
				Error("error retrieving arvs for patient")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		patientInfo, err := a.AcsisDb.FindPatientBasicInfo(patientId)
		if err != nil {
			log.WithFields(log.Fields{"patientId": patientId, "user": user}).
				WithError(err).
				Error("error retrieving patient basic info")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		arvsResponse := ArvsResponse{
			Arvs:    arvs,
			Patient: *patientInfo,
		}
		result, err := json.Marshal(arvsResponse)
		if err != nil {
			log.WithFields(log.Fields{"patientId": patientId, "user": user}).
				WithError(err).
				Error("marshalling arvs response failed")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(result))
	}
}
