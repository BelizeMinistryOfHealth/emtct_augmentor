// Package api contains the full set of handler functions and routes
// supported by the web api.
package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/uris77/auth0"

	"moh.gov.bz/mch/emtct/internal/app"
	"moh.gov.bz/mch/emtct/internal/business/data/admissions"
	"moh.gov.bz/mch/emtct/internal/business/data/homeVisits"
	"moh.gov.bz/mch/emtct/internal/business/data/infant"
	"moh.gov.bz/mch/emtct/internal/business/data/patient"
	"moh.gov.bz/mch/emtct/internal/business/data/pregnancy"
)

func API(app app.App) *mux.Router {
	r := mux.NewRouter()

	// Instantiate an aut0 client with a Cache with a key capacity of
	// 60 tokens and a ttl of 24 hours.
	auth0Client := auth0.NewAuth0(60, 518400)

	// Middleware that verifies JWT token and also enables CORS.
	authMid := NewChain(EnableCors(), VerifyToken(app.Auth.JwkUrl, app.Auth.Aud, app.Auth.Iss, auth0Client))

	// ETL
	etl := Etl{
		AcsisDb: *app.AcsisDb,
		EmtctDb: *app.EmtctDb,
	}
	eltRouter := r.PathPrefix("/api/etl").Subrouter()
	eltRouter.HandleFunc("/pregnancies", authMid.Then(etl.PregnancyEtlHandler)).
		Methods(http.MethodOptions, http.MethodPost)

	// Infants
	pregnancies := pregnancy.Pregnancies{EmtctDb: app.EmtctDb}
	inf := infant.New(app.AcsisDb.DB)
	infantRoutes := InfantRoutes{
		Infant:      infant.Infants{Acsis: inf.Acsis},
		AcsisDb:     *app.AcsisDb,
		Pregnancies: pregnancies,
	}
	infantRouter := r.PathPrefix("/api/infants").Subrouter()
	infantRouter.HandleFunc("/diagnoses/{infantId}", authMid.Then(infantRoutes.InfantDiagnosesHandler)).
		Methods(http.MethodOptions, http.MethodGet)
	infantRouter.HandleFunc("/{patientId}", authMid.Then(infantRoutes.InfantHandlers)).
		Methods(http.MethodOptions, http.MethodGet)

	// Patients
	patients := patient.New(app.AcsisDb.DB)
	patientRouter := r.PathPrefix("/api/patients").Subrouter()

	// HomeVisits
	visits := homeVisits.New(app.EmtctDb.DB)
	homeVisitRoutes := HomeVisitRoutes{
		HomeVisits: visits,
		Patients:   patients,
	}
	homeVisitsRouter := r.PathPrefix("/api/homeVisits").Subrouter()
	homeVisitsRouter.HandleFunc("/{homeVisitId}", authMid.Then(homeVisitRoutes.HomeVisitsHandler)).
		Methods(http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut)
	patientRouter.HandleFunc("/{id}/homeVisits", authMid.Then(homeVisitRoutes.FindByPatientHandler)).
		Methods(http.MethodPost, http.MethodGet, http.MethodOptions)

	// Admissions
	hospitalAdmissions := admissions.New(app.EmtctDb.DB)
	admissionRoutes := AdmissionRoutes{
		Admissions: hospitalAdmissions,
		Patients:   patients,
	}
	admissionRouter := r.PathPrefix("/api/hospitalAdmissions").Subrouter()
	admissionRouter.HandleFunc("/", authMid.Then(admissionRoutes.AdmissionsHandler)).
		Methods(http.MethodPost, http.MethodPut, http.MethodOptions)
	patientRouter.HandleFunc("/{patientId}/hospitalAdmissions", authMid.Then(admissionRoutes.AdmissionsByPatientHandler)).
		Methods(http.MethodOptions, http.MethodGet)

	return r
}
