package main

import (
	"context"
	"flag"
)

func prepareContext() context.Context {
	verbose := false
	flag.BoolVar(&verbose, "verbose", verbose, "turn on loading process logging")

	noTablesExport := false
	flag.BoolVar(&noTablesExport, "no-tables", noTablesExport, "do not export detailed tables information")

	noGoroutines := false
	flag.BoolVar(&noGoroutines, "no-goroutines", noGoroutines, "do not use goroutines while loading databases information")

	jsonLines := false
	flag.BoolVar(&jsonLines, "json-lines", jsonLines, "export each database information on a separate json line")

	jsonPretty := false
	flag.BoolVar(&jsonPretty, "pretty", jsonPretty, "prettify json output")

	flag.Parse()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "verbose", verbose)
	ctx = context.WithValue(ctx, "no-tables", noTablesExport)
	ctx = context.WithValue(ctx, "no-goroutines", noGoroutines)

	ctx = context.WithValue(ctx, "json-lines", jsonLines)
	ctx = context.WithValue(ctx, "json-pretty", jsonPretty)

	return ctx
}
