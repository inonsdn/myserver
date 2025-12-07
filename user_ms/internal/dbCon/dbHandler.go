package dbcon

import (
	"database/sql"
	"fmt"
	"strings"

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

type Inserter interface {
	TableName() string
	ColumnNames() []string
	ColumnValues() []interface{}
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

func (s *SqlCon) execute(statement string, args []any) (any, error) {

	stmt, err := s.db.Prepare(statement)
	if err != nil {
		return nil, err
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		fmt.Println("Cannot execute statement.")
		return nil, err
	}

	return result.RowsAffected()
}

func (s *SqlCon) query(dest SqlResult, statement string, args ...any) error {

	rows, err := s.db.Query(statement, args...)
	defer rows.Close()

	if err != nil {
		fmt.Println("Found error when query", err)
		return nil
	}

	if !rows.Next() {
		// no rows returned
		if err := rows.Err(); err != nil {
			fmt.Println("rows error:", err)
			return err
		}
		// no error, just no data
		return sql.ErrNoRows
	}

	return dest.update(rows)
}

// //////////////////////////////////////
//
// Public Function
func (s *SqlCon) InitializeSchema() error {
	db, err := sql.Open(s.driver, s.dataSource)
	if err != nil {
		fmt.Println("Cannot open sql")
		return err
	}
	defer db.Close()

	for _, statement := range schemaStrList {
		_, err := s.db.Exec(statement)
		if err != nil {
			fmt.Println("Found error when init", err)
			break
		}
	}

	return nil
}

type UserResult struct {
	Name string
	Id   string
}

type UserInfo struct {
	Id       string `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u UserInfo) TableName() string {
	return "user"
}

func (u UserInfo) ColumnNames() []string {
	return []string{"id", "username", "password"}
}

func (u UserInfo) ColumnValues() []interface{} {
	return []interface{}{u.Id, u.Username, u.Password}
}

// implement method update of interface sql result
func (u *UserResult) update(rows *sql.Rows) error {

	err := rows.Scan(&u.Id, &u.Name)
	if err != nil {
		fmt.Println("Found error when get from rows", err)
		return err
	}
	return nil
}

func (s *SqlCon) GetUserById(userId string) (*UserResult, error) {
	userResult := UserResult{}
	err := s.query(&userResult, "SELECT id, username FROM user WHERE id = ?", userId)
	return &userResult, err
}

func (s *SqlCon) CreateUser(userInfo UserInfo) error {
	columns := strings.Join(userInfo.ColumnNames(), ",")
	var placeholders []string
	var values []any
	for i, val := range userInfo.ColumnValues() {
		placeholders = append(placeholders, "?")
		fmt.Println(i)
		values = append(values, val)
	}
	params := strings.Join(placeholders, ",")
	statement := fmt.Sprintf("INSERT INTO %s ( %s ) VALUES ( %s )", userInfo.TableName(), columns, params)
	fmt.Println(statement)
	rowAffected, err := s.execute(statement, values)

	fmt.Println("Executed, effect row", rowAffected)

	return err
}

func (s *SqlCon) UserAuthentication(username string, password string) (string, error) {
	userResult := UserResult{}
	err := s.query(&userResult, "SELECT id, username FROM user WHERE username = ? AND password = ?", username, password)
	if err != nil {
		return "", err
	}
	return userResult.Id, err
}
