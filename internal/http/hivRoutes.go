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

// HIV SCREENINGS
type NewHivScreeningRequest struct {
	PatientId              int        `json:"patientId"`
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
	MchEncounterId         int        `json:"mchEncounterId"`
}

func (a *App) CreateHivScreening(user string, r NewHivScreeningRequest, timely bool, dueDate time.Time) (*models.HivScreening, error) {
	id := uuid.New().String()

	s := models.HivScreening{
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
		MchEncounterId:         r.MchEncounterId,
		DueDate:                &dueDate,
		CreatedAt:              time.Now(),
		UpdatedAt:              nil,
		CreatedBy:              user,
		UpdatedBy:              nil,
		Timely:                 timely,
	}

	err := a.Db.CreateHivScreening(s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (a *App) findBirthForEncounter(patientId, mchEncounterId int) (*models.Birth, error) {
	anc, err := a.AcsisDb.FindAntenatalEncounterById(patientId, mchEncounterId)
	if err != nil {
		return nil, fmt.Errorf("error fetching antenatal encounter: %+v", err)
	}
	if anc == nil {
		return nil, fmt.Errorf("error: no antenatal encounter was found: %+v", err)
	}
	birth, err := a.AcsisDb.FindLatestBirth(patientId, anc.Id)
	if err != nil {
		return nil, fmt.Errorf("error fetching births related to an antenatal encounter: %+v", err)
	}
	if birth == nil {
		return nil, fmt.Errorf("error: no birth found for antenatal encounter")
	}
	return birth, nil
}

func (a *App) CreateHivScreeningHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodPost:
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		var req NewHivScreeningRequest
		err := parseBody(r.Body, &req)
		if err != nil {
			log.WithError(err).Error("error parsing request body for creating an hiv screening")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// Calculate if the sample was taken in a timely manner. We need the birth date for this.
		birth, err := a.findBirthForEncounter(req.PatientId, req.MchEncounterId)
		if err != nil {
			log.WithFields(log.Fields{
				"user":    user,
				"request": req,
				"handler": "CreateHivScreeningHandler",
			}).WithError(err).Error("")
			http.Error(w, fmt.Sprintf("no birth was found for this mch encounter: %s", req.MchEncounterId), http.StatusBadRequest)
			return
		}
		timely := models.IsHivScreeningTimely(*birth, req.TestName, req.DateSampleTaken)
		dueDate := models.HivScreeningDueDate(req.TestName, birth.BirthDate)
		screening, err := a.CreateHivScreening(user, req, timely, dueDate)
		if err != nil {
			log.WithFields(log.Fields{
				"hivScreeningRequest": req,
				"handler":             "CreateHivScreeningHandler",
				"user":                user,
			}).WithError(err).Error("failed to create an hiv screening")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		result, err := json.Marshal(screening)
		if err != nil {
			log.WithFields(log.Fields{
				"screening": screening,
				"handler":   "CreateHivScreeningHandler",
				"user":      user,
			}).WithError(err).Error("error marshalling result for newly created hiv screening")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, string(result))

	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

type hivScreeningsResponse struct {
	HivScreenings []models.HivScreening   `json:"hivScreenings"`
	Patient       models.PatientBasicInfo `json:"patient"`
}

func (a *App) HivScreeningsByPatientIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	patientId := vars["patientId"]
	id, err := strconv.Atoi(patientId)
	if err != nil {
		log.WithError(err).Error("patient id must be a number")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	screenings, err := a.Db.FindHivScreeningsByPatient(id)
	if err != nil {
		log.WithFields(log.Fields{
			"patientId": patientId,
			"handler":   "HivScreeningsByPatientIdHandler",
		}).WithError(err).Error("error retrieving hiv screenings for patient")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	patient, err := a.AcsisDb.FindPatientBasicInfo(id)
	if err != nil {
		log.WithFields(log.Fields{"patientId": id, "screenings": screenings}).WithError(err).
			Error("error retrieving patient when fetching hiv screenings")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := hivScreeningsResponse{
		HivScreenings: screenings,
		Patient:       *patient,
	}

	results, err := json.Marshal(response)
	if err != nil {
		log.WithFields(log.Fields{
			"patientId":  patientId,
			"screenings": screenings,
			"handler":    "HivScreeningsByPatientIdHandler",
		}).WithError(err).Error("error marshalling screenings")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(results))
}

func (a *App) HivScreeningApi(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodPut:
		vars := mux.Vars(r)
		screeningId := vars["screeningId"]
		token := r.Context().Value("user").(JwtToken)
		user := token.Email
		s, err := a.Db.FindHivScreeningById(screeningId)
		if err != nil {
			log.WithFields(log.Fields{
				"screeningId": screeningId,
				"user":        user,
			}).WithError(err).Error("error when querying database for hiv screening")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if s == nil {
			log.WithFields(log.Fields{
				"screeningId": screeningId,
				"user":        user,
			}).Error("tried to update an hiv screening that does not exist")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		var req models.HivScreening
		err = parseBody(r.Body, &req)
		if err != nil {
			log.WithFields(log.Fields{
				"screeningId": screeningId,
				"user":        user,
			}).WithError(err).Error("failure parsing the body for editing an hiv screening")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if req.Id != screeningId {
			log.WithFields(log.Fields{
				"screeningId": screeningId,
				"request":     req,
				"user":        user,
			}).Error("the screening id in the body does not match the resource screening id")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		birth, err := a.findBirthForEncounter(req.PatientId, req.MchEncounterId)
		if err != nil {
			log.WithFields(log.Fields{
				"user":    user,
				"request": req,
				"handler": "CreateHivScreeningHandler",
			}).WithError(err).Error("")
			http.Error(w, fmt.Sprintf("no birth was found for this mch encounter: %s", req.MchEncounterId), http.StatusBadRequest)
			return
		}
		timely := models.IsHivScreeningTimely(*birth, req.TestName, *req.DateSampleTaken)
		req.UpdatedBy = &user
		req.Timely = timely
		saved, err := a.Db.EditHivScreening(req)
		if err != nil {
			log.WithFields(log.Fields{
				"screeningId": screeningId,
				"user":        user,
				"request":     req,
			}).WithError(err).Error("db failure while editing an hiv screening")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		result, err := json.Marshal(saved)
		if err != nil {
			log.WithFields(log.Fields{
				"screeningId":   screeningId,
				"user":          user,
				"editedRequest": saved,
			}).WithError(err).Error("failure marshalling the edited hiv screening")
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
