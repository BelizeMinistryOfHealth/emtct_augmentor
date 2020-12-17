package db

import (
	"cloud.google.com/go/firestore"
	"context"
)

type FirestoreClient struct {
	Client    *firestore.Client
	Ctx       context.Context
	projectId string
}

func (c *FirestoreClient) Close() error {
	return c.Client.Close()
}

func NewFirestore(ctx context.Context, projectId string) (*FirestoreClient, error) {
	c, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		return nil, err
	}
	client := &FirestoreClient{
		Client:    c,
		Ctx:       ctx,
		projectId: projectId,
	}
	return client, nil
}
