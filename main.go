package main

import (
	"context"
	"fmt"
)

func main() {
	// some windows stuff
	// os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "C:\\Users\\alex-\\OneDrive\\Desktop\\google-cloud-key.json")

	ctx := context.Background()

	databases := ListDatabases(ctx)
	databaseInfos := GetDatabaseInfos(ctx, databases)

	for _, databaseInfo := range databaseInfos {
		fmt.Println(databaseInfo.ToJson())
	}
}
