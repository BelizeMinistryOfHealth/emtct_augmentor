package models

import "moh.gov.bz/mch/emtct/internal/models/permissions"

type Role struct {
	App        string                 `json:"app" firestore:"app"`
	Permission permissions.Permission `json:"permission" firestore:"permission"`
}

type User struct {
	Email     string `json:"email" firestore:"email"`
	FirstName string `json:"firstName" firestore:"firstName"`
	LastName  string `json:"lastName" firestore:"lastName"`
	Roles     []Role `json:"roles" firestore:"roles"`
}
