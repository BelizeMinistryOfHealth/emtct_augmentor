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

type HomeVisitResponse struct {
	HomeVisits []models.HomeVisit      `json:"homeVisits"`
	Patient    models.PatientBasicInfo `json:"patient"`
}

func (a *App) FindHomeVisitsByPatient(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

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
	homeVisits, err := a.Db.FindHomeVisitsByPatientId(patientId)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId}).
			WithError(err).
			Error("database error while retrieving home visits")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	patient, err := a.AcsisDb.FindPatientBasicInfo(patientId)
	if err != nil {
		log.WithFields(log.Fields{
			"patientId": patient,
			"handler":   "FindHomeVisitsByPatient",
		}).WithError(err).Error("error retrieving patient information")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := HomeVisitResponse{
		HomeVisits: homeVisits,
		Patient:    *patient,
	}
	results, err := json.Marshal(response)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId, "homeVisits": homeVisits}).
			WithError(err).
			Error("error marshalling the home visits results")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(results))
}

func (a *App) HomeVisitApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["homeVisitId"]

	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		homeVisit, err := a.Db.FindHomeVisitById(id)
		if err != nil {
			log.WithFields(log.Fields{"homeVisitId": id}).
				WithError(err).
				Error("error retrieving home visit from the database")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(homeVisit)
		if err != nil {
			log.WithFields(log.Fields{"homeVisit": homeVisit, "homeVisitId": id}).
				WithError(err).
				Error("failed to marshal home visit")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, string(result))
	case http.MethodPut:
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		var req HomeVisitRequest
		err := parseBody(r.Body, &req)
		if err != nil {
			log.WithFields(log.Fields{
				"homeVisitId": id,
				"user":        user,
			}).WithError(err).Error("parsing body failed")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		homeVisit, err := a.EditHomeVisit(id, user, req)
		if err != nil {
			log.WithFields(log.Fields{
				"user":        user,
				"homeVisitId": id,
			}).WithError(err).Error("error editing home visit")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// Return the home visit
		result, err := json.Marshal(homeVisit)
		if err != nil {
			log.WithFields(log.Fields{
				"user":        user,
				"homeVisitId": id,
				"homeVisit":   homeVisit,
			}).
				WithError(err).
				Error("error marshalling home visit into json")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, string(result))
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

type HomeVisitRequest struct {
	Reason      string    `json:"reason"`
	Comments    string    `json:"comments"`
	DateOfVisit time.Time `json:"dateOfVisit"`
}

func (a *App) EditHomeVisit(id, user string, r HomeVisitRequest) (*models.HomeVisit, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("the home visit id can not be empty")
	}

	if len(user) == 0 {
		return nil, fmt.Errorf("must provide the email of the user making the update")
	}

	h, err := a.Db.FindHomeVisitById(id)
	if err != nil {
		return nil, fmt.Errorf("home visit with given id does not exist: %+v", err)
	}
	homeVisit, err := a.Db.EditHomeVisit(models.HomeVisit{
		Id:          h.Id,
		PatientId:   h.PatientId,
		Reason:      r.Reason,
		Comments:    r.Comments,
		DateOfVisit: r.DateOfVisit,
		CreatedAt:   h.CreatedAt,
		UpdatedAt:   h.UpdatedAt,
		CreatedBy:   h.CreatedBy,
		UpdatedBy:   &user,
	})
	return homeVisit, err
}

type NewHomeVisitRequest struct {
	PatientId      int       `json:"patientId"`
	Reason         string    `json:"reason"`
	Comments       string    `json:"comments"`
	DateOfVisit    time.Time `json:"dateOfVisit"`
	MchEncounterId int       `json:"mchEncounterId"`
}

func (a *App) PostHomeVisit(w http.ResponseWriter, r *http.Request) {

	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodPost:
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		var req NewHomeVisitRequest
		err := parseBody(r.Body, &req)
		if err != nil {
			log.WithFields(log.Fields{
				"user": user,
			}).WithError(err).Error("failed to parse body for creating a home visit")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		homeVisit, err := a.CreateHomeVisit(user, req)
		if err != nil {
			log.WithFields(log.Fields{
				"request": req,
			}).WithError(err).Error("error creating home visit")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		//Marshal and return the newly created home visit
		result, err := json.Marshal(homeVisit)
		if err != nil {
			log.WithFields(log.Fields{
				"homeVisit": homeVisit,
			}).WithError(err).Error("failed to marshal the newly created home visit")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, string(result))

	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

}

func (a *App) CreateHomeVisit(user string, r NewHomeVisitRequest) (*models.HomeVisit, error) {
	id := uuid.New().String()

	if len(user) == 0 {
		return nil, fmt.Errorf("user did not provide an email")
	}

	homeVisit := models.HomeVisit{
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

	err := a.Db.CreateHomeVisit(homeVisit)
	if err != nil {
		return nil, fmt.Errorf("error creating home visit: %+v", err)
	}

	return &homeVisit, nil
}
