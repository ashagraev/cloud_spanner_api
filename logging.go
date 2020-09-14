package main

import (
	"context"
	"log"
)

func logLoad(ctx context.Context, title string, name string) {
	if ctx.Value("verbose") != true {
		return
	}
	log.Println("loaded " + title + ": " + name)
}

func logProjectLoad(ctx context.Context, projectName string) {
	logLoad(ctx, "project", projectName)
}

func logInstanceLoad(ctx context.Context, instanceName string) {
	logLoad(ctx, "instance", instanceName)
}

func logDataBaseLoad(ctx context.Context, databaseName string) {
	logLoad(ctx, "database", databaseName)
}

func logDatabaseStateLoad(ctx context.Context, databaseName string) {
	logLoad(ctx, "database state", databaseName)
}

func logTableRowsCountLoad(ctx context.Context, tableName string) {
	logLoad(ctx, "table rows count", tableName)
}

func logTableInfoLoad(ctx context.Context, tableName string) {
	logLoad(ctx, "tables info", tableName)
}
