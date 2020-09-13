package main

import (
	"bytes"
	"context"
	"encoding/json"

	"golang.org/x/sync/errgroup"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	databasepb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

type DatabaseInfo struct {
	Path  string
	State string

	Tables []TableInfo
}

func (databaseInfo DatabaseInfo) ToJson() (string, error) {
	resp, err := json.Marshal(databaseInfo)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func (databaseInfo DatabaseInfo) ToJsonPretty() (string, error) {
	simpleJson, err := json.Marshal(databaseInfo)
	if err != nil {
		return "", err
	}

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, simpleJson, "", "  ")

	return prettyJSON.String(), nil
}

func GetDatabaseInfo(ctx context.Context, databasePath string) (DatabaseInfo, error) {
	var databaseInfo DatabaseInfo
	databaseInfo.Path = databasePath

	getDatabaseRequest := &databasepb.GetDatabaseRequest{
		Name: databasePath,
	}

	databaseAdminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return databaseInfo, err
	}
	defer databaseAdminClient.Close()

	resp, err := databaseAdminClient.GetDatabase(ctx, getDatabaseRequest)
	if err != nil {
		return databaseInfo, err
	}
	LogDatabaseStateLoad(ctx, databasePath)

	databaseInfo.State = resp.GetState().String()

	if ctx.Value("no-tables") == false {
		tables, err := GetTableInfos(ctx, databasePath)
		if err != nil {
			return databaseInfo, err
		}
		databaseInfo.Tables = tables
	}

	return databaseInfo, nil
}

func GetDatabaseInfos(ctx context.Context, databasePaths []string) ([]DatabaseInfo, error) {
	databaseInfos := make([]DatabaseInfo, len(databasePaths))
	errs, ctx := errgroup.WithContext(ctx)
	for databaseIdx := range databasePaths {
		databaseIdx := databaseIdx // https://golang.org/doc/faq#closures_and_goroutines

		setupDatabaseInfo := func(idx int) error {
			dbInfo, err := GetDatabaseInfo(ctx, databasePaths[idx])
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
