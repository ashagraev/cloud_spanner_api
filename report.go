package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

func ReportDatabases(ctx context.Context, databaseInfos []DatabaseInfo) error {
	if ctx.Value("json-lines") == true {
		for _, databaseInfo := range databaseInfos {
			databaseJson, err := databaseInfo.ToJson(false)
			if err != nil {
				return err
			}
			fmt.Println(databaseJson)
		}
		return nil
	}

	simpleJSON, err := json.Marshal(databaseInfos)
	if err != nil {
		return err
	}

	if ctx.Value("json-pretty") == false {
		fmt.Println(string(simpleJSON))
		return nil
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, simpleJSON, "", "  "); err != nil {
		return err
	}
	fmt.Println(prettyJSON.String())

	return nil
}
