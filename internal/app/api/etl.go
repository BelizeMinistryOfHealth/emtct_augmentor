package api

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/business/data/pregnancy"
)

type Etl struct {
	Pregnancies pregnancy.Pregnancies
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
