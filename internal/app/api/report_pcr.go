package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"moh.gov.bz/mch/emtct/internal/business/data/reports"
	"net/http"
	"strconv"
)

type ReportsRoutes struct {
	Report reports.Reports
}

// MissingPcrReport is the HTTP Handler for retrieving all HIV Screenings that have a missing PCR Test.
func (i ReportsRoutes) MissingPcrReport(w http.ResponseWriter, r *http.Request) {
	handlerName := "MissingPcrReport"
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	year, _ := strconv.Atoi(vars["year"])
	reports, err := i.Report.MissingPcrs(r.Context(), year)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": handlerName,
			"method":  r.Method,
			"year":    year,
		}).WithError(err).Error("failed to get missing pcrs")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(reports); err != nil {
		log.WithFields(log.Fields{
			"handler": handlerName,
			"method":  r.Method,
			"year":    year,
			"reports": reports,
		}).WithError(err).Error("could not encode the report")
		return
	}
}
