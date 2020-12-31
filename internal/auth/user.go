package auth

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
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

// NewStore creates a new store that provides ways for creating and mutating a user.
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

// CreateUser creates an auth user and also creates a record in the user collection.
// It will send a password reset email so the user, with a link that allows them to set a password
// for their account.
func (s *UserStore) CreateUser(user User) error {
	u := (&auth.UserToCreate{}).Email(user.Email).DisplayName(fmt.Sprintf("%s %s", user.FirstName, user.LastName)).Disabled(false)
	userRecord, err := s.authClient.CreateUser(s.ctx(), u)
	if err != nil {
		return AuthError{
			Reason: "failed creating auth user",
			Inner:  err,
		}
	}
	userRecord.CustomClaims = map[string]interface{}{"permissions": user.Permissions}
	ur, err := s.authClient.UpdateUser(s.ctx(), userRecord.UID, (&auth.UserToUpdate{}).CustomClaims(map[string]interface{}{"permissions": user.Permissions}))
	if err != nil {
		return AuthError{
			Reason: "failed updating auth user claims",
			Inner:  err,
		}
	}
	log.WithFields(log.Fields{
		"userRecord": *ur,
	}).Info("updated user with claims")
	//Create userRecord in firestore
	user.ID = userRecord.UID
	_, err = s.db.Client.Collection(s.collection).Doc(user.ID).Set(s.ctx(), user)
	if err != nil {
		return AuthError{
			Reason: "failed inserting user to collection",
			Inner:  err,
		}
	}
	if err := s.SendPasswordResetEmail(user.Email); err != nil {
		return AuthError{
			Reason: "failed sending password reset email",
			Inner:  err,
		}
	}
	return nil
}

// SendPasswordResetEmail sends an email to the user with a link to reset their password.
func (s *UserStore) SendPasswordResetEmail(email string) error {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=%s", s.apiKey)
	reqBody, _ := json.Marshal(map[string]string{
		"requestType": "PASSWORD_RESET",
		"email":       email,
	})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return AuthError{
			Reason: "failed posting PASSWORD_RESET",
			Inner:  err,
		}
	}
	defer resp.Body.Close()
	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return AuthError{
			Reason: "failed decoding PASSWORD_RESET response",
			Inner:  err,
		}
	}
	log.WithFields(log.Fields{
		"body": body,
	}).Info("result")
	return nil
}

// DeleteUser deletes a the user from the firestore collection and also from Firestore Auth.
func (s *UserStore) DeleteUser(user User) error {
	_, err := s.authClient.DeleteUsers(s.ctx(), []string{user.ID})
	if err != nil {
		return AuthError{
			Reason: "failed deleting auth user",
			Inner:  err,
		}
	}
	if _, err := s.db.Client.Collection(s.collection).Doc(user.ID).Delete(s.ctx()); err != nil {
		return AuthError{
			Reason: "failed deleting user from collection",
			Inner:  err,
		}
	}
	return nil
}

// GetUserByEmail gets a user's record from firebase.
func (s *UserStore) GetUserByEmail(email string) (User, error) {
	iter := s.db.Client.Collection(s.collection).Where("email", "==", email).Limit(1).Documents(s.ctx())
	var user User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return User{}, AuthError{
				Reason: "failed to fetch user by email",
				Inner:  err,
			}
		}
		if err := doc.DataTo(&user); err != nil {
			return User{}, AuthError{
				Reason: "failed to convert user data",
				Inner:  err,
			}
		}
	}
	u, _ := s.authClient.GetUser(s.ctx(), user.ID)
	log.WithFields(log.Fields{
		"user": *u,
	}).Info("user from auth client")
	return user, nil
}

// UpdateUser updates a user's permissions.
func (s *UserStore) UpdateUser(user User) error {
	userToUpdate := (&auth.UserToUpdate{}).
		CustomClaims(map[string]interface{}{"permissions": user.Permissions})
	if _, err := s.authClient.UpdateUser(s.ctx(), user.ID, userToUpdate); err != nil {
		return AuthError{
			Reason: "failed to update auth user",
			Inner:  err,
		}
	}
	if _, err := s.db.Client.Collection(s.collection).Doc(user.ID).Update(s.ctx(), []firestore.Update{
		{
			Path:  "permissions",
			Value: user.Permissions,
		},
	}); err != nil {
		return AuthError{
			Reason: "failed updating user collection",
			Inner:  err,
		}
	}
	return nil
}
