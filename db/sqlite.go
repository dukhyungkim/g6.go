package db

import "fmt"

func sqliteTablesQuery() string {
	return `SELECT
    name
FROM
    sqlite_master
WHERE
    type ='table' AND
    name NOT LIKE 'sqlite_%'`
}

func sqliteHasTableQuery(tableName string) string {
	return fmt.Sprintf(`SELECT
    count(*)
FROM
    sqlite_master
WHERE
    type ='table' AND
    name = '%s'`, tableName)
}
