package tests

import (
	"context"
	firebase "firebase.google.com/go/v4"
	log "github.com/sirupsen/logrus"
	"moh.gov.bz/mch/emtct/internal/auth"
	"moh.gov.bz/mch/emtct/internal/db"
	"os"
	"testing"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint: true,
	})
	log.SetOutput(os.Stdout)
}

func Test_Auth_CreateUser(t *testing.T) {
	apiKey := os.Getenv("API_KEY")
	projectId := os.Getenv("GCP_PROJECT_ID")

	ctx := context.Background()
	firestoreDb, err := db.NewFirestore(ctx, projectId)
	if err != nil {
		t.Fatalf("failed to create firestore connection: %v", err)
	}
	config := firebase.Config{
		ProjectID: projectId,
	}
	userStore, err := auth.NewStore(firestoreDb, &config, apiKey)
	if err != nil {
		t.Fatalf("failed to create user store: %v", err)
	}
	user := auth.User{
		FirstName:   "Uris",
		LastName:    "Guerra",
		Email:       "uris77@gmail.com",
		Permissions: []string{"app:read", "app:write"},
	}
	err = userStore.CreateUser(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
}

func Test_Auth_GetUserByEmail(t *testing.T) {
	apiKey := os.Getenv("API_KEY")
	projectId := os.Getenv("GCP_PROJECT_ID")
	ctx := context.Background()
	firestoreDb, err := db.NewFirestore(ctx, projectId)
	if err != nil {
		t.Fatalf("failed to create firestore connection: %v", err)
	}
	config := firebase.Config{
		ProjectID: projectId,
	}
	userStore, err := auth.NewStore(firestoreDb, &config, apiKey)
	if err != nil {
		t.Fatalf("failed to create user store: %v", err)
	}
	email := "uris77@gmail.com"
	user, err := userStore.GetUserByEmail(email)
	if err != nil {
		t.Fatalf("GetUserByEmail error: %v", err)
	}
	t.Logf("user: %v", user)
}

func Test_Auth_DeleteUser(t *testing.T) {
	apiKey := os.Getenv("API_KEY")
	projectId := os.Getenv("GCP_PROJECT_ID")
	ctx := context.Background()
	firestoreDb, err := db.NewFirestore(ctx, projectId)
	if err != nil {
		t.Fatalf("failed to create firestore connection: %v", err)
	}
	config := firebase.Config{
		ProjectID: projectId,
	}
	userStore, err := auth.NewStore(firestoreDb, &config, apiKey)
	if err != nil {
		t.Fatalf("failed to create user store: %v", err)
	}
	email := "uris77@gmail.com"
	user, err := userStore.GetUserByEmail(email)
	if err != nil {
		t.Fatalf("GetUserByEmail error: %v", err)
	}
	err = userStore.DeleteUser(user)
	if err != nil {
		t.Fatalf("DeleteUser error: %v", err)
	}
}

func Test_Auth_UpdateUser(t *testing.T) {
	apiKey := os.Getenv("API_KEY")
	projectId := os.Getenv("GCP_PROJECT_ID")
	ctx := context.Background()
	firestoreDb, err := db.NewFirestore(ctx, projectId)
	if err != nil {
		t.Fatalf("failed to create firestore connection: %v", err)
	}
	config := firebase.Config{
		ProjectID: projectId,
	}
	userStore, err := auth.NewStore(firestoreDb, &config, apiKey)
	if err != nil {
		t.Fatalf("failed to create user store: %v", err)
	}
	email := "uris77@gmail.com"
	user, err := userStore.GetUserByEmail(email)
	if err != nil {
		t.Fatalf("GetUserByEmail error: %v", err)
	}
	user.Permissions = []string{"admin:read", "admin:write"}
	if err := userStore.UpdateUser(user); err != nil {
		t.Fatalf("failed to update user: %v", err)
	}
	t.Logf("user permissions: %v", user.Permissions)
}
