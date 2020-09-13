package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	databasepb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

type DatabaseInfo struct {
	Path  string
	State string

	Tables []TableInfo
}

func (databaseInfo DatabaseInfo) ToJson() string {
	resp, _ := json.Marshal(databaseInfo)
	return string(resp)
}

func (databaseInfo DatabaseInfo) ToJsonPretty() string {
	simpleJson, _ := json.Marshal(databaseInfo)

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, simpleJson, "", "  ")

	return prettyJSON.String()
}

func GetDatabaseInfo(ctx context.Context, databasePath string) DatabaseInfo {
	var databaseInfo DatabaseInfo
	databaseInfo.Path = databasePath

	getDatabaseRequest := &databasepb.GetDatabaseRequest{
		Name: databasePath,
	}

	databaseAdminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		fmt.Println(err)
		return DatabaseInfo{}
	}
	defer databaseAdminClient.Close()

	resp, err := databaseAdminClient.GetDatabase(ctx, getDatabaseRequest)
	if err != nil {
		fmt.Println(err)
		return databaseInfo
	}

	databaseInfo.State = resp.GetState().String()
	databaseInfo.Tables = GetTableInfos(ctx, databasePath)

	return databaseInfo
}

func GetDatabaseInfos(ctx context.Context, databasePaths []string) []DatabaseInfo {
	databaseInfos := make([]DatabaseInfo, len(databasePaths))
	var wg sync.WaitGroup
	for databaseIdx := range databasePaths {
		wg.Add(1)
		go func(idx int) {
			databaseInfos[idx] = GetDatabaseInfo(ctx, databasePaths[idx])
			wg.Done()
		}(databaseIdx)
	}
	wg.Wait()
	return databaseInfos
}
