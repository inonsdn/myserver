package dbcon

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type SqlCon struct {
	driver     string
	dataSource string
}

type SqlResult interface {
	update(rows *sql.Rows)
}

type DbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewSqlCon(cfg *DbConfig) *SqlCon {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	return &SqlCon{
		driver:     "mysql",
		dataSource: dsn,
	}
}

// //////////////////////////////////////
//
// Helper Function

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

	dest.update(rows)

	return nil
}

// //////////////////////////////////////
//
// Public Function

type UserResult struct {
	name string
	id   string
}

// implement method update of interface sql result
func (u *UserResult) update(rows *sql.Rows) {
	err := rows.Scan(u.name, u.id)
	if err != nil {
		fmt.Println("Found error when get from rows", err)
	}
}

func (s *SqlCon) GetUserById(userId string) *UserResult {
	userResult := UserResult{}
	s.query(&userResult, "SELECT * FROM User WHERE id = ?", userId)
	return &userResult
}
