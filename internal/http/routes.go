package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
		user := r.Context().Value("user").(JwtToken)
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
		homeVisit, err := a.EditHomeVisit(id, user.Email, req)
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
