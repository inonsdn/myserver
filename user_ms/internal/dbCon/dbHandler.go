package dbcon

import (
	"database/sql"
	"fmt"
)

type SqlCon struct {
	driver     string
	dataSource string
}

type SqlResult interface {
	Update(rows *sql.Rows)
}

type UserResult struct {
	name string
	id   string
}

func (u *UserResult) Update(rows *sql.Rows) {
	err := rows.Scan(u.name, u.id)
	if err != nil {
		fmt.Println("Found error when get from rows", err)
	}
}

// func NewDbHandler() *SqlCon {
// 	db, err := sql.Open("mysql", "user:password@/database-name")
// 	if err != nil {
// 		fmt.Println("Cannot open sql")
// 		return nil
// 	}
// 	return &SqlCon{
// 		db: db,
// 	}
// }

func (s *SqlCon) execute(statement string, args ...any) (any, error) {
	db, err := sql.Open(s.driver, s.dataSource)
	if err != nil {
		fmt.Println("Cannot open sql")
		return nil, err
	}
	defer db.Close()

	result, err := db.Exec(statement, args...)
	if err != nil {
		fmt.Println("Cannot execute statement.")
		return nil, err
	}

	return result.RowsAffected()
}

func (s *SqlCon) query(dest SqlResult, statement string, args ...any) error {
	db, err := sql.Open(s.driver, s.dataSource)
	if err != nil {
		fmt.Println("Cannot open sql")
		return err
	}
	defer db.Close()

	rows, err := db.Query(statement, args...)

	dest.Update(rows)

	return nil
}

func (s *SqlCon) GetUserById(userId string) {
	userResult := UserResult{}
	s.query(&userResult, "SELECT * FROM User WHERE id = ?", userId)
}
