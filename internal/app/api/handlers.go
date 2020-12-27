// Package api contains the full set of handler functions and routes
// supported by the web api.
package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"moh.gov.bz/mch/emtct/internal/app"
	"moh.gov.bz/mch/emtct/internal/business/data/admissions"
	"moh.gov.bz/mch/emtct/internal/business/data/contactTracing"
	"moh.gov.bz/mch/emtct/internal/business/data/contraceptives"
	"moh.gov.bz/mch/emtct/internal/business/data/hiv"
	"moh.gov.bz/mch/emtct/internal/business/data/homeVisits"
	"moh.gov.bz/mch/emtct/internal/business/data/infant"
	"moh.gov.bz/mch/emtct/internal/business/data/labs"
	"moh.gov.bz/mch/emtct/internal/business/data/partners"
	"moh.gov.bz/mch/emtct/internal/business/data/patient"
)

func API(app app.App) *mux.Router {
	r := mux.NewRouter()

	// Instantiate an aut0 client with a Cache with a key capacity of
	// 60 tokens and a ttl of 24 hours.
	//auth0Client := auth0.NewAuth0(60, 518400)

	// Middleware that verifies JWT token and also enables CORS.
	authMid := NewChain(EnableCors(), VerifyToken(app.Firestore))

	// Patients
	patients := patient.New(app.Firestore)
	patientRouter := r.PathPrefix("/api/patients").Subrouter()

	lab := labs.New(app.AcsisDb)

	// Infants
	inf := infant.New(app.Firestore)
	infantRoutes := InfantRoutes{
		Infant:  inf,
		Labs:    lab,
		Patient: patients,
	}

	// HomeVisits
	visits := homeVisits.New(app.Firestore)
	homeVisitRoutes := HomeVisitRoutes{
		HomeVisits: visits,
		Patients:   patients,
	}

	//Returns a home visit that matches the specific homeVisitId
	patientRouter.HandleFunc("/{patientId}/pregnancy/{pregnancyId}/homeVisits/{homeVisitId}", authMid.Then(homeVisitRoutes.HomeVisitsHandler)).
		Methods(http.MethodOptions, http.MethodGet)
	//Mutation Routes
	patientRouter.HandleFunc("/{patientId}/pregnancy/{pregnancyId}/homeVisits", authMid.Then(homeVisitRoutes.HomeVisitsHandler)).
		Methods(http.MethodOptions, http.MethodPost, http.MethodPut)
	//Returns a list of home visits for the pregnancy id
	patientRouter.HandleFunc("/{patientId}/pregnancy/{pregnancyId}/homeVisits", authMid.Then(homeVisitRoutes.FindByPregnancyHandler)).
		Methods(http.MethodGet, http.MethodOptions)

	// Admissions
	hospitalAdmissions := admissions.New(app.Firestore)
	admissionRoutes := AdmissionRoutes{
		Admissions: hospitalAdmissions,
		Patients:   patients,
	}
	patientRouter.HandleFunc("/{patientId}/pregnancy/{pregnancyId}/hospitalAdmissions", authMid.Then(admissionRoutes.AdmissionsHandler)).
		Methods(http.MethodPost, http.MethodPut, http.MethodOptions)
	patientRouter.HandleFunc("/{patientId}/pregnancy/{pregnancyId}/hospitalAdmissions", authMid.Then(admissionRoutes.AdmissionsByPregnancyHandler)).
		Methods(http.MethodOptions, http.MethodGet)

	// Contraceptives
	contraceptive := contraceptives.New(app.Firestore)
	contraceptiveRoutes := ContraceptivesRoutes{
		Contraceptives: contraceptive,
		Patients:       patients,
	}
	patientRouter.HandleFunc("/{patientId}/pregnancy/{pregnancyId}/contraceptivesUsed", authMid.Then(contraceptiveRoutes.ContraceptivesHandler)).
		Methods(http.MethodPost, http.MethodPut, http.MethodOptions)
	patientRouter.HandleFunc("/{patientId}/pregnancy/{pregnancyId}/contraceptivesUsed", authMid.Then(contraceptiveRoutes.ContraceptivesByPatientHandler)).
		Methods(http.MethodOptions, http.MethodGet)

	// Partners Router
	partnerRoutes := partnersRoutes{
		Patient:  patients,
		Partners: partners.New(app.Firestore),
	}

	// Contact Tracing
	tracing := contactTracing.New(app.Firestore)
	tracingRoutes := ContactTracingRoutes{
		ContactTracings: tracing,
		Patient:         patients,
	}
	patientRouter.HandleFunc("/{patientId}/pregnancy/{pregnancyId}/contactTracing", authMid.Then(tracingRoutes.ContactTracingHandler)).
		Methods(http.MethodOptions, http.MethodPost, http.MethodPut)
	patientRouter.HandleFunc("/{patientId}/pregnancy/{pregnancyId}/contactTracing", authMid.Then(tracingRoutes.ContactTracingHandler)).
		Methods(http.MethodOptions, http.MethodGet)
	patientRouter.HandleFunc("/{patientId}/pregnancy/{pregnancyId}/syphilisTreatments", authMid.Then(partnerRoutes.SyphilisTreatmentHandler)).
		Methods(http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut)

	// Pregnancies
	Hiv := hiv.New(app.AcsisDb)
	pregRoutes := pregnancyRoutes{Patient: patients, Lab: lab, Hiv: Hiv}
	patientRouter.HandleFunc("/{patientId}/pregnancy/{pregnancyId}",
		authMid.Then(pregRoutes.GetPregnancy)).Methods(http.MethodOptions, http.MethodGet)
	patientRouter.HandleFunc("/{patientId}/pregnancy/{pregnancyId}/labResults",
		authMid.Then(pregRoutes.FindPregnancyLabResults)).Methods(http.MethodOptions, http.MethodGet)

	patientRouter.HandleFunc("/{patientId}/pregnancy/{pregnancyId}/arvs", authMid.Then(pregRoutes.ArvsHandler)).
		Methods(http.MethodOptions, http.MethodGet)
	patientRouter.HandleFunc("/{patientId}/pregnancy/{pregnancyId}/syphilisTreatments", authMid.Then(pregRoutes.PatientSyphilisTreatmentHandler)).
		Methods(http.MethodOptions, http.MethodGet)
	patientRouter.HandleFunc("/{motherId}/pregnancy/{pregnancyId}/infant", authMid.Then(pregRoutes.InfantHandlers)).
		Methods(http.MethodOptions, http.MethodGet)
	patientRouter.HandleFunc("/{motherId}/pregnancy/{pregnancyId}/infant/{infantId}", authMid.Then(pregRoutes.InfantByIdHandlers)).
		Methods(http.MethodOptions, http.MethodGet)
	patientRouter.HandleFunc("/{patientId}/infant/{infantId}/diagnoses", authMid.Then(pregRoutes.InfantDiagnosesHandler)).
		Methods(http.MethodOptions, http.MethodGet)
	patientRouter.HandleFunc("/{motherId}/infant/{infantId}/syphilisTreatments", authMid.Then(pregRoutes.InfantSyphilisTreatmentHandler)).
		Methods(http.MethodGet, http.MethodOptions)
	patientRouter.HandleFunc("/{motherId}/infant/{infantId}/syphilisScreenings", authMid.Then(pregRoutes.InfantSyphilisScreeningHandler)).
		Methods(http.MethodOptions, http.MethodGet)

	patientRouter.HandleFunc("/{motherId}/infant/{infantId}/hivScreenings", authMid.Then(infantRoutes.HivScreeningHandler)).
		Methods(http.MethodOptions, http.MethodPost, http.MethodPut, http.MethodGet)

	patientRouter.HandleFunc("/{id}", authMid.Then(pregRoutes.RetrievePatientHandler)).
		Methods(http.MethodOptions, http.MethodGet)

	return r
}
