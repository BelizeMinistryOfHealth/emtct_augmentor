package api

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/business/data/patient"
	"moh.gov.bz/mch/emtct/internal/business/data/pregnancy"
)

type Etl struct {
	Pregnancies pregnancy.Pregnancies
	patients    patient.Patients
}

type pregnancyEtlRequest struct {
	Year int `json:"year"`
}

func (e Etl) PregnancyEtlHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	handlerName := "PregnancyEtlHandler"
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodPost:
		method := "POST"

		var req pregnancyEtlRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.WithFields(log.Fields{
				"method":  method,
				"handler": handlerName,
			}).WithError(err).Error("could not decode request")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		yr := req.Year
		existingPregnancies, err := e.Pregnancies.FindExistingPregnanciesByYear(yr)
		if err != nil {
			log.WithFields(log.Fields{
				"year":    yr,
				"method":  method,
				"handler": handlerName,
			}).WithError(err).Error("error while fetching existing pregnancies")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		acsisPregnancies, err := e.Pregnancies.FindPregnanciesInBhisByYear(yr)
		if err != nil {
			log.WithFields(log.Fields{
				"year":    yr,
				"method":  method,
				"handler": handlerName,
			}).WithError(err).Error("error retrieving pregnancies from acsis")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		var pregs []pregnancy.Pregnancy
		for _, p := range acsisPregnancies {
			if !p.Include(existingPregnancies) {
				pregs = append(pregs, p)
			}
		}
		err = e.Pregnancies.Create(r.Context(), pregs)
		if err != nil {
			log.WithFields(log.Fields{
				"year":        yr,
				"pregnancies": pregs,
				"method":      method,
				"handler":     handlerName,
			}).WithError(err).Error("error inserting pregnancy into emtct db")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(pregs); err != nil {
			log.WithFields(log.Fields{
				"year":        yr,
				"pregnancies": pregs,
				"method":      method,
				"handler":     handlerName,
			}).WithError(err).Error("failed to encode pregnancies")
			http.Error(w, "inserted pregnancies but failed to encode them", http.StatusInternalServerError)
			return
		}
	}
}

// PatientEtlHandler retrieves patients from acsis, and inserts them into the emtct database.
func (e Etl) PatientEtlHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "PatientEtlHandler"
	defer r.Body.Close()
	type patientRequest struct {
		Year int `json:"year"`
	}
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodPost:
		method := "POST"
		var req patientRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.WithFields(log.Fields{
				"handler": handlerName,
				"method":  method,
			}).WithError(err).Error("failed to decode the request")
			http.Error(w, "could not decode the request", http.StatusBadRequest)
			return
		}
		// BHIS implementation began in 2008, so we should not allow queries for years before that.
		if req.Year < 2008 {
			http.Error(w, "only years greater than 2007 are valid", http.StatusBadRequest)
			return
		}
		acsisPatients, err := e.patients.FindInBhisByYear(req.Year)
		if err != nil {
			log.WithFields(log.Fields{
				"handler": handlerName,
				"method":  method,
				"request": req,
			}).WithError(err).Error("patient query by year failed")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Infof("inserting %d patients into emtct database", len(acsisPatients))
		if err := e.patients.Create(r.Context(), acsisPatients); err != nil {
			log.WithFields(log.Fields{
				"handler":  handlerName,
				"method":   method,
				"patients": len(acsisPatients),
			}).WithError(err).Error("error inserting acsis patients into emtct database")
			http.Error(w, "inserting into emtct database failed", http.StatusInternalServerError)
			return
		}
		log.Infof("created %d patients into emtct database", len(acsisPatients))
		response := map[string]interface{}{"total": len(acsisPatients)}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.WithFields(log.Fields{
				"handler":  handlerName,
				"method":   method,
				"patients": len(acsisPatients),
			}).WithError(err).Error("error decoding response")
			http.Error(w, "could not decode the response", http.StatusInternalServerError)
			return
		}
	}
}
