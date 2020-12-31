package api

import (
	"encoding/json"
	"moh.gov.bz/mch/emtct/internal/auth"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/business/data/partners"
	"moh.gov.bz/mch/emtct/internal/business/data/patient"
	"moh.gov.bz/mch/emtct/internal/business/data/prescription"
)

type partnersRoutes struct {
	Patient  patient.Patients
	Partners partners.Partners
}

type newSyphilisTreatmentRequest struct {
	PatientId  int       `json:"patientId"`
	Medication string    `json:"medication"`
	Dosage     string    `json:"dosage"`
	Comments   string    `json:"comments"`
	Date       time.Time `json:"date"`
}

func (p *partnersRoutes) SyphilisTreatmentHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "SyphilisTreatmentHandler"
	w.Header().Add("Content-type", "application/json")
	defer r.Body.Close()
	vars := mux.Vars(r)
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		token := r.Context().Value("user").(auth.JwtToken)
		user := token.Email
		id := vars["patientId"]
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
		treatments, err := p.Partners.FindByPregnancyId(pregnancyId)
		if err != nil {
			log.WithFields(log.Fields{
				"user":        user,
				"pregnancyId": pregnancyId,
				"handler":     handlerName,
				"method":      r.Method,
			}).WithError(err).Error("error while finding partner's syphilis treatment")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		patient, err := p.Patient.FindByPatientId(id)
		if err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"patientId": id,
				"handler":   handlerName,
				"method":    r.Method,
			}).WithError(err).Error("error querying patient's basic info")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response := map[string]interface{}{
			"patient":    patient,
			"treatments": treatments,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{
				"user":     user,
				"response": response,
				"handler":  handlerName,
				"method":   r.Method,
			}).WithError(err).Error("error encoding response")
		}
	case http.MethodPost:
		token := r.Context().Value("user").(auth.JwtToken)
		user := token.Email
		id := vars["patientId"]
		pregnancyId, err := strconv.Atoi(vars["pregnancyId"])
		if err != nil {
			log.WithFields(log.Fields{
				"user":        user,
				"handler":     handlerName,
				"patientId":   id,
				"pregnancyId": vars["pregnancyId"],
				"method":      r.Method,
			}).WithError(err).Error("patient id must be a valid number")
			http.Error(w, "patient id must be a valid number", http.StatusBadRequest)
			return
		}
		var treatmentReq newSyphilisTreatmentRequest
		if err := json.NewDecoder(r.Body).Decode(&treatmentReq); err != nil {
			log.WithFields(log.Fields{
				"user":    user,
				"body":    r.Body,
				"handler": handlerName,
				"method":  r.Method,
			}).WithError(err).Error("error decoding request")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		location, _ := time.LoadLocation("Local")
		treatment := prescription.SyphilisTreatment{
			ID:          uuid.New().String(),
			PatientId:   treatmentReq.PatientId,
			PregnancyId: pregnancyId,
			Medication:  treatmentReq.Medication,
			Dosage:      treatmentReq.Dosage,
			Comments:    treatmentReq.Comments,
			Date:        treatmentReq.Date.In(location),
			CreatedBy:   user,
			CreatedAt:   time.Now(),
		}
		if err := p.Partners.Create(treatment); err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"request":   treatmentReq,
				"treatment": treatment,
				"handler":   handlerName,
				"method":    r.Method,
			}).WithError(err).Error("error adding a partner's syphilis treatment")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(treatment); err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"treatment": treatment,
				"handler":   handlerName,
				"method":    r.Method,
			}).WithError(err).Error("error encoding response")
		}
	case http.MethodPut:
		token := r.Context().Value("user").(auth.JwtToken)
		user := token.Email
		var treatment prescription.SyphilisTreatment
		if err := json.NewDecoder(r.Body).Decode(&treatment); err != nil {
			log.WithFields(log.Fields{
				"user":    user,
				"body":    r.Body,
				"handler": handlerName,
			}).WithError(err).Error("error decoding the request")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		treatment.UpdatedBy = user
		today := time.Now()
		treatment.UpdatedAt = &today
		location, _ := time.LoadLocation("Local")
		treatment.Date = treatment.Date.In(location)
		if err := p.Partners.Update(treatment); err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"treatment": treatment,
				"handler":   handlerName,
			}).WithError(err).Error("failed to update treatment")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(treatment); err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"treatment": treatment,
				"handler":   handlerName,
			}).WithError(err).Error("")
		}
	}
}
