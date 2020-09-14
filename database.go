package main

import (
	"context"
	"fmt"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
)

type DatabaseClient struct {
	databaseAdminClient *database.DatabaseAdminClient
	instanceAdminClient *instance.InstanceAdminClient
}

func NewDatabaseClient(ctx context.Context) (*DatabaseClient, error) {
	databaseAdminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewDatabaseAdminClient() error: %v", err)
	}

	instanceAdminClient, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewInstanceAdminClient error: %v", err)
	}

	return &DatabaseClient{databaseAdminClient: databaseAdminClient, instanceAdminClient: instanceAdminClient}, nil
}

func (db *DatabaseClient) Close() {
	db.databaseAdminClient.Close()
	db.instanceAdminClient.Close()
}
