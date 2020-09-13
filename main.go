package main

import (
	"fmt"
	"os"
)

func main() {
	ctx := PrepareContext()

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

	err = ReportDatabases(ctx, databaseInfos)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
