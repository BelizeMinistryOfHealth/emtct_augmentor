package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/models"
)

type NewContraceptiveUsedRequest struct {
	PatientId      int       `json:"patientId"`
	MchEncounterId int       `json:"mchEncounterId"'`
	Contraceptive  string    `json:"contraceptive"`
	Comments       string    `json:"comments"`
	DateUsed       time.Time `json:"dateUsed"`
}

func (a *App) CreateContraceptiveUsedHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodPost:
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		var req NewContraceptiveUsedRequest
		err := parseBody(r.Body, &req)
		if err != nil {
			log.WithError(err).Error("failed to parse the body for creating a contraceptive")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		contraceptive := models.ContraceptiveUsed{
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
		err = a.Db.CreateContraceptiveUsed(contraceptive)
		if err != nil {
			log.WithFields(log.Fields{
				"user":          user,
				"request":       req,
				"contraceptive": contraceptive,
			}).WithError(err).Error("failed to create new contraceptive used")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(contraceptive)
		if err != nil {
			log.WithFields(log.Fields{
				"user":          user,
				"contraceptive": contraceptive,
			}).WithError(err).Error("failed to marshal new contraceptive created")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, string(resp))
	}
}

type contraceptivesResponse struct {
	Contraceptives []models.ContraceptiveUsed `json:"contraceptives"`
	Patient        models.PatientBasicInfo    `json:"patient"`
}

func (a *App) ContraceptivesByPatientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
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
	contraceptives, err := a.Db.ContraceptivesUsedByPatientId(id)
	if err != nil {
		log.WithFields(log.Fields{
			"patientId": patientId,
		}).WithError(err).Error("failure retrieving a patient's contraceptives used")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	patient, err := a.AcsisDb.FindPatientBasicInfo(id)
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
	results, err := json.Marshal(response)
	if err != nil {
		log.WithFields(log.Fields{
			"patientId":      patientId,
			"contraceptives": contraceptives,
			"response":       response,
		}).WithError(err).Error("failed to marshal list of contraceptives")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(results))
}

func (a *App) ContraceptivesApiHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodPut:
		vars := mux.Vars(r)
		contraceptiveId := vars["contraceptiveId"]
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		c, err := a.Db.FindContraceptiveById(contraceptiveId)
		if err != nil {
			log.WithFields(log.Fields{
				"contraceptiveId": contraceptiveId,
				"user":            user,
			}).WithError(err).Error("error searching for contraceptive when performing an update")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		var req models.ContraceptiveUsed
		err = parseBody(r.Body, &req)
		if err != nil {
			log.WithFields(log.Fields{
				"contraceptiveId": contraceptiveId,
				"user":            user,
				"contraceptive":   c,
			}).WithError(err).Error("failed to parse request body for updating a contraceptive")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if req.Id != contraceptiveId {
			log.WithFields(log.Fields{
				"contraceptiveId": contraceptiveId,
				"user":            user,
				"request":         c,
			}).Error("the contraceptive id and the request id must be equal when updating a contraceptive")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		now := time.Now()
		req.UpdatedBy = &user
		req.UpdatedAt = &now
		err = a.Db.EditContraceptiveUsed(req)
		if err != nil {
			log.WithFields(log.Fields{
				"user":            user,
				"contraceptiveId": contraceptiveId,
				"request":         req,
			}).WithError(err).Error("error trying to update a contraceptive")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		result, err := json.Marshal(req)
		if err != nil {
			log.WithFields(log.Fields{
				"contraceptiveId": contraceptiveId,
				"user":            user,
				"request":         req,
			}).WithError(err).Error("failed to marshal contraceptive when editing a contraceptive")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, string(result))
	default:
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		return
	}
}
