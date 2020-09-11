package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
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

func GetRowsCount(client *spanner.Client, table string) int64 {
	ctx := context.Background()
	stmt := spanner.Statement{SQL: `SELECT COUNT(*) as count FROM ` + table}

	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		var rowsCount int64
		if err := row.Columns(&rowsCount); err != nil {
			fmt.Println(err)
			break
		}

		return rowsCount
	}

	return -1
}

func GetTableInfos(databasePath string) []TableInfo {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, databasePath)
	if err != nil {
		fmt.Println(err)
		return []TableInfo{}
	}
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

	tables := make(map[string]TableInfo)

	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		var columnName, tableName, spannerType string
		if err := row.Columns(&columnName, &tableName, &spannerType); err != nil {
			fmt.Println(err)
			break
		}

		tableInfo := tables[tableName]
		tableInfo.Name = tableName
		tableInfo.Columns = append(tableInfo.Columns, ColumnInfo{Name: columnName, Type: spannerType})
		tables[tableName] = tableInfo
	}

	var tableInfos []TableInfo
	for _, tableInfo := range tables {
		tableInfo.RowsCount = GetRowsCount(client, tableInfo.Name)
		tableInfos = append(tableInfos, tableInfo)
	}
	return tableInfos
}
