package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

func reportDatabases(ctx context.Context, databaseInfos []*DatabaseInfo) error {
	if ctx.Value("json-lines") == true {
		for _, databaseInfo := range databaseInfos {
			databaseJSON, err := databaseInfo.ToJSON(false)
			if err != nil {
				return err
			}
			fmt.Println(databaseJSON)
		}
		return nil
	}

	simpleJSON, err := json.Marshal(databaseInfos)
	if err != nil {
		return fmt.Errorf("json.Marshal() databases error: %v", err)
	}

	if ctx.Value("json-pretty") == false {
		fmt.Println(string(simpleJSON))
		return nil
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, simpleJSON, "", "  "); err != nil {
		return fmt.Errorf("json.Indent() databases error: %v", err)
	}
	fmt.Println(prettyJSON.String())

	return nil
}
