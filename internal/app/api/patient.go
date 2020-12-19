package api

import (
	"encoding/json"
	"moh.gov.bz/mch/emtct/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/app"
)

type patientResponse struct {
	Patient     *models.Patient    `json:"patient"`
	Pregnancies []models.Pregnancy `json:"pregnancies"`
}

func (a *pregnancyRoutes) RetrievePatientHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		patientId := vars["id"]
		patient, err := a.Patient.FindByPatientId(patientId)
		if err != nil {
			log.WithFields(
				log.Fields{"request": r}).WithError(err).Error("could not find patient with specified id")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		v, err := a.Patient.GetPregnancies(patientId)
		if patient == nil {
			emptyResponse := patientResponse{
				Patient: nil,
			}
			_ = json.NewEncoder(w).Encode(emptyResponse)
			return
		}
		response := patientResponse{
			Patient:     patient,
			Pregnancies: v,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{"patientId": patientId, "response": response}).WithError(err).Error("error marshalling response")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

type arvsResponse struct {
	Arvs    []models.Prescription `json:"arvs"`
	Patient models.Patient        `json:"patient"`
}

func (a *pregnancyRoutes) ArvsHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "ArvsHandler"
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		method := "Get"
		token := r.Context().Value("user").(app.JwtToken)
		user := token.Email
		vars := mux.Vars(r)
		id := vars["patientId"]
		patientId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).
				WithError(err).
				Error("patient id is not a valid number")
			http.Error(w, "patientId is not a valid number", http.StatusBadRequest)
			return
		}
		// Find the pregnancy and the lmp so we can get the date bounds
		pregs, err := a.Pregnancies.FindCurrentPregnancy(patientId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).
				WithError(err).
				Error("error retrieving patient's current pregnancy")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if pregs == nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).
				WithError(err).
				Error("patient does not have a current pregnancy")
			http.Error(w, "patient does not have a current pregnancy", http.StatusNotFound)
			return
		}
		lmp := pregs.Lmp
		nextDate := lmp.Add(time.Hour * 24 * 7 * 54)
		arvs, err := a.Patient.FindArvsByPatient(id, *lmp, nextDate)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).
				WithError(err).
				Error("patient does not have a current pregnancy")
			http.Error(w, "patient does not have a current pregnancy", http.StatusNotFound)
			return
		}
		patientInfo, err := a.Patient.FindByPatientId(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).
				WithError(err).
				Error("error retrieving patient's basic info")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		arvsResponse := arvsResponse{
			Arvs:    arvs,
			Patient: *patientInfo,
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(arvsResponse); err != nil {
			log.WithFields(log.Fields{"patientId": patientId, "user": user}).
				WithError(err).
				Error("marshalling arvs response failed")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}
}

type treatmentResponse struct {
	Prescriptions []models.Prescription `json:"prescriptions"`
	Patient       models.Patient        `json:"patient"`
}

func (a *pregnancyRoutes) PatientSyphilisTreatmentHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "PatientSyphilisTreatmentHandler"
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		method := "GET"
		vars := mux.Vars(r)
		id := vars["patientId"]
		token := r.Context().Value("user").(app.JwtToken)
		user := token.Email
		patientId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).WithError(err).Error("patient id is not a valid number")
			http.Error(w, "patient id is not a valid number", http.StatusBadRequest)
			return
		}
		patientInfo, err := a.Patient.FindByPatientId(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"handler":   handlerName,
				"method":    method,
			}).WithError(err).Error("failed to retrieve patient basic info")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		pregs, err := a.Patient.GetPregnancies(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"handler":   handlerName,
				"method":    method,
			}).WithError(err).Error("failed to retrieve patient latest pregnancy")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if pregs == nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
				"user":      user,
			}).Error("patient has no pregnancy")
			http.Error(w, "patient has no pregnancy", http.StatusNoContent)
			return
		}
		preg := pregs[0]
		lmp := preg.Lmp
		endDate := lmp.Add(time.Hour * 24 * 7 * 52)
		treatments, err := a.Patient.FindSyphilisTreatment(id, lmp, &endDate)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"handler":   handlerName,
				"method":    method,
				"lmp":       lmp,
				"endDate":   endDate,
			}).WithError(err).Error("failed to retrieve patient syphilis treatment")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response := treatmentResponse{
			Prescriptions: treatments,
			Patient:       *patientInfo,
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"handler":   handlerName,
				"method":    method,
				"lmp":       lmp,
				"endDate":   endDate,
				"response":  response,
			}).WithError(err).Error("failed to encode syphilis treatment")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}
}
