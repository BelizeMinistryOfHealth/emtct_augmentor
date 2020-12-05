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
	"moh.gov.bz/mch/emtct/internal/business/data/contraceptives"
	"moh.gov.bz/mch/emtct/internal/business/data/patient"
)

type ContraceptivesRoutes struct {
	Contraceptives contraceptives.Contraceptives
	Patients       patient.Patients
}

type newContraceptivesRequest struct {
	PatientId      int       `json:"patientId"`
	MchEncounterId int       `json:"mchEncounterId"`
	Contraceptive  string    `json:"contraceptive"`
	Comments       string    `json:"comments"`
	DateUsed       time.Time `json:"dateUsed"`
}

type contraceptivesResponse struct {
	Contraceptives []contraceptives.ContraceptiveUsed `json:"contraceptives"`
	Patient        patient.BasicInfo                  `json:"patient"`
}

func (a *ContraceptivesRoutes) ContraceptivesByPatientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		patientId := vars["patientId"]
		id, err := strconv.Atoi(patientId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
			}).Error("patientId provided is not a valid number")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		contraceptives, err := a.Contraceptives.FindByPatientId(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
			}).WithError(err).Error("failure retrieving a patient's contraceptives used")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		patient, err := a.Patients.FindBasicInfo(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId":      id,
				"contraceptives": contraceptives,
				"handler":        "ContraceptivesByPatientHandler",
			}).WithError(err).Error("error retrieving patient's information")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response := contraceptivesResponse{
			Contraceptives: contraceptives,
			Patient:        *patient,
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{
				"patientId":      patientId,
				"contraceptives": contraceptives,
				"response":       response,
			}).WithError(err).Error("failed to marshal list of contraceptives")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (a *ContraceptivesRoutes) ContraceptivesHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodPost:
		token := r.Context().Value("user").(app.JwtToken)
		user := token.Email
		var req newContraceptivesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.WithError(err).Error("failed to parse the body for creating a contraceptive")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		contraceptive := contraceptives.ContraceptiveUsed{
			Id:             uuid.New().String(),
			PatientId:      req.PatientId,
			MchEncounterId: req.MchEncounterId,
			Contraceptive:  req.Contraceptive,
			Comments:       req.Comments,
			DateUsed:       req.DateUsed,
			CreatedAt:      time.Now(),
			UpdatedAt:      nil,
			CreatedBy:      user,
			UpdatedBy:      nil,
		}
		err := a.Contraceptives.Create(contraceptive)
		if err != nil {
			log.WithFields(log.Fields{
				"user":          user,
				"request":       req,
				"contraceptive": contraceptive,
			}).WithError(err).Error("failed to create new contraceptive used")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(contraceptive); err != nil {
			log.WithFields(log.Fields{
				"user":          user,
				"contraceptive": contraceptive,
			}).WithError(err).Error("failed to marshal new contraceptive created")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPut:
		token := r.Context().Value("user").(app.JwtToken)
		user := token.Email
		var req contraceptives.ContraceptiveUsed
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.WithError(err).Error("failed to parse the body for creating a contraceptive")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		updated := time.Now()
		contraceptive := contraceptives.ContraceptiveUsed{
			Id:             req.Id,
			PatientId:      req.PatientId,
			MchEncounterId: req.MchEncounterId,
			Contraceptive:  req.Contraceptive,
			Comments:       req.Comments,
			DateUsed:       req.DateUsed,
			UpdatedAt:      &updated,
			UpdatedBy:      &user,
		}

		if err := a.Contraceptives.Edit(contraceptive); err != nil {
			log.WithFields(log.Fields{
				"user":          user,
				"request":       req,
				"contraceptive": contraceptive,
			}).WithError(err).Error("failed to create new contraceptive used")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(contraceptive); err != nil {
			log.WithFields(log.Fields{
				"user":          user,
				"contraceptive": contraceptive,
			}).WithError(err).Error("failed to marshal new contraceptive created")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
