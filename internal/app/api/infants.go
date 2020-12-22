package api

import (
	"encoding/json"
	"fmt"
	"moh.gov.bz/mch/emtct/internal/business/data/patient"
	"moh.gov.bz/mch/emtct/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/app"
	"moh.gov.bz/mch/emtct/internal/business/data/infant"
	"moh.gov.bz/mch/emtct/internal/business/data/labs"
	"moh.gov.bz/mch/emtct/internal/business/data/pregnancy"
)

type InfantRoutes struct {
	Infant      infant.Infants
	Patient     patient.Patients
	Pregnancies pregnancy.Pregnancies
	Labs        labs.Labs
}

func (i InfantRoutes) InfantHandlers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, nil)
	}
}

func (i InfantRoutes) InfantDiagnosesHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, nil)
	}
}

// HIV SCREENINGS
type newHivScreeningRequest struct {
	PatientId              string     `json:"patientId"`
	MotherId               string     `json:"motherId"`
	TestName               string     `json:"testName"`
	ScreeningDate          time.Time  `json:"screeningDate"`
	DateSampleReceivedAtHq *time.Time `json:"dateSampleReceivedAtHq,omitempty"`
	SampleCode             string     `json:"sampleCode"`
	DateSampleShipped      *time.Time `json:"dateSampleShipped"`
	Destination            string     `json:"destination"`
	DateResultReceived     *time.Time `json:"dateResultReceived,omitempty"`
	Result                 string     `json:"result"`
	DateResultShared       *time.Time `json:"dateResultShared,omitempty"`
	DateSampleTaken        time.Time  `json:"dateSampleTaken"`
}

type hivScreeningsResponse struct {
	HivScreenings []infant.HivScreening `json:"hivScreenings"`
	Infant        models.Infant         `json:"patient"`
}

func (i InfantRoutes) CreateHivScreening(user string, r newHivScreeningRequest, timely bool, dueDate time.Time) (*infant.HivScreening, error) {
	id := uuid.New().String()

	s := infant.HivScreening{
		Id:                     id,
		PatientId:              r.PatientId,
		TestName:               r.TestName,
		ScreeningDate:          r.ScreeningDate,
		DateSampleReceivedAtHq: r.DateSampleReceivedAtHq,
		SampleCode:             r.SampleCode,
		DateSampleShipped:      r.DateSampleShipped,
		Destination:            r.Destination,
		DateResultReceived:     r.DateResultReceived,
		Result:                 r.Result,
		DateResultShared:       r.DateResultShared,
		DateSampleTaken:        &r.DateSampleTaken,
		MotherId:               r.MotherId,
		DueDate:                &dueDate,
		CreatedAt:              time.Now(),
		UpdatedAt:              nil,
		CreatedBy:              user,
		UpdatedBy:              nil,
		Timely:                 timely,
	}

	err := i.Infant.SaveHivScreening(s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (i InfantRoutes) HivScreeningHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "HivScreeningHandler"
	w.Header().Add("Content-Type", "application/json")
	var user string
	defer r.Body.Close()
	// Only try to extract the Jwt Token info if the request is not an OPTIONS request.
	// The middleware that verifies the token immediately returns, without inspecting the token.
	// This means that trying to extract the token information in an OPTIONS request will always fail.
	if r.Method != http.MethodOptions {
		token := r.Context().Value("user").(app.JwtToken)
		user = token.Email
	}

	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodPost:
		var req newHivScreeningRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.WithError(err).Error("error parsing request body for creating an hiv screening")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		pId, err := strconv.Atoi(req.PatientId)
		if err != nil {
			log.WithFields(log.Fields{
				"handler": handlerName,
				"method":  r.Method,
				"request": req,
			}).WithError(err).Error("patient id is not a valid number")
			http.Error(w, "patient id is not a valid number", http.StatusBadRequest)
			return
		}
		infant, err := i.Infant.FindInfant(pId)
		if err != nil {
			log.WithFields(log.Fields{
				"user":    user,
				"request": req,
				"handler": "CreateHivScreeningHandler",
			}).WithError(err).Error("")
			http.Error(w, fmt.Sprintf("no birth was found for this infant id: %d", req.PatientId), http.StatusBadRequest)
			return
		}
		timely := i.Infant.IsHivScreeningTimely(*infant.Infant.Dob, req.TestName, req.DateSampleTaken)
		dueDate := i.Infant.HivScreeningDueDate(req.TestName, *infant.Infant.Dob)
		screening, err := i.CreateHivScreening(user, req, timely, dueDate)
		if err != nil {
			log.WithFields(log.Fields{
				"hivScreeningRequest": req,
				"handler":             "CreateHivScreeningHandler",
				"user":                user,
			}).WithError(err).Error("failed to create an hiv screening")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(screening); err != nil {
			log.WithFields(log.Fields{
				"screening": screening,
				"handler":   "CreateHivScreeningHandler",
				"user":      user,
			}).WithError(err).Error("error marshalling result for newly created hiv screening")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPut:
		var screening infant.HivScreening
		if err := json.NewDecoder(r.Body).Decode(&screening); err != nil {
			log.WithFields(log.Fields{
				"request": r.Body,
				"user":    user,
			}).WithError(err).Error("failure parsing the body for editing an hiv screening")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		s, err := i.Infant.FindHivScreeningById(screening.Id)
		if err != nil {
			log.WithFields(log.Fields{
				"screeningRequest": screening,
				"user":             user,
			}).WithError(err).Error("error when querying database for hiv screening")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if s == nil {
			log.WithFields(log.Fields{
				"screeningRequest": screening,
				"user":             user,
			}).Error("tried to update an hiv screening that does not exist")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		infant, err := i.Patient.GetInfant(screening.PatientId)
		if err != nil {
			log.WithFields(log.Fields{
				"user":    user,
				"request": screening,
				"handler": "CreateHivScreeningHandler",
			}).WithError(err).Error("")
			http.Error(w, fmt.Sprintf("no birth was found for infant Id: %s", screening.PatientId), http.StatusBadRequest)
			return
		}
		timely := i.Infant.IsHivScreeningTimely(*infant.Dob, screening.TestName, *screening.DateSampleTaken)
		screening.UpdatedBy = &user
		screening.Timely = timely
		// Never allow the user to modify the due date because this is a computed value
		screening.DueDate = s.DueDate
		screening.CreatedBy = s.CreatedBy
		screening.CreatedAt = s.CreatedAt
		err = i.Infant.SaveHivScreening(screening)
		if err != nil {
			log.WithFields(log.Fields{
				"screeningId": screening.Id,
				"user":        user,
				"request":     screening,
			}).WithError(err).Error("db failure while editing an hiv screening")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(screening); err != nil {
			log.WithFields(log.Fields{
				"screening": screening,
				"user":      user,
			}).WithError(err).Error("failure marshalling the edited hiv screening")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodGet:
		vars := mux.Vars(r)
		patientId := vars["infantId"]
		id, err := strconv.Atoi(patientId)
		if err != nil {
			log.WithError(err).Error("infant id must be a number")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		screenings, err := i.Infant.FindHivScreeningsByPatient(patientId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
				"handler":   "HivScreeningsByPatientIdHandler",
			}).WithError(err).Error("error retrieving hiv screenings for patient")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		infant, err := i.Patient.GetInfant(patientId)
		if err != nil {
			log.WithFields(log.Fields{"patientId": id, "screenings": screenings}).WithError(err).
				Error("error retrieving patient when fetching hiv screenings")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response := hivScreeningsResponse{
			HivScreenings: screenings,
			Infant:        *infant,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{
				"patientId":  patientId,
				"screenings": screenings,
				"handler":    "HivScreeningsByPatientIdHandler",
			}).WithError(err).Error("error marshalling screenings")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}
}

type infantSyphilisScreeningResponse struct {
	Infant     infant.Infant            `json:"infant"`
	Screenings []labs.SyphilisScreening `json:"screenings"`
}

func (i InfantRoutes) InfantSyphilisScreeninngHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		infantId := vars["infantId"]
		id, err := strconv.Atoi(infantId)
		if err != nil {
			log.WithFields(log.Fields{
				"infantId": infantId,
				"handler":  "InfantSyphilisScreeningHandler",
			}).WithError(err).Error("infantId is not a valid number: %+v", err)
			http.Error(w, "infantId is not a valid number", http.StatusBadRequest)
			return
		}
		infantInfo, err := i.Infant.FindInfant(id)
		if err != nil {
			log.WithFields(log.Fields{
				"infantId": id,
				"handler":  "InfantSyphilisScreeningHandler",
			}).WithError(err).Error("error retrieving infant information")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if infantInfo == nil {
			log.WithFields(log.Fields{
				"infantId": id,
				"handler":  "InfantSyphilisScreeningHandler",
			}).Error("no infant exists with given id")
			http.Error(w, fmt.Sprintf("infant with id %d does not exist", id), http.StatusNotFound)
			return
		}
		screenings, err := i.Labs.FindInfantSyphilisScreenings(id, *infantInfo.Infant.Dob)
		if err != nil {
			log.WithFields(log.Fields{
				"infantId":   id,
				"handler":    "InfantSyphilisScreeningHandler",
				"infantInfo": infantInfo,
			}).WithError(err).Error("error retrieving syphilis screenings for infant")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response := infantSyphilisScreeningResponse{
			Infant:     *infantInfo,
			Screenings: screenings,
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{
				"infantId":   id,
				"screenings": screenings,
				"handler":    "InfantSyphilisScreeningHandler",
			}).WithError(err).Error("error marshalling screening data")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
