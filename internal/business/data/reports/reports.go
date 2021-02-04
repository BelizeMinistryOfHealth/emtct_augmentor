package reports

import "moh.gov.bz/mch/emtct/internal/db"

type DbCollections struct {
	HivScreening string
	Patients     string
	Infants      string
}

type Reports struct {
	firestore   *db.FirestoreClient
	collections DbCollections
}

func New(firestore *db.FirestoreClient) Reports {
	return Reports{
		firestore: firestore,
		collections: DbCollections{
			HivScreening: "hivScreenings",
			Patients:     "patients",
			Infants:      "infants",
		},
	}
}
