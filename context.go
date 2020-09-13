package main

import (
	"context"
	"flag"
)

func PrepareContext() context.Context {
	// some windows stuff
	// os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "C:\\Users\\alex-\\OneDrive\\Desktop\\google-cloud-key.json")

	verbose := false
	flag.BoolVar(&verbose, "verbose", verbose, "be verbose")

	noTablesExport := false
	flag.BoolVar(&noTablesExport, "no-tables", noTablesExport, "do not export detailed tables information")

	flag.Parse()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "verbose", verbose)
	ctx = context.WithValue(ctx, "no-tables", noTablesExport)

	return ctx
}
