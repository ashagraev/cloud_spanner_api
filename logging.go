package main

import (
	"context"
	"log"
)

func LogLoad(ctx context.Context, title string, name string) {
	if ctx.Value("verbose") != true {
		return
	}
	log.Println("loaded " + title + ": " + name)
}

func LogProjectLoad(ctx context.Context, projectName string) {
	LogLoad(ctx, "project", projectName)
}

func LogInstanceLoad(ctx context.Context, instanceName string) {
	LogLoad(ctx, "instance", instanceName)
}

func LogDataBaseLoad(ctx context.Context, databaseName string) {
	LogLoad(ctx, "database", databaseName)
}

func LogDatabaseStateLoad(ctx context.Context, databaseName string) {
	LogLoad(ctx, "database state", databaseName)
}

func LogTableRowsCountLoad(ctx context.Context, tableName string) {
	LogLoad(ctx, "table rows count", tableName)
}

func LogTableInfoLoad(ctx context.Context, tableName string) {
	LogLoad(ctx, "tables info", tableName)
}
