package app

import "moh.gov.bz/mch/emtct/internal/db"

type App struct {
	Firestore       *db.FirestoreClient
	Auth            Auth
	ProjectID       string
	FirestoreApiKey string
}

type Auth struct {
	JwkUrl string
	Iss    string
	Aud    string
}
