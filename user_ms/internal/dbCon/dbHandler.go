package dbcon

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type SqlCon struct {
	driver     string
	dataSource string
	db         *sql.DB
}

type SqlResult interface {
	update(rows *sql.Rows) error
}

type DbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewSqlCon(cfg *DbConfig) (*SqlCon, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Cannot open sql")
		return nil, err
	}
	return &SqlCon{
		driver:     "mysql",
		dataSource: dsn,
		db:         db,
	}, nil
}

// //////////////////////////////////////
//
// Helper Function

func (s *SqlCon) execute(statement string, args ...any) (any, error) {
	// db, err := sql.Open(s.driver, s.dataSource)
	// if err != nil {
	// 	fmt.Println("Cannot open sql")
	// 	return nil, err
	// }
	// defer db.Close()

	result, err := s.db.Exec(statement, args...)
	if err != nil {
		fmt.Println("Cannot execute statement.")
		return nil, err
	}

	return result.RowsAffected()
}

func (s *SqlCon) query(dest SqlResult, statement string, args ...any) error {
	// db, err := sql.Open(s.driver, s.dataSource)
	// if err != nil {
	// 	fmt.Println("Cannot open sql")
	// 	return err
	// }
	// defer db.Close()

	rows, err := s.db.Query(statement, args...)
	defer rows.Close()

	return dest.update(rows)
}

// //////////////////////////////////////
//
// Public Function
func (s *SqlCon) InitializeSchema() error {
	// db, err := sql.Open(s.driver, s.dataSource)
	// if err != nil {
	// 	fmt.Println("Cannot open sql")
	// 	return err
	// }
	// defer db.Close()

	for _, statement := range schemaStrList {
		_, err := s.db.Exec(statement)
		if err != nil {
			break
		}
	}

	return nil
}

type UserResult struct {
	Name string
	Id   string
}

// implement method update of interface sql result
func (u *UserResult) update(rows *sql.Rows) error {
	if !rows.Next() {
		// no rows returned
		if err := rows.Err(); err != nil {
			fmt.Println("rows error:", err)
			return err
		}
		// no error, just no data
		return sql.ErrNoRows
	}
	err := rows.Scan(&u.Name, &u.Id)
	if err != nil {
		fmt.Println("Found error when get from rows", err)
		return err
	}
	return nil
}

func (s *SqlCon) GetUserById(userId string) (*UserResult, error) {
	userResult := UserResult{}
	err := s.query(&userResult, "SELECT * FROM user WHERE id = ?", userId)
	return &userResult, err
}
