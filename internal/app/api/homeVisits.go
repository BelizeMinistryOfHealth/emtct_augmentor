package api

import (
	"encoding/json"
	"fmt"
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
	Patient    patient.BasicInfo      `json:"patient"`
}

func (h HomeVisitRoutes) FindByPatientHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		id := vars["id"]
		patientId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
			}).WithError(err).Error("error trying to retrieve home visits using a non-numeric patient id")
			http.Error(w, "the patient id must be a valid number", http.StatusBadRequest)
			return
		}
		homeVisits, err := h.HomeVisits.FindByPatientId(patientId)
		if err != nil {
			log.WithFields(log.Fields{"patientId": patientId}).
				WithError(err).
				Error("database error while retrieving home visits")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		patient, err := h.Patients.FindBasicInfo(patientId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": patient,
				"handler":   "FindHomeVisitsByPatient",
			}).WithError(err).Error("error retrieving patient information")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response := homeVisitResponse{
			HomeVisits: homeVisits,
			Patient:    *patient,
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{"patientId": patientId, "homeVisits": homeVisits}).
				WithError(err).
				Error("error marshalling the home visits results")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
	defer r.Body.Close()
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
		w.Header().Add("Content-Type", "application/json")
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
				"handler": "HomeVisitsHandler",
				"Method":  "PUT",
			}).WithError(err).Error("failed to decode the home visit request")
			http.Error(w, "could not decode your request", http.StatusInternalServerError)
			return
		}
		visit, err := h.editHomeVisit(user, req)
		if err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"homeVisit": req,
			}).WithError(err).Error("error editing home visit")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(visit); err != nil {
			log.WithFields(log.Fields{
				"user":      user,
				"homeVisit": visit,
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
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.WithFields(log.Fields{
				"user": user,
			}).WithError(err).Error("failed to parse body for creating a home visit")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		visit, err := h.createHomeVisit(user, req)
		if err != nil {
			log.WithFields(log.Fields{
				"request": req,
			}).WithError(err).Error("error creating home visit")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(visit); err != nil {
			log.WithFields(log.Fields{
				"homeVisit": visit,
			}).WithError(err).Error("failed to marshal the newly created home visit")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}
}

func (h HomeVisitRoutes) editHomeVisit(user string, r homeVisitRequest) (*homeVisits.HomeVisit, error) {

	v, err := h.HomeVisits.FindById(r.ID)
	if err != nil {
		return nil, fmt.Errorf("home visit with given id does not exist: %+v", err)
	}
	modified, err := h.HomeVisits.Edit(homeVisits.HomeVisit{
		Id:             v.Id,
		PatientId:      v.PatientId,
		MchEncounterId: v.MchEncounterId,
		Reason:         r.Reason,
		Comments:       r.Comments,
		DateOfVisit:    r.DateOfVisit,
		CreatedAt:      v.CreatedAt,
		UpdatedAt:      v.UpdatedAt,
		CreatedBy:      v.CreatedBy,
		UpdatedBy:      &user,
	})
	return modified, err
}

func (h HomeVisitRoutes) createHomeVisit(user string, r newHomeVisitRequest) (*homeVisits.HomeVisit, error) {
	id := uuid.New().String()
	if len(user) == 0 {
		return nil, fmt.Errorf("user did not provide an email")
	}

	visit := homeVisits.HomeVisit{
		Id:             id,
		PatientId:      r.PatientId,
		MchEncounterId: r.MchEncounterId,
		Reason:         r.Reason,
		Comments:       r.Comments,
		DateOfVisit:    r.DateOfVisit,
		CreatedAt:      time.Now(),
		UpdatedAt:      nil,
		CreatedBy:      user,
		UpdatedBy:      nil,
	}

	err := h.HomeVisits.Create(visit)
	if err != nil {
		return nil, fmt.Errorf("error creating home visit: %+v", err)
	}

	return &visit, nil
}
