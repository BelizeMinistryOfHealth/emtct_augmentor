package db

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type FirestoreClient struct {
	Client      *firestore.Client
	AuthClient  *auth.Client
	AdminClient *firebase.App
	Ctx         context.Context
	projectId   string
}

func (c *FirestoreClient) Close() error {
	return c.Client.Close()
}

func NewFirestore(ctx context.Context, projectId string) (*FirestoreClient, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}
	authClient, _ := app.Auth(ctx)

	c, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		return nil, err
	}
	client := &FirestoreClient{
		Client:      c,
		AuthClient:  authClient,
		AdminClient: app,
		Ctx:         ctx,
		projectId:   projectId,
	}
	return client, nil
}
