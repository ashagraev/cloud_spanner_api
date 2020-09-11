package main

import (
	"bytes"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"context"
	"encoding/json"
	"fmt"
	databasepb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

type DatabaseInfo struct {
	Path string
	State string

	BytesSize int64
	Tables []TableInfo
}

func (databaseInfo DatabaseInfo) ToJson() []byte {
	resp, err := json.Marshal(databaseInfo)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	return resp
}

func (databaseInfo DatabaseInfo) ToJsonPretty() string {
	simpleJson := databaseInfo.ToJson()

	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, simpleJson, "", "  ")
	if error != nil {
		return ""
	}

	return prettyJSON.String()
}

func GetDatabaseInfo(databaseAdminClient *database.DatabaseAdminClient, databasePath string) DatabaseInfo {
	var databaseInfo DatabaseInfo
	databaseInfo.Path = databasePath

	getDatabaseRequest := &databasepb.GetDatabaseRequest{
		Name: databasePath,
	}

	ctx := context.Background()
	resp, err := databaseAdminClient.GetDatabase(ctx, getDatabaseRequest)
	if err != nil {
		fmt.Println(err)
		return databaseInfo
	}

	databaseInfo.State = resp.GetState().String()
	databaseInfo.Tables = GetTableInfos(databasePath)

	return databaseInfo
}

func GetDatabaseInfos(databasePaths []string) []DatabaseInfo {
	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		fmt.Println(err)
		return []DatabaseInfo{}
	}
	defer adminClient.Close()

	var databaseInfos []DatabaseInfo
	for _, databasePath := range databasePaths {
		databaseInfos = append(databaseInfos, GetDatabaseInfo(adminClient, databasePath))
	}
	return databaseInfos
}