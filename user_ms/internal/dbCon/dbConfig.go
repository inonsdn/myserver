package dbcon

var schemaStrList = []string{
	`CREATE TABLE IF NOT EXISTS user (
		id VARCHAR(36) NOT NULL,
		username VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		PRIMARY KEY (id)
	)`,
}
