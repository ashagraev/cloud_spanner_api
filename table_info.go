package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
)

// ColumnInfo stores basic information about a specific column in a Cloud Spanner table.
type ColumnInfo struct {
	// Name of a column.
	Name string
	// Type of column elements.
	Type string
}

// TableInfo stores basic information about Cloud Spanner table.
type TableInfo struct {
	// Name of a table.
	Name string

	// List of table columns.
	Columns []*ColumnInfo

	// Number of table rows.
	RowsCount int64
}

// TableClient exposes an interface for collecting Cloud Spanner tables information.
type TableClient struct {
	client *spanner.Client
	path   string
}

// NewTableClient returns a TableClient that holds a connection to Cloud Spanner.
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

// Close() closes TableClient's connection to Cloud Spanner.
func (tc *TableClient) Close() {
	tc.client.Close()
}

// GetRowsCount() returns number of rows in a specific Cloud Spanner table.
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

// GetTableInfos() collects basic information about a list of Spanner tables.
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

			rowsCount, err := tc.GetRowsCount(ctx, tableName)
			if err != nil {
				return err
			}
			logTableRowsCountLoad(ctx, tc.path+"/"+tableName)

			tableInfo := &TableInfo{Name: tableName, RowsCount: rowsCount}
			tableInfos = append(tableInfos, tableInfo)
			logTableInfoLoad(ctx, tc.path+"/"+tableName)
		}

		tableInfos[index].Columns = append(tableInfos[index].Columns, &ColumnInfo{Name: columnName, Type: spannerType})

		return nil
	})

	return tableInfos, err
}
