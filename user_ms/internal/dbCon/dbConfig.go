package dbcon

var schemaStrList = []string{
	`CREATE TABLE IF NOT EXISTS user (
		id BLOB NOT NULL,
		name VARCHAR(255) NOT NULL,
		PRIMARY KEY (id)
	)`,
}
