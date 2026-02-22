package database

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLStatementInfo struct {
	cols []string
	vals []any
}

func (s *SQLStatementInfo) GetStatementCols() string {
	return strings.Join(s.cols, ",")
}

func (s *SQLStatementInfo) GetStatementArgs() string {
	argsStatements := []string{}
	for i := 1; i <= len(s.vals); i++ {
		argsStatements = append(argsStatements, fmt.Sprintf("$%d", i))
	}
	return strings.Join(argsStatements, ",")
}

func (s *SQLStatementInfo) GetUpdateStatementArgs() string {
	argsStatements := []string{}
	for i := 1; i <= len(s.vals); i++ {
		argsStatements = append(argsStatements, fmt.Sprintf("$%d", i))
	}
	return strings.Join(argsStatements, ",")
}

func SQLStructExtraction[T any](s T) SQLStatementInfo {
	key := reflect.TypeOf(s)
	val := reflect.ValueOf(s)

	colNames := []string{}
	valNames := []any{}

	for i := 0; i < key.NumField(); i++ {
		colName := key.Field(i).Tag.Get("json")
		colVal := val.Field(i)
		colNames = append(colNames, colName)
		valNames = append(valNames, colVal)
	}

	return SQLStatementInfo{
		cols: colNames,
		vals: valNames,
	}
}

type SQLValue[T any] struct {
	Val    T
	isNull bool
}

func SQLVal[T any](val T) SQLValue[T] {
	return SQLValue[T]{
		Val: val,
	}
}

func SQLNull[T any]() SQLValue[T] {
	return SQLValue[T]{
		isNull: true,
	}
}

type Uuid struct {
	uuid uuid.UUID
}

func (u *Uuid) ToString() string {
	return u.uuid.String()
}

func UuidCvt(v [16]uint8) uuid.UUID {
	return uuid.UUID(v)
}

func UuidCvtFromDb(v any) uuid.UUID {
	if v == nil {
		return uuid.Nil
	}
	return uuid.UUID(v.([16]uint8))
}

// func UuidCvt(v any) Uuid {
// 	var u uuid.UUID
// 	switch x := v.(type) {
// 	case [16]uint8:
// 		u = uuid.UUID(x)
// 	}
// 	return Uuid{
// 		uuid: u,
// 	}
// }

type BaseRepo struct {
	pool     *pgxpool.Pool
	executor DBExecutor
}

type Record struct {
	id                  uuid.UUID
	createdTimestamp    time.Time
	lastUpdateTimestamp time.Time
}

// Interface for execute to databse, expected method that must have
//
//	Connect to connect to database server
//	Disconnect to disconnect from database server
//	Query for query item to database and expected to return map of column name to value and error if existed
type DBExecutor interface {
	Connect() error
	Disconnect()
	Execute(context.Context, string, ...any) (int64, error)
	Query(context.Context, string, ...any) ([]map[string]any, error)
	QueryRow(context.Context, []any, string, ...any) error
}

// One of executor is Postgres database to connect to Supabase
type PGExecutor struct {
	config *pgxpool.Config
	pool   *pgxpool.Pool
}

func NewPGExecutor(config *pgxpool.Config) *PGExecutor {
	return &PGExecutor{
		config: config,
	}
}

func (pg *PGExecutor) Execute(ctx context.Context, statement string, args ...any) (int64, error) {

	tag, err := pg.pool.Exec(ctx, statement, args...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (pg *PGExecutor) QueryRow(ctx context.Context, dest []any, statement string, args ...any) error {
	err := pg.pool.QueryRow(ctx, statement, args...).Scan(dest...)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PGExecutor) Query(ctx context.Context, statement string, args ...any) ([]map[string]any, error) {
	allRows := []map[string]any{}
	rows, err := pg.pool.Query(ctx, statement, args...)
	if err != nil {
		return allRows, err
	}
	defer rows.Close()
	fds := rows.FieldDescriptions()
	for rows.Next() {
		values, _ := rows.Values()
		row := map[string]any{}
		for i, fd := range fds {
			row[string(fd.Name)] = values[i]
		}
		allRows = append(allRows, row)
	}
	return allRows, nil
}

func (pg *PGExecutor) Connect() error {
	poolCon, err := pgxpool.NewWithConfig(context.Background(), pg.config)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	pg.pool = poolCon
	return nil
}

func (pg *PGExecutor) Disconnect() {
	if pg.pool != nil {
		pg.pool.Close()
	}
}

// func (r *BaseRepo) WithTransaction(ctx context.Context, fn func(pgx.Tx) (any, error)) error {
// 	tx, err := r.pool.Begin(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback(ctx)
// 	_, err = fn(tx)
// 	if err != nil {
// 		tx.Rollback(ctx)
// 	} else {
// 		tx.Commit(ctx)
// 	}
// 	return nil
// }

func (r *BaseRepo) execute(statement string, args ...any) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r.pool.Exec(ctx, statement, args...)
}

func (r *BaseRepo) query(statement string, args ...any) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

}

type DatabaseHandler struct {
	db        DBExecutor
	user      *UserRepo
	scheduler *SchedulerJobRepo
	notes     *NotesRepo
}

func registerRepo(dh *DatabaseHandler) {
	RegisterRepo_User(dh)
	RegisterRepo_Notes(dh)
}

func Connect(db DBExecutor) *DatabaseHandler {
	db.Connect()
	databaseHandler := DatabaseHandler{
		db: db,
	}
	registerRepo(&databaseHandler)

	return &databaseHandler
}

func (d *DatabaseHandler) GetUserConnection() *UserRepo {
	return d.user
}

func (d *DatabaseHandler) GetSchedulerJobConnection() *SchedulerJobRepo {
	return d.scheduler
}

func (d *DatabaseHandler) GetNotesConnection() *NotesRepo {
	return d.notes
}
