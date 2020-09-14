package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
)

type ColumnInfo struct {
	Name string
	Type string
}

type TableInfo struct {
	Name string

	Columns   []ColumnInfo
	RowsCount int64
}

func GetRowsCount(ctx context.Context, client *spanner.Client, table string) (int64, error) {
	stmt := spanner.Statement{SQL: `SELECT COUNT(*) as count FROM ` + table}

	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var rowsCount int64
	err := iter.Do(func(row *spanner.Row) error {
		if err := row.Columns(&rowsCount); err != nil {
			return fmt.Errorf("spanner.Row.Columns() error for table %v: %v", table, err)
		}
		return nil
	})
	return rowsCount, err
}

func GetTableInfos(ctx context.Context, databasePath string) ([]TableInfo, error) {
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

	error := iter.Do(func(row *spanner.Row) error {
		var columnName, tableName, spannerType string
		row.Columns(&columnName, &tableName, &spannerType)

		index, ok := tables[tableName]
		if !ok {
			index = len(tableInfos)
			tables[tableName] = index

			var tableInfo TableInfo
			tableInfo.Name = tableName
			rowsCount, err := GetRowsCount(ctx, client, tableName)
			if err != nil {
				return err
			}
			LogTableRowsCountLoad(ctx, databasePath+"/"+tableName)

			tableInfo.RowsCount = rowsCount
			tableInfos = append(tableInfos, tableInfo)
			LogTableInfoLoad(ctx, databasePath+"/"+tableName)
		}

		tableInfos[index].Columns = append(tableInfos[index].Columns, ColumnInfo{Name: columnName, Type: spannerType})

		return nil
	})

	return tableInfos, error
}
