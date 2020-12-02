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
	"moh.gov.bz/mch/emtct/internal/business/data/infant"
	"moh.gov.bz/mch/emtct/internal/business/data/pregnancy"
	"moh.gov.bz/mch/emtct/internal/business/data/prescription"
)

type InfantRoutes struct {
	Infant      infant.Infants
	Pregnancies pregnancy.Pregnancies
}

func (i InfantRoutes) InfantHandlers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		motherId := vars["patientId"]
		id, err := strconv.Atoi(motherId)
		if err != nil {
			log.WithFields(log.Fields{
				"motherId": motherId,
			}).WithError(err).Error("motherId is not a valid number")
			http.Error(w, "motherId is not a valid number", http.StatusBadRequest)
			return
		}
		// Find current pregnancy
		preg, err := i.Pregnancies.FindLatest(id)
		if err != nil {
			log.WithFields(log.Fields{
				"motherId": id,
			}).WithError(err).Error("error retrieving this patient's latest pregnancy")
			http.Error(w, "could not retrieve the mother's latest pregnancy", http.StatusInternalServerError)
			return
		}
		infant, err := i.Infant.FindPregnancyInfant(*preg)
		if err != nil {
			log.WithFields(log.Fields{
				"motherId": motherId,
			}).WithError(err).Error("error retrieving pregnancy infant")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if infant == nil {
			log.WithFields(log.Fields{
				"motherId": motherId,
			}).Error("no infant exists for current pregnancy")
			http.Error(w, "no infant exists for relevant pregnancy", http.StatusNotFound)
			return
		}
		result, err := json.Marshal(infant)
		if err != nil {
			log.WithFields(log.Fields{
				"infant": infant,
			}).WithError(err).Error("error marshalling infant data")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(result))
	}
}

type infantDiagnosesResponse struct {
	Diagnoses []infant.Diagnoses `json:"diagnoses"`
	Infant    infant.Infant      `json:"infant"`
}

func (i InfantRoutes) InfantDiagnosesHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		id := vars["infantId"]
		infantId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{"infantId": id}).WithError(err).
				Error("the infant id is not a valid number")
			http.Error(w, "infant id must be a numeric value", http.StatusBadRequest)
			return
		}
		mId := vars["motherId"]
		motherId, err := strconv.Atoi(mId)
		if err != nil {
			log.WithFields(log.Fields{"motherId": id}).WithError(err).
				Error("the mother id is not a valid number")
			http.Error(w, "mother id must be a numeric value", http.StatusBadRequest)
			return
		}
		token := r.Context().Value("user").(app.JwtToken)
		user := token.Email
		diagnoses, err := i.Infant.FindInfantDiagnoses(infantId)
		if err != nil {
			log.WithFields(log.Fields{
				"infantId": infantId,
				"user":     user,
				"handler":  "InfantDiagnosesHandler",
			}).
				WithError(err).
				Error("error while fetching infant diagnoses")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// Return an empty array if no results are found
		if diagnoses == nil {
			diagnoses = []infant.Diagnoses{}
		}

		// Find current pregnancy
		preg, err := i.Pregnancies.FindLatest(motherId)
		if err != nil {
			log.WithFields(log.Fields{
				"motherId": id,
			}).WithError(err).Error("error retrieving this patient's latest pregnancy")
			http.Error(w, "could not retrieve the mother's latest pregnancy", http.StatusInternalServerError)
			return
		}
		infantInfo, err := i.Infant.FindPregnancyInfant(*preg)
		if err != nil {
			log.WithFields(log.Fields{
				"infantId": infantId,
				"motherId": motherId,
				"user":     user,
			}).WithError(err).Error("could not find infant info")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		response := infantDiagnosesResponse{
			Diagnoses: diagnoses,
			Infant:    *infantInfo,
		}
		result, err := json.Marshal(response)
		if err != nil {
			log.WithFields(log.Fields{
				"motherId": id,
				"user":     user,
				"response": response,
			}).
				WithError(err).
				Error("error while marshalling the infant diagnoses")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(result))
	}
}

// HIV SCREENINGS
type newHivScreeningRequest struct {
	PatientId              int        `json:"patientId"`
	MotherId               int        `json:"motherId"`
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
	HivScreenings []infant.HivScreening `json:"hivScreening"`
	Infant        infant.Infant         `json:"patient"`
}

type infantTreatmentResponse struct {
	Prescriptions []prescription.Prescription `json:"prescription"`
	Infant        infant.Infant               `json:"infant"`
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

	err := i.Infant.CreateHivScreening(s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (i InfantRoutes) HivScreeningHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	token := r.Context().Value("user").(app.JwtToken)
	user := token.Email

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
		infant, err := i.Infant.FindInfant(req.PatientId)
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
		w.Header().Add("Content-Type", "application/json")

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
		infant, err := i.Infant.FindInfant(screening.PatientId)
		if err != nil {
			log.WithFields(log.Fields{
				"user":    user,
				"request": screening,
				"handler": "CreateHivScreeningHandler",
			}).WithError(err).Error("")
			http.Error(w, fmt.Sprintf("no birth was found for infant Id: %s", screening.PatientId), http.StatusBadRequest)
			return
		}
		timely := i.Infant.IsHivScreeningTimely(*infant.Infant.Dob, screening.TestName, *screening.DateSampleTaken)
		screening.UpdatedBy = &user
		screening.Timely = timely
		saved, err := i.Infant.EditHivScreening(screening)
		if err != nil {
			log.WithFields(log.Fields{
				"screeningId": screening.Id,
				"user":        user,
				"request":     screening,
			}).WithError(err).Error("db failure while editing an hiv screening")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(saved); err != nil {
			log.WithFields(log.Fields{
				"screening":     saved,
				"user":          user,
				"editedRequest": saved,
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
		screenings, err := i.Infant.FindHivScreeningsByPatient(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
				"handler":   "HivScreeningsByPatientIdHandler",
			}).WithError(err).Error("error retrieving hiv screenings for patient")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		infant, err := i.Infant.FindInfant(id)
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

		w.Header().Add("Content-Type", "application/json")

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

func (i InfantRoutes) InfantSyphilisTreatmentHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		id := vars["infantId"]
		token := r.Context().Value("user").(app.JwtToken)
		user := token.Email
		infantId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{
				"infantId": id,
				"user":     user,
			}).WithError(err).Error("patient id is not a valid number")
			http.Error(w, "patient id is not a valid number", http.StatusBadRequest)
			return
		}
		treatments, err := i.Infant.FindInfantSyphilisTreatment(infantId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": infantId,
				"user":      user,
			}).WithError(err).Error("error retrieving syphilis treatment")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		infant, err := i.Infant.FindInfant(infantId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": infantId,
				"user":      user,
				"handler":   "SyphilisTreatmentHandler",
			}).WithError(err).Error("error retrieving patient information")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response := infantTreatmentResponse{
			Prescriptions: treatments,
			Infant:        *infant,
		}
		w.Header().Add("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{
				"response":  response,
				"user":      user,
				"patientId": infantId,
				"handler":   "SyphilisTreatmentHandler",
			}).WithError(err).Error("error marshaling syphilis treatment response")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
