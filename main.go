package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	// some windows stuff
	// os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "C:\\Users\\alex-\\OneDrive\\Desktop\\google-cloud-key.json")

	ctx := context.Background()

	databases, err := ListDatabases(ctx)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	databaseInfos, err := GetDatabaseInfos(ctx, databases)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	for _, databaseInfo := range databaseInfos {
		databaseJson, err := databaseInfo.ToJson()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			continue
		}
		fmt.Println(databaseJson)
	}
}
