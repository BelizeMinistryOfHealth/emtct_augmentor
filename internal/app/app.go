package app

import "moh.gov.bz/mch/emtct/internal/db"

type App struct {
	Firestore *db.FirestoreClient
	Auth      Auth
}

type Auth struct {
	JwkUrl string
	Iss    string
	Aud    string
}

type JwtToken struct {
	Email string
}
