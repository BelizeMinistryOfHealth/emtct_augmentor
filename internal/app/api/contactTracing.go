package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/app"
	"moh.gov.bz/mch/emtct/internal/business/data/contactTracing"
	"moh.gov.bz/mch/emtct/internal/business/data/patient"
)

type ContactTracingRoutes struct {
	ContactTracings contactTracing.ContactTracings
	Patient         patient.Patients
}

type contactTracingRequest struct {
	PatientId  int       `json:"patientId"`
	Test       string    `json:"test"`
	TestResult string    `json:"testResult"`
	Comments   string    `json:"comments"`
	Date       time.Time `json:"date"`
}

func (a *ContactTracingRoutes) ContactTracingHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "ContactTracingHandler"
	w.Header().Add("Content-Type", "application/json")
	defer r.Body.Close()

	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodPost:
		token := r.Context().Value("user").(app.JwtToken)
		user := token.Email
		vars := mux.Vars(r)
		pregnancyId, err := strconv.Atoi(vars["pregnancyId"])
		if err != nil {
			log.WithFields(log.Fields{
				"handler":     handlerName,
				"method":      r.Method,
				"pregnancyId": vars["pregnancyId"],
			}).WithError(err).Error("pregnancy id must be a valid number")
			http.Error(w, "pregnancy id must be a valid number", http.StatusBadRequest)
			return
		}
		var request contactTracingRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.WithFields(log.Fields{
				"user":    user,
				"body":    r.Body,
				"handler": handlerName,
				"method":  r.Method,
			}).WithError(err).Error("error decoding json payload")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		location, _ := time.LoadLocation("Local")
		contactTracing := contactTracing.ContactTracing{
			ID:          uuid.New().String(),
			PatientId:   request.PatientId,
			PregnancyId: pregnancyId,
			Test:        request.Test,
			TestResult:  request.TestResult,
			Comments:    request.Comments,
			Date:        request.Date.In(location),
			CreatedBy:   user,
			CreatedAt:   time.Now(),
		}
		if err := a.ContactTracings.Create(contactTracing); err != nil {
			log.WithFields(log.Fields{
				"user":           user,
				"contactTracing": contactTracing,
				"handler":        handlerName,
				"method":         r.Method,
			}).WithError(err).Error("error when inserting contact tracing")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(contactTracing); err != nil {
			log.WithFields(log.Fields{
				"user":           user,
				"contactTracing": contactTracing,
				"handler":        handlerName,
				"method":         r.Method,
			}).WithError(err).Error("error encoding contact tracing response")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodGet:
		token := r.Context().Value("user").(app.JwtToken)
		user := token.Email
		vars := mux.Vars(r)
		id := vars["patientId"]
		pregnancyId, err := strconv.Atoi(vars["pregnancyId"])
		if err != nil {
			log.WithFields(log.Fields{
				"handler":     handlerName,
				"method":      r.Method,
				"user":        user,
				"pregnancyId": vars["pregnancyId"],
			}).WithError(err).Error("pregnancy id must be a valid number")
			http.Error(w, "pregnancy id must be a valid number", http.StatusBadRequest)
			return
		}
		contacts, err := a.ContactTracings.FindByPregnancyId(pregnancyId)
		if err != nil {
			log.WithFields(log.Fields{
				"user":        user,
				"pregnancyId": pregnancyId,
				"handler":     handlerName,
				"method":      r.Method,
			}).WithError(err).Error("error retrieving patient's contact tracings")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		patient, err := a.Patient.FindByPatientId(id)
		if err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"patientId": id,
				"handler":   handlerName,
				"method":    r.Method,
			}).WithError(err).Error("error retrieving patient's basic information")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response := map[string]interface{}{
			"patient":        patient,
			"contactTracing": contacts,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{
				"user":     user,
				"response": response,
				"handler":  handlerName,
			}).WithError(err).Error("error encoding contact tracings")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPut:
		token := r.Context().Value("user").(app.JwtToken)
		user := token.Email
		today := time.Now()
		var contactTracing contactTracing.ContactTracing
		if err := json.NewDecoder(r.Body).Decode(&contactTracing); err != nil {
			log.WithFields(log.Fields{
				"user":    user,
				"handler": handlerName,
				"body":    r.Body,
			}).WithError(err).Error("error decoding payload")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		contactTracing.UpdatedBy = user
		contactTracing.UpdatedAt = &today
		location, _ := time.LoadLocation("Local")
		contactTracing.Date = contactTracing.Date.In(location)
		if err := a.ContactTracings.Update(contactTracing); err != nil {
			log.WithFields(log.Fields{
				"user":           user,
				"contactTracing": contactTracing,
				"handler":        handlerName,
			}).WithError(err).Error("error updaging contact tracing record")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(contactTracing); err != nil {
			log.WithFields(log.Fields{
				"user":           user,
				"contactTracing": contactTracing,
				"handler":        handlerName,
			}).WithError(err).Error("error encoding updated contact tracing")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
