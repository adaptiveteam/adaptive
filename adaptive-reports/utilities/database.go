package utilities

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Table map[string]string

type Database struct {
	db *sql.DB
}

// WrapDB converts an existing database connection to our type
func WrapDB(db *sql.DB) *Database {
	return &Database{db: db}
}

// ConnectionString concatenates all arguments into a single connection string
func ConnectionString(
	endPoint,
	userName,
	password,
	port,
	database string) string {
	return userName + ":" + password + "@tcp(" + endPoint + ":" + port + ")/" + database
}

func SQLOpenUnsafe(driver,
	connectionString string) *sql.DB {
	db, err := sql.Open(driver, connectionString)
	if err != nil {
		log.Panicf("Error creating database: %+v", err)
	}
	return db
}

func NewDatabase(
	driver,
	endPoint,
	userName,
	password,
	port,
	database string,
) *Database {
	db := SQLOpenUnsafe(driver, ConnectionString(
		endPoint,
		userName,
		password,
		port,
		database))
	return &Database{db: db}
}

func CloseUnsafe(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Panicf("Error closing database: %+v", err)
	}
}

func (db *Database) CloseDatabase() {
	CloseUnsafe(db.db)
}

func (db *Database) GetTable(query string, arguments ...interface{}) (columns []string, rows [][]string, tableMap Table, err error) {
	defer recoverToErrorVar("GetTable", &err)
	// Execute the query
	var stmtOut *sql.Stmt
	stmtOut, err = db.db.Prepare(query)
	defer stmtOut.Close()
	if err == nil {
		var databaseRows *sql.Rows
		databaseRows, err = stmtOut.Query(arguments...)
		if err == nil {
			columns, err = databaseRows.Columns()
			if err == nil {
				rows = make([][]string, 0)
				tableMap = make(Table)
				// Make a slice for the values
				values := make([]sql.RawBytes, len(columns))

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
						rows = append(rows, make([]string, len(columns)))
						columnNum := 0
						for i, col := range values {
							// Here we can check if the value is nil (NULL value)
							if col == nil {
								rows[rowNum][i] = "NULL"
							} else {
								rows[rowNum][i] = string(col)
							}
							tableMap[GetIndex(columns[columnNum], rowNum)] = rows[rowNum][i]
							columnNum++
						}
						rowNum++
						err = databaseRows.Err()
					}
				}
			}
		}
	} else {
		err = errors.Wrap(err, "Couldn't prepare query")
	}

	return
}

func (table Table) GetValue(column string, rowNum int) (rv string) {
	var ok bool
	rv, ok = table[GetIndex(column, rowNum)]
	if !ok {
		log.Panic("no value at " + GetIndex(column, rowNum))
	}
	return
}

func GetIndex(column string, row int) (rv string) {
	return column + ":" + strconv.Itoa(row)
}

type Query struct {
	query     string
	arguments []string
}


type QueryResult struct {
	Columns []string
	Rows    [][]string
	Table
}

func (q QueryResult) GetColumns() []string {
	return q.Columns
}

func (q QueryResult) GetRows() [][]string {
	return q.Rows
}

func (q QueryResult) GetTable() Table {
	return q.Table
}

type QueryResultMap map[string]QueryResult

func (queries QueryInvocations)Run(
	db *Database,
) (rv QueryResultMap, err error) {
	type results struct {
		purpose     string
		queryResult QueryResult
	}

	r := make(chan results, len(queries))
	var errGroup errgroup.Group
	for _, q := range queries {
		q2 := q
		errGroup.Go(func() error {
			columns, rows, table, tableErr := db.GetTable(q2.SQL, q2.ArgumentValues...)
			if tableErr == nil {
				n := results{
					purpose: q2.Name,
					queryResult: QueryResult{
						Columns: columns,
						Rows:    rows,
						Table:   table,
					},
				}
				r <- n
			}
			return tableErr
		})
	}

	err = errGroup.Wait()
	close(r)
	if err == nil {
		rv = make(QueryResultMap)
		for i := range r {
			rv[i.purpose] = i.queryResult
		}
	}
	return rv, err
}
