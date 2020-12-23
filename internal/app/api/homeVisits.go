package api

import (
	"encoding/json"
	"fmt"
	"moh.gov.bz/mch/emtct/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/app"
	"moh.gov.bz/mch/emtct/internal/business/data/homeVisits"
	"moh.gov.bz/mch/emtct/internal/business/data/patient"
)

type HomeVisitRoutes struct {
	HomeVisits homeVisits.HomeVisits
	Patients   patient.Patients
}

type homeVisitResponse struct {
	HomeVisits []homeVisits.HomeVisit `json:"homeVisits"`
	Patient    models.Patient         `json:"patient"`
}

func (h HomeVisitRoutes) FindByPregnancyHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "FindByPregnancyHandler"
	w.Header().Add("Content-Type", "application/json")
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		patientId := vars["patientId"]
		id := vars["pregnancyId"]
		pregnancyId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{
				"handler":     handlerName,
				"method":      r.Method,
				"pregnancyId": pregnancyId,
			}).WithError(err).Error("pregnancy id must be a valid number: %w", err)
			http.Error(w, "pregnancy id must be a valid number", http.StatusInternalServerError)
			return
		}
		patient, err := h.Patients.FindByPatientId(patientId)
		homeVisits, err := h.HomeVisits.FindByPregnancyId(pregnancyId)
		if err != nil {
			log.WithFields(log.Fields{
				"handler":     handlerName,
				"method":      r.Method,
				"pregnancyId": pregnancyId,
			}).WithError(err).Error("")
		}
		resp := homeVisitResponse{
			HomeVisits: homeVisits,
			Patient:    *patient,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.WithFields(log.Fields{
				"handler":  handlerName,
				"method":   r.Method,
				"response": resp,
			}).WithError(err).Error("failed to encode response for home visits")
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

type homeVisitRequest struct {
	ID          string    `json:"id"`
	Reason      string    `json:"reason"`
	Comments    string    `json:"comments"`
	DateOfVisit time.Time `json:"dateOfVisit"`
}

type newHomeVisitRequest struct {
	PatientId      int       `json:"patientId"`
	Reason         string    `json:"reason"`
	Comments       string    `json:"comments"`
	DateOfVisit    time.Time `json:"dateOfVisit"`
	MchEncounterId int       `json:"mchEncounterId"`
}

func (h HomeVisitRoutes) HomeVisitsHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "HomeVisitsHandler"
	defer r.Body.Close()
	w.Header().Add("Content-Type", "application/json")

	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		id := vars["homeVisitId"]
		visit, err := h.HomeVisits.FindById(id)
		if err != nil {
			log.WithFields(log.Fields{"homeVisitId": id}).
				WithError(err).
				Error("error retrieving home visit from the database")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(visit); err != nil {
			log.WithFields(log.Fields{"homeVisit": visit, "homeVisitId": id}).
				WithError(err).
				Error("failed to marshal home visit")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPut:
		token := r.Context().Value("user").(app.JwtToken)
		user := token.Email
		var req homeVisitRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.WithFields(log.Fields{
				"user":    user,
				"handler": handlerName,
				"Method":  r.Method,
			}).WithError(err).Error("failed to decode the home visit request")
			http.Error(w, "could not decode your request", http.StatusInternalServerError)
			return
		}
		err := h.editHomeVisit(user, req)
		if err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"homeVisit": req,
			}).WithError(err).Error("error editing home visit")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(req); err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"homeVisit": req,
			}).
				WithError(err).
				Error("error marshalling home visit into json")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		token := r.Context().Value("user").(app.JwtToken)
		user := token.Email
		var req newHomeVisitRequest
		vars := mux.Vars(r)
		pregnancyId, err := strconv.Atoi(vars["pregnancyId"])
		if err != nil {
			log.WithFields(log.Fields{
				"handler":     handlerName,
				"method":      r.Method,
				"pregnancyId": vars["pregnancyId"],
			}).WithError(err).Error("pregnancy id is not a valid number")
			http.Error(w, "prengnancy id is not a valid number", http.StatusInternalServerError)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.WithFields(log.Fields{
				"user": user,
			}).WithError(err).Error("failed to parse body for creating a home visit")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		visit, err := h.createHomeVisit(user, pregnancyId, req)
		if err != nil {
			log.WithFields(log.Fields{
				"request": req,
			}).WithError(err).Error("error creating home visit")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(visit); err != nil {
			log.WithFields(log.Fields{
				"homeVisit": visit,
			}).WithError(err).Error("failed to marshal the newly created home visit")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}
}

func (h HomeVisitRoutes) editHomeVisit(user string, r homeVisitRequest) error {

	v, err := h.HomeVisits.FindById(r.ID)
	if err != nil {
		return fmt.Errorf("home visit with given id does not exist: %+v", err)
	}
	err = h.HomeVisits.Save(homeVisits.HomeVisit{
		Id:          v.Id,
		PatientId:   v.PatientId,
		PregnancyId: v.PregnancyId,
		Reason:      r.Reason,
		Comments:    r.Comments,
		DateOfVisit: r.DateOfVisit,
		CreatedAt:   v.CreatedAt,
		UpdatedAt:   v.UpdatedAt,
		CreatedBy:   v.CreatedBy,
		UpdatedBy:   &user,
	})
	return err
}

func (h HomeVisitRoutes) createHomeVisit(user string, pregnancyId int, r newHomeVisitRequest) (*homeVisits.HomeVisit, error) {
	id := uuid.New().String()
	if len(user) == 0 {
		return nil, fmt.Errorf("user did not provide an email")
	}

	visit := homeVisits.HomeVisit{
		Id:          id,
		PatientId:   r.PatientId,
		PregnancyId: pregnancyId,
		Reason:      r.Reason,
		Comments:    r.Comments,
		DateOfVisit: r.DateOfVisit,
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
		CreatedBy:   user,
		UpdatedBy:   nil,
	}

	err := h.HomeVisits.Save(visit)
	if err != nil {
		return nil, fmt.Errorf("error creating home visit: %+v", err)
	}

	return &visit, nil
}
