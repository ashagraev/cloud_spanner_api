package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/sync/errgroup"

	databasepb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

type DatabaseInfo struct {
	Path  string
	State string

	Tables []*TableInfo `json:"Tables,omitempty"`
}

func (databaseInfo DatabaseInfo) ToJSON(pretty bool) (string, error) {
	simpleJSON, err := json.Marshal(databaseInfo)
	if err != nil {
		return "", fmt.Errorf("json.Marshal(%v) error: %v", databaseInfo.Path, err)
	}

	if !pretty {
		return string(simpleJSON), nil
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, simpleJSON, "", "  "); err != nil {
		return "", fmt.Errorf("json.Indent(%v) error: %v", databaseInfo.Path, err)
	}
	return prettyJSON.String(), nil
}

func (db *DatabaseClient) GetDatabaseInfo(ctx context.Context, databasePath string) (*DatabaseInfo, error) {
	var databaseInfo DatabaseInfo
	databaseInfo.Path = databasePath

	getDatabaseRequest := &databasepb.GetDatabaseRequest{
		Name: databasePath,
	}

	resp, err := db.databaseAdminClient.GetDatabase(ctx, getDatabaseRequest)
	if err != nil {
		return &databaseInfo, fmt.Errorf("DatabaseAdminClient.GetDatabase(%v) error: %v", databasePath, err)
	}
	LogDatabaseStateLoad(ctx, databasePath)

	databaseInfo.State = resp.GetState().String()

	tc, err := NewTableClient(ctx, databasePath)
	if err != nil {
		return &databaseInfo, err
	}

	if ctx.Value("no-tables") == false {
		tables, err := tc.GetTableInfos(ctx)
		if err != nil {
			return &databaseInfo, err
		}
		databaseInfo.Tables = tables
	}

	return &databaseInfo, nil
}

func (db *DatabaseClient) GetDatabaseInfos(ctx context.Context, databasePaths []string) ([]*DatabaseInfo, error) {
	databaseInfos := make([]*DatabaseInfo, len(databasePaths))
	errs, ctx := errgroup.WithContext(ctx)
	for databaseIdx := range databasePaths {
		databaseIdx := databaseIdx // https://golang.org/doc/faq#closures_and_goroutines

		setupDatabaseInfo := func(idx int) error {
			dbInfo, err := db.GetDatabaseInfo(ctx, databasePaths[idx])
			if err != nil {
				return err
			}
			databaseInfos[idx] = dbInfo
			return nil
		}

		if ctx.Value("no-goroutines") == true {
			if err := setupDatabaseInfo(databaseIdx); err != nil {
				return databaseInfos, err
			}
		} else {
			errs.Go(func() error {
				return setupDatabaseInfo(databaseIdx)
			})
		}
	}
	err := errs.Wait()
	return databaseInfos, err
}
