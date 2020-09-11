package main

import (
	"cloud.google.com/go/spanner"
	"context"
)

type ColumnInfo struct {
	Name string
	Type string
}

type TableInfo struct {
	Name string

	Columns []ColumnInfo
	RowsCount int64
}

func GetRowsCount(ctx context.Context, client *spanner.Client, table string) int64 {
	stmt := spanner.Statement{SQL: `SELECT COUNT(*) as count FROM ` + table}

	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var rowsCount int64
	iter.Do(func(row *spanner.Row) error {
		if err := row.Columns(&rowsCount); err != nil {
			return err
		}
		return nil
	})
	return rowsCount
}

func GetTableInfos(ctx context.Context, databasePath string) []TableInfo {
	client, _ := spanner.NewClientWithConfig(ctx, databasePath, spanner.ClientConfig{
		SessionPoolConfig: spanner.SessionPoolConfig{
			MinOpened: 1,
			MaxOpened: 2,
			MaxIdle:   1,
			MaxBurst:  2,
		}})
	defer client.Close()

	stmt := spanner.Statement{SQL: `
		SELECT
			column_name,
			table_name,
			spanner_type
		FROM
			information_schema.columns
		WHERE table_catalog = '' AND table_schema = ''
	`}

	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var tableInfos []TableInfo
	tables := make(map[string]int)

	iter.Do(func(row *spanner.Row) error {
		var columnName, tableName, spannerType string
		row.Columns(&columnName, &tableName, &spannerType)

		index, ok := tables[tableName]; if !ok {
			index = len(tableInfos)
			tables[tableName] = index

			var tableInfo TableInfo
			tableInfo.Name = tableName
			tableInfo.RowsCount = GetRowsCount(ctx, client, tableName)
			tableInfos = append(tableInfos, tableInfo)
		}

		tableInfos[index].Columns = append(tableInfos[index].Columns, ColumnInfo{Name: columnName, Type: spannerType})

		return nil
	})

	return tableInfos
}
