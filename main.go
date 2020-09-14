package main

import (
	"fmt"
	"os"
)

func main() {
	ctx := PrepareContext()

	db, err := NewDatabaseClient(ctx)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()

	databases, err := db.ListDatabases(ctx)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	databaseInfos, err := db.GetDatabaseInfos(ctx, databases)
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
