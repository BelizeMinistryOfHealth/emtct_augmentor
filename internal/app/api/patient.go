package api

import (
	"encoding/json"
	"moh.gov.bz/mch/emtct/internal/auth"
	"moh.gov.bz/mch/emtct/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type patientResponse struct {
	Patient     *models.Patient    `json:"patient"`
	Pregnancies []models.Pregnancy `json:"pregnancies"`
}

func (a *pregnancyRoutes) RetrievePatientHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		patientId := vars["id"]
		patient, err := a.Patient.FindByPatientId(patientId)
		if err != nil {
			log.WithFields(
				log.Fields{"request": r}).WithError(err).Error("could not find patient with specified id")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		v, err := a.Patient.GetPregnancies(patientId)
		if patient == nil {
			emptyResponse := patientResponse{
				Patient: nil,
			}
			_ = json.NewEncoder(w).Encode(emptyResponse)
			return
		}
		response := patientResponse{
			Patient:     patient,
			Pregnancies: v,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{"patientId": patientId, "response": response}).WithError(err).Error("error marshalling response")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

type arvsResponse struct {
	Arvs    []models.Prescription `json:"arvs"`
	Patient models.Patient        `json:"patient"`
}

func (a *pregnancyRoutes) ArvsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	handlerName := "ArvsHandler"
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		method := r.Method
		token := r.Context().Value("user").(auth.JwtToken)
		user := token.Email
		vars := mux.Vars(r)
		id := vars["patientId"]
		patientId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).
				WithError(err).
				Error("patient id is not a valid number")
			http.Error(w, "patientId is not a valid number", http.StatusBadRequest)
			return
		}
		// Find the pregnancy and the lmp so we can get the date bounds
		pregId := vars["pregnancyId"]
		pregnancyId, err := strconv.Atoi(pregId)
		if err != nil {
			log.WithFields(log.Fields{
				"handler":     handlerName,
				"method":      method,
				"user":        user,
				"pregnancyId": pregId,
				"patientId":   patientId,
			}).WithError(err).Error("pregnancy id is not a valid number")
			http.Error(w, "pregnancy id is not a valid number", http.StatusBadRequest)
			return
		}
		pregs, err := a.Patient.GetPregnancy(pregnancyId)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).
				WithError(err).
				Error("error retrieving patient's current pregnancy")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if pregs == nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).
				WithError(err).
				Error("patient does not have a current pregnancy")
			http.Error(w, "patient does not have a current pregnancy", http.StatusNotFound)
			return
		}
		lmp := pregs.Lmp
		nextDate := lmp.Add(time.Hour * 24 * 7 * 54)
		arvs, err := a.Patient.FindArvsByPatient(id, *lmp, nextDate)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).
				WithError(err).
				Error("patient does not have a current pregnancy")
			http.Error(w, "patient does not have a current pregnancy", http.StatusNotFound)
			return
		}
		patientInfo, err := a.Patient.FindByPatientId(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).
				WithError(err).
				Error("error retrieving patient's basic info")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		arvsResponse := arvsResponse{
			Arvs:    arvs,
			Patient: *patientInfo,
		}
		if err := json.NewEncoder(w).Encode(arvsResponse); err != nil {
			log.WithFields(log.Fields{"patientId": patientId, "user": user}).
				WithError(err).
				Error("marshalling arvs response failed")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}
}

type treatmentResponse struct {
	Prescriptions []models.Prescription `json:"prescriptions"`
	Patient       models.Patient        `json:"patient"`
}

func (a *pregnancyRoutes) PatientSyphilisTreatmentHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "PatientSyphilisTreatmentHandler"
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		method := "GET"
		vars := mux.Vars(r)
		id := vars["patientId"]
		token := r.Context().Value("user").(auth.JwtToken)
		user := token.Email
		patientId, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"user":      user,
				"handler":   handlerName,
				"method":    method,
			}).WithError(err).Error("patient id is not a valid number")
			http.Error(w, "patient id is not a valid number", http.StatusBadRequest)
			return
		}
		patientInfo, err := a.Patient.FindByPatientId(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"handler":   handlerName,
				"method":    method,
			}).WithError(err).Error("failed to retrieve patient basic info")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		pregs, err := a.Patient.GetPregnancies(id)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"handler":   handlerName,
				"method":    method,
			}).WithError(err).Error("failed to retrieve patient latest pregnancy")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if pregs == nil {
			log.WithFields(log.Fields{
				"patientId": patientId,
				"user":      user,
			}).Error("patient has no pregnancy")
			http.Error(w, "patient has no pregnancy", http.StatusNoContent)
			return
		}
		preg := pregs[0]
		lmp := preg.Lmp
		endDate := lmp.Add(time.Hour * 24 * 7 * 52)
		treatments, err := a.Patient.FindSyphilisTreatment(id, lmp, &endDate)
		if err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"handler":   handlerName,
				"method":    method,
				"lmp":       lmp,
				"endDate":   endDate,
			}).WithError(err).Error("failed to retrieve patient syphilis treatment")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response := treatmentResponse{
			Prescriptions: treatments,
			Patient:       *patientInfo,
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{
				"patientId": id,
				"handler":   handlerName,
				"method":    method,
				"lmp":       lmp,
				"endDate":   endDate,
				"response":  response,
			}).WithError(err).Error("failed to encode syphilis treatment")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}
}

func (a *pregnancyRoutes) InfantHandlers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Add("Content-Type", "application/json")
	handlerName := "InfantHandlers"
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		motherId := vars["motherId"]
		patientId, err := strconv.Atoi(motherId)
		if err != nil {
			log.WithFields(log.Fields{
				"handler":   handlerName,
				"method":    r.Method,
				"patientId": motherId,
			}).WithError(err).Error("motherId must be a valid number")
			http.Error(w, "mother id must be a valid number", http.StatusBadRequest)
			return
		}
		pregnancyId := vars["pregnancyId"]
		pregId, err := strconv.Atoi(pregnancyId)
		if err != nil {
			log.WithFields(log.Fields{
				"handler":     handlerName,
				"method":      r.Method,
				"motherId":    motherId,
				"pregnancyId": pregnancyId,
			}).WithError(err).Error("pregnancy id must be a valid number")
			http.Error(w, "pregnancy id must be a valid number", http.StatusBadRequest)
			return
		}
		//Get Pregnancy so we can filter infants who are born about 9 months after LMP
		preg, err := a.Patient.GetPregnancy(pregId)
		if err != nil {
			log.WithFields(log.Fields{
				"handler":     handlerName,
				"method":      r.Method,
				"pregnancyId": pregnancyId,
				"motherId":    motherId,
			}).WithError(err).Error("error retrieving pregnancy")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		lmp := preg.Lmp
		if lmp == nil {
			log.WithFields(log.Fields{
				"handler":   handlerName,
				"method":    r.Method,
				"pregnancy": preg,
			}).Error("can not retrieve patient if lmp is missing")
			http.Error(w, "can not retrieve patiet if pregnancy is missing an lmp", http.StatusInternalServerError)
			return
		}
		infant, err := a.Patient.GetInfantForPregnancy(patientId, *lmp)
		if err != nil {
			log.WithFields(log.Fields{
				"handler":   handlerName,
				"method":    r.Method,
				"pregnancy": preg,
			}).WithError(err).Error("failed to retrieve infant for pregnancy")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(infant); err != nil {
			log.WithFields(log.Fields{
				"handler":   handlerName,
				"method":    r.Method,
				"pregnancy": preg,
				"infant":    infant,
			}).WithError(err).Error("failed to json encode infant")
			http.Error(w, "failed to json encode infant", http.StatusInternalServerError)
			return
		}
	}
}

func (a *pregnancyRoutes) InfantByIdHandlers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Add("Content-Type", "application/json")
	handlerName := "InfantByIdHandlers"
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		infantId := vars["infantId"]
		infant, err := a.Patient.GetInfant(infantId)
		if err != nil {
			log.WithFields(log.Fields{
				"handler":  handlerName,
				"method":   r.Method,
				"infantId": infantId,
			}).WithError(err).Error("failed to retrieve infant for pregnancy")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(infant); err != nil {
			log.WithFields(log.Fields{
				"handler": handlerName,
				"method":  r.Method,
				"infant":  infant,
			}).WithError(err).Error("failed to json encode infant")
			http.Error(w, "failed to json encode infant", http.StatusInternalServerError)
			return
		}
	}
}

func (a *pregnancyRoutes) InfantDiagnosesHandler(w http.ResponseWriter, r *http.Request) {
	type infantDiagnosesResponse struct {
		Diagnoses []models.Diagnosis `json:"diagnoses"`
		Infant    *models.Infant     `json:"infant"`
	}
	defer r.Body.Close()
	w.Header().Add("Content-Type", "application/json")
	handlerName := "InfantDiagnosesHandler"
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		infantId := vars["infantId"]
		diagnoses, err := a.Patient.GetDiagnoses(infantId)
		if err != nil {
			log.WithFields(log.Fields{
				"handler":  handlerName,
				"infantId": infantId,
				"method":   r.Method,
			}).WithError(err).Error("failed to fetch infant's diagnoses")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		infant, err := a.Patient.GetInfant(infantId)
		if err != nil {
			log.WithFields(log.Fields{
				"handler":  handlerName,
				"method":   r.Method,
				"infantId": infantId,
			}).WithError(err).Error("failed to retrieve infant")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		resp := infantDiagnosesResponse{
			Diagnoses: diagnoses,
			Infant:    infant,
		}
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			log.WithFields(log.Fields{
				"handler":  handlerName,
				"method":   r.Method,
				"infantId": infantId,
				"response": resp,
			}).WithError(err).Error("failed to encode infant diagnoses response")
			http.Error(w, "failed to encode infant diagnoses", http.StatusInternalServerError)
			return
		}
	}
}

func (a *pregnancyRoutes) InfantSyphilisTreatmentHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Add("Content-Type", "application/json")
	handlerName := "InfantSyphilisTreatmentHandler"

	type infantTreatmentResponse struct {
		Prescriptions []models.Prescription `json:"prescriptions"`
		Infant        models.Infant         `json:"infant"`
	}

	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		infantId := vars["infantId"]
		infant, err := a.Patient.GetInfant(infantId)
		if err != nil {
			log.WithFields(log.Fields{
				"handler":  handlerName,
				"method":   r.Method,
				"infantId": infantId,
			}).WithError(err).Error("")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		endDate := infant.Dob.Add(time.Hour * 24 * 7 * 54)
		prescriptions, err := a.Patient.FindSyphilisTreatment(infantId, infant.Dob, &endDate)
		if err != nil {
			log.WithFields(log.Fields{
				"handler":  handlerName,
				"method":   r.Method,
				"infantId": infantId,
			}).WithError(err).Error("failed to retrieve syphilis treatments")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		resp := infantTreatmentResponse{
			Prescriptions: prescriptions,
			Infant:        *infant,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.WithFields(log.Fields{
				"handler":  handlerName,
				"method":   r.Method,
				"infantId": infantId,
				"response": resp,
			}).WithError(err).Error("failed to encode response")
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func (a *pregnancyRoutes) InfantSyphilisScreeningHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Add("Content-Type", "application/json")
	handlerName := "InfantSyphilisScreeningHandler"

	type infantSyphilisScreeningResponse struct {
		Infant     models.Infant      `json:"infant"`
		Screenings []models.LabResult `json:"screenings"`
	}

	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		vars := mux.Vars(r)
		infantId := vars["infantId"]
		infant, err := a.Patient.GetInfant(infantId)
		if err != nil {
			log.WithFields(log.Fields{
				"handler":  handlerName,
				"method":   r.Method,
				"infantId": infantId,
			}).WithError(err).Error("")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		endDate := infant.Dob.Add(time.Hour * 24 * 7 * 54)
		tests, err := a.Patient.FindSyphilisLabTests(infantId, *infant.Dob, endDate)
		if err != nil {
			log.WithFields(log.Fields{
				"handler":  handlerName,
				"method":   r.Method,
				"infantId": infantId,
			}).WithError(err).Error("failed to retrieve syphilis lab tests")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		resp := infantSyphilisScreeningResponse{
			Infant:     *infant,
			Screenings: tests,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.WithFields(log.Fields{
				"handler":  handlerName,
				"method":   r.Method,
				"response": resp,
			}).WithError(err).Error("failed to encode response")
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
