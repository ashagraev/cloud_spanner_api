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
	flag.Parse()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "verbose", verbose)

	return ctx
}
