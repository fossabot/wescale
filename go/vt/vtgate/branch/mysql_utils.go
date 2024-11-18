package branch

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

func connectToMysql(host string, port int, username, password, connectionOpt string) (*sql.DB, error) {
	var dsn string

	// Handle different authentication scenarios
	switch {
	case username == "" && password == "":
		// Connect without authentication
		dsn = fmt.Sprintf("@tcp(%s:%d)/%s", host, port, connectionOpt)
	case username == "":
		// Username is required if password is provided
		return nil, fmt.Errorf("username is required when password is provided")
	case password == "":
		// Connect with username but no password
		dsn = fmt.Sprintf("%s@tcp(%s:%d)/%s", username, host, port, connectionOpt)
	default:
		// Connect with both username and password
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, connectionOpt)
	}

	// Establish database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping MySQL: %v", err)
	}
	return db, nil
}

// GetAllDatabases retrieves all database names from MySQL
func GetAllDatabases(host string, port int, username, password string) ([]string, error) {
	db, err := connectToMysql(host, port, username, password, "")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Execute query to get all database names
	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		return nil, fmt.Errorf("failed to query database list: %v", err)
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, fmt.Errorf("failed to scan database name: %v", err)
		}
		databases = append(databases, dbName)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating database list: %v", err)
	}

	return databases, nil
}

// GetAllCreateTableStatements retrieves CREATE TABLE statements for all tables in all databases
// Returns a nested map where the first level key is the database name,
// second level key is the table name, and the value is the CREATE TABLE statement
func GetAllCreateTableStatements(host string, port int, username, password string, databasesExclude []string) (map[string]map[string]string, error) {
	db, err := connectToMysql(host, port, username, password, "information_schema?multiStatements=true")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// First step: Get information about all tables and build the combined query
	buildQuery := `
        SELECT CONCAT( 'SHOW CREATE TABLE ', TABLE_SCHEMA, '.', TABLE_NAME, ';' ) AS show_stmt,
               TABLE_SCHEMA,
               TABLE_NAME
        FROM information_schema.TABLES 
        WHERE TABLE_TYPE = 'BASE TABLE'
        `

	if databasesExclude != nil && len(databasesExclude) > 0 {
		buildQuery += fmt.Sprintf(" AND TABLE_SCHEMA NOT IN ('%s')", strings.Join(databasesExclude, "','"))
	}

	rows, err := db.Query(buildQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query table information: %v", err)
	}
	defer rows.Close()

	// Collect all statements and table information
	var showStatements []string
	type tableInfo struct {
		schema string
		name   string
	}
	tableInfos := make([]tableInfo, 0)

	for rows.Next() {
		var showStmt, schema, tableName string
		if err := rows.Scan(&showStmt, &schema, &tableName); err != nil {
			return nil, fmt.Errorf("failed to scan query result: %v", err)
		}
		showStatements = append(showStatements, showStmt)
		tableInfos = append(tableInfos, tableInfo{schema: schema, name: tableName})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating query results: %v", err)
	}

	// Build the combined query
	combinedQuery := strings.Join(showStatements, "")

	// Execute the combined query to get all CREATE TABLE statements at once
	multiRows, err := db.Query(combinedQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute combined query: %v", err)
	}
	defer multiRows.Close()

	// Initialize result map
	result := make(map[string]map[string]string)

	// Process each result set
	for i := 0; i < len(tableInfos); i++ {
		schema := tableInfos[i].schema
		tableName := tableInfos[i].name

		// Ensure database map is initialized
		if _, exists := result[schema]; !exists {
			result[schema] = make(map[string]string)
		}

		// Each SHOW CREATE TABLE result has two columns: table name and create statement
		if !multiRows.Next() {
			return nil, fmt.Errorf("unexpected end of result sets while processing %s.%s", schema, tableName)
		}

		var tableNameResult, createTableStmt string
		if err := multiRows.Scan(&tableNameResult, &createTableStmt); err != nil {
			return nil, fmt.Errorf("failed to scan create table result for %s.%s: %v", schema, tableName, err)
		}

		// Store the result
		result[schema][tableName] = createTableStmt

		// Move to next result set
		if i < len(tableInfos)-1 {
			if !multiRows.NextResultSet() {
				return nil, fmt.Errorf("failed to move to next result set after processing %s.%s", schema, tableName)
			}
		}
	}

	return result, nil
}

func ExecuteSQL(host string, port int, username, password, query string) error {
	db, err := connectToMysql(host, port, username, password, "information_schema?multiStatements=true")
	if err != nil {
		return err
	}
	defer db.Close()

	// use Exec instead of Query since we're not expecting any rows to be returned
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to execute SQL statements: %w", err)
	}

	return nil
}

func ExecuteSQLInTxn(host string, port int, username, password string, queries []string) error {
	db, err := connectToMysql(host, port, username, password, "")
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, query := range queries {
		_, err := tx.Exec(query)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
