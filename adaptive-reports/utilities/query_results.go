package utilities

import (
	"database/sql"

	"github.com/pkg/errors"
)

// QueryResults encapsulates query results in a convenient tabular form
type QueryResults struct {
	Columns           []string
	Rows              [][]string
	ColumnIndexByName map[string]int
}

// GetValue returns string value at the given row/column
func (qr QueryResults) GetValue(column string, row int) string {
	col, ok := qr.ColumnIndexByName[column]
	if !ok {
		panic(errors.Errorf("There is no column %s", column))
	}
	return qr.Rows[row][col]
}

// RunQuery runs the provided query and returns QueryResults
func (db *Database) RunQuery(query string, arguments ...interface{}) (queryResults QueryResults, err error) {
	defer recoverToErrorVar("RunQuery", &err)
	// Execute the query
	var stmtOut *sql.Stmt
	stmtOut, err = db.db.Prepare(query)
	err = errors.Wrap(err, "Couldn't prepare query")
	defer stmtOut.Close()
	if err == nil {
		var databaseRows *sql.Rows
		databaseRows, err = stmtOut.Query(arguments...)
		err = errors.Wrapf(err, "Couldn't Query with arguments %v", arguments)
		if err == nil {
			defer databaseRows.Close()
			queryResults, err = ConvertRowsToQueryResults(databaseRows)
		}
	}
	return
}

// ConvertRowsToQueryResults converts sql.Rows to
func ConvertRowsToQueryResults(databaseRows *sql.Rows) (queryResults QueryResults, err error) {
	defer recoverToErrorVar("ConvertRowsToQueryResults", &err)
	queryResults.Columns, err = databaseRows.Columns()
	err = errors.Wrap(err, "Couldn't get Columns")
	if err == nil {
		queryResults.Rows = make([][]string, 0)
		
		// Make a slice for the values
		values := make([]sql.RawBytes, len(queryResults.Columns))

		// rows.Scan wants '[]interface{}' as an argument, so we must copy the
		// references into such a slice
		// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		// Fetch rows
		rowNum := 0
		for databaseRows.Next() {
			// get RawBytes from data
			err = databaseRows.Scan(scanArgs...)
			if err == nil {
				queryResults.Rows = append(queryResults.Rows, make([]string, len(queryResults.Columns)))
				columnNum := 0
				for i, col := range values {
					// Here we can check if the value is nil (NULL value)
					if col == nil {
						queryResults.Rows[rowNum][i] = "NULL"
					} else {
						queryResults.Rows[rowNum][i] = string(col)
					}
					columnNum++
				}
				rowNum++
				err = databaseRows.Err()
			}
		}
		queryResults.ColumnIndexByName = ReverseIndex(queryResults.Columns)
	}

	return
}

// ReverseIndex builds an index from string to it's position
func ReverseIndex(strs []string) (index map[string]int) {
	index = make(map[string]int)
	for i, c := range strs {
		index[c] = i
	}
	return
}
