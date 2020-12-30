package auth

import (
	"bytes"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	log "github.com/sirupsen/logrus"
	"moh.gov.bz/mch/emtct/internal/db"
	"net/http"
)

type User struct {
	ID          string   `json:"id" firestore:"id"`
	FirstName   string   `json:"firstName" firestore:"firstName"`
	LastName    string   `json:"lastName" firestore:"lastName"`
	Email       string   `json:"email" firestore:"email"`
	Permissions []string `json:"permissions" firestore:"permissions"`
}

type UserStore struct {
	db          *db.FirestoreClient
	collection  string
	adminClient *firebase.App
	authClient  *auth.Client
	apiKey      string
}

func NewStore(db *db.FirestoreClient, config *firebase.Config, apiKey string) (UserStore, error) {
	app, err := firebase.NewApp(db.Ctx, config)
	if err != nil {
		return UserStore{}, err
	}
	authClient, err := app.Auth(db.Ctx)
	if err != nil {
		return UserStore{}, err
	}

	return UserStore{
		db:          db,
		collection:  "emtct_users",
		authClient:  authClient,
		adminClient: app,
		apiKey:      apiKey,
	}, nil
}

func (s *UserStore) ctx() context.Context {
	return s.db.Ctx
}

func (s *UserStore) CreateUser(user User) error {
	u := (&auth.UserToCreate{}).Email(user.Email).DisplayName(fmt.Sprintf("%s %s", user.FirstName, user.LastName)).Disabled(false)
	userRecord, err := s.authClient.CreateUser(s.ctx(), u)
	if err != nil {
		return err
	}
	//Create userRecord in firestore
	user.ID = userRecord.UID
	_, err = s.db.Client.Collection(s.collection).Doc(user.ID).Set(s.ctx(), user)
	if err != nil {
		return err
	}
	err = s.SendPasswordResetEmail(user.Email)
	return err
}

func (s *UserStore) SendPasswordResetEmail(email string) error {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=%s", s.apiKey)
	reqBody, _ := json.Marshal(map[string]string{
		"requestType": "PASSWORD_RESET",
		"email":       email,
	})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"body": body,
	}).Info("result")
	return nil
}
