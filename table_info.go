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

type TableClient struct {
	client *spanner.Client
	path   string
}

func NewTableClient(ctx context.Context, databasePath string) (*TableClient, error) {
	client, err := spanner.NewClientWithConfig(ctx, databasePath, spanner.ClientConfig{
		SessionPoolConfig: spanner.SessionPoolConfig{
			MinOpened: 1,
			MaxOpened: 2,
			MaxIdle:   1,
			MaxBurst:  2,
		}})
	if err != nil {
		return nil, fmt.Errorf("spanner.NewClientWithConfig(%v) error: %v", databasePath, err)
	}
	return &TableClient{client: client}, nil
}

func (tc *TableClient) Close() {
	tc.client.Close()
}

func (tc *TableClient) GetRowsCount(ctx context.Context, table string) (int64, error) {
	stmt := spanner.Statement{SQL: `SELECT COUNT(*) as count FROM ` + table}

	iter := tc.client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var rowsCount int64
	row, err := iter.Next()
	if err != nil {
		return rowsCount, fmt.Errorf("iter.Next() error for table %v: %v", table, err)
	}

	if err := row.Columns(&rowsCount); err != nil {
		return rowsCount, fmt.Errorf("spanner.Row.Columns() error for table %v: %v", table, err)
	}

	return rowsCount, nil
}

func (tc *TableClient) GetTableInfos(ctx context.Context) ([]*TableInfo, error) {
	stmt := spanner.Statement{SQL: `
		SELECT
			column_name,
			table_name,
			spanner_type
		FROM
			information_schema.columns
		WHERE table_catalog = '' AND table_schema = ''
	`}

	iter := tc.client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var tableInfos []*TableInfo
	tables := make(map[string]int)

	err := iter.Do(func(row *spanner.Row) error {
		var columnName, tableName, spannerType string
		err := row.Columns(&columnName, &tableName, &spannerType)
		if err != nil {
			return fmt.Errorf("spanner.Row.Columns() error processing %v: %v", tc.path, err)
		}

		index, ok := tables[tableName]
		if !ok {
			index = len(tableInfos)
			tables[tableName] = index

			var tableInfo TableInfo
			tableInfo.Name = tableName
			rowsCount, err := tc.GetRowsCount(ctx, tableName)
			if err != nil {
				return err
			}
			LogTableRowsCountLoad(ctx, tc.path+"/"+tableName)

			tableInfo.RowsCount = rowsCount
			tableInfos = append(tableInfos, &tableInfo)
			LogTableInfoLoad(ctx, tc.path+"/"+tableName)
		}

		tableInfos[index].Columns = append(tableInfos[index].Columns, ColumnInfo{Name: columnName, Type: spannerType})

		return nil
	})

	return tableInfos, err
}
