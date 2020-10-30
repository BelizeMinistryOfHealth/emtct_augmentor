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

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

// TestAuth tests that authentication is working
func TestAuth(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	log.Printf("user: %+v", user)
	fmt.Fprintf(w, "TEST")
}

type PatientResponse struct {
	Patient          *models.Patient           `json:"patient"`
	ObstetricHistory []models.ObstetricHistory `json:"obstetricHistory"`
	Diagnoses        []models.Diagnosis        `json:"diagnoses"`
}

func (a *App) RetrievePatient(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	patientId := vars["id"]
	patient, err := a.Db.FindPatientById(patientId)
	if err != nil {
		log.WithFields(
			log.Fields{"request": r}).WithError(err).Error("could not find patient with specified id")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	diagnoses, err := a.Db.FindDiagnoses(patientId)
	if err != nil {
		log.WithFields(
			log.Fields{"request": r}).WithError(err).Error("could not retrieve obstetric history for the patient")
	}

	obstetricHistory, err := a.Db.FindObstetricHistory(patientId)
	if err != nil {
		log.WithFields(log.Fields{"request": r}).WithError(err).Error("could not retrieve obstetric history")
	}

	if patient == nil {
		w.Header().Add("Content-Type", "application/json")
		emptyResponse := PatientResponse{
			Patient:          nil,
			ObstetricHistory: nil,
			Diagnoses:        nil,
		}
		resp, _ := json.Marshal(emptyResponse)
		fmt.Fprint(w, string(resp))
		return
	}
	response := PatientResponse{
		Patient:          patient,
		ObstetricHistory: obstetricHistory,
		Diagnoses:        diagnoses,
	}
	resp, err := json.Marshal(response)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId, "response": response}).WithError(err).Error("error marshalling response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}

type PregnancyResponse struct {
	Vitals    *models.PregnancyVitals `json:"vitals"`
	Diagnoses []models.Diagnosis      `json:"diagnoses"`
}

func (a *App) FindCurrentPregnancy(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}
	vars := mux.Vars(r)
	patientId := vars["id"]
	pregnancy, err := a.Db.FindCurrentPregnancy(patientId)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId}).
			WithError(err).
			Error("error retrieving current pregnancy from database")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	diagnoses, err := a.Db.FindPregnancyDiagnoses(patientId)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId, "pregnancy": pregnancy}).
			WithError(err).
			Error("error fetching diagnoses for a pregnancy")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if pregnancy == nil {
		w.Header().Add("Content-Type", "application/json")
		r, _ := json.Marshal(nil)
		fmt.Fprintf(w, string(r))
	}

	response := PregnancyResponse{
		Vitals:    pregnancy,
		Diagnoses: diagnoses,
	}

	resp, err := json.Marshal(response)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId, "pregnancy": pregnancy, "diagnoses": diagnoses}).
			WithError(err).
			Error("error marshalling pregnancy")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}

func (a *App) FindPregnancyLabResults(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	patientId := vars["id"]
	labResults, err := a.Db.FindPregnancyLabResults(patientId)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId}).
			WithError(err).
			Error("error while retrieving lab results")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	results, err := json.Marshal(labResults)
	if err != nil {
		log.WithFields(log.Fields{"labResults": labResults, "patientId": patientId}).
			WithError(err).
			Error("error marshalling lab results for a pregnancy")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(results))
}

func (a *App) FindHomeVisitsByPatient(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	patientId := vars["id"]
	homeVisits, err := a.Db.FindHomeVisitsByPatientId(patientId)
	if err != nil {
		log.WithFields(log.Fields{"patientId": patientId}).
			WithError(err).
			Error("database error while retrieving home visits")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	results, err := json.Marshal(homeVisits)
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
	PatientId   int       `json:"patientId"`
	Reason      string    `json:"reason"`
	Comments    string    `json:"comments"`
	DateOfVisit time.Time `json:"dateOfVisit"`
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
		Id:          id,
		PatientId:   r.PatientId,
		Reason:      r.Reason,
		Comments:    r.Comments,
		DateOfVisit: r.DateOfVisit,
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
		CreatedBy:   user,
		UpdatedBy:   nil,
	}

	err := a.Db.CreateHomeVisit(homeVisit)
	if err != nil {
		return nil, fmt.Errorf("error creating home visit: %+v", err)
	}

	return &homeVisit, nil
}

// HIV SCREENINGS
type NewHivScreeningRequest struct {
	PatientId              int        `json:"patientId"`
	TestName               string     `json:"testName"`
	ScreeningDate          time.Time  `json:"screeningDate"`
	DateSampleReceivedAtHq *time.Time `json:"dateSampleReceivedAtHq,omitempty"`
	SampleCode             string     `json:"sampleCode"`
	DateSampleShipped      time.Time  `json:"dateSampleShipped"`
	Destination            string     `json:"destination"`
	DateResultReceived     *time.Time `json:"dateResultReceived,omitempty"`
	Result                 string     `json:"result"`
	DateResultShared       *time.Time `json:"dateResultShared,omitempty"`
}

func (a *App) CreateHivScreening(user string, r NewHivScreeningRequest) (*models.HivScreening, error) {
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
		CreatedAt:              time.Now(),
		UpdatedAt:              nil,
		CreatedBy:              user,
		UpdatedBy:              nil,
	}

	err := a.Db.CreateHivScreening(s)
	if err != nil {
		return nil, err
	}

	return &s, nil
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
		screening, err := a.CreateHivScreening(user, req)
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

	results, err := json.Marshal(screenings)
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
		req.UpdatedBy = &user
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

// CONTRACEPTIVES USED

type NewContraceptiveUsedRequest struct {
	PatientId     int       `json:"patientId"`
	Contraceptive string    `json:"contraceptive"`
	Comments      string    `json:"comments"`
	DateUsed      time.Time `json:"dateUsed"`
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
			Id:            uuid.New().String(),
			PatientId:     req.PatientId,
			Contraceptive: req.Contraceptive,
			Comments:      req.Comments,
			DateUsed:      req.DateUsed,
			CreatedAt:     time.Now(),
			UpdatedAt:     nil,
			CreatedBy:     user,
			UpdatedBy:     nil,
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
	results, err := json.Marshal(contraceptives)
	if err != nil {
		log.WithFields(log.Fields{
			"patientId":      patientId,
			"contraceptives": contraceptives,
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
