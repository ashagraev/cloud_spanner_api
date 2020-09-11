package main

import (
	"fmt"
)

func main() {
	// some windows stuff
	// os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "C:\\Users\\alex-\\OneDrive\\Desktop\\google-cloud-key.json")

	projects := ListProjects()
	instances := ListInstances(projects)
	databases := ListDatabases(instances)
	databaseInfos := GetDatabaseInfos(databases)

	for _, databaseInfo := range databaseInfos {
		fmt.Println(databaseInfo.ToJsonPretty())
	}
}
