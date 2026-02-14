package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

const (
	REPO_USER_NAME = "user"
)

type UserRepo struct {
	*BaseRepo
}

type User struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func (u *User) GetId() string {
	return u.Id.String()
}

func RegisterRepo_User(dh *DatabaseHandler) {
	dh.user = &UserRepo{
		BaseRepo: &BaseRepo{
			executor: dh.db,
		},
	}
}

func (u *UserRepo) CreateNewUser(username string) uuid.UUID {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	statement := "INSERT INTO users (username) VALUES ($1) RETURNING id"
	columnToValues, err := u.executor.Query(ctx, statement)
	// err := u.pool.QueryRow(ctx, statement, username).Scan(&userId)
	if err != nil {
		slog.Error(err.Error())
		return uuid.Nil
	}
	userId := UuidCvt(columnToValues[0]["id"].([16]uint8))
	return userId
}

func (u *UserRepo) GetAllUser() []User {
	allUsers := []User{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	statement := "SELECT id, username FROM users"
	columnToValues, err := u.executor.Query(ctx, statement)
	if err != nil {
		slog.Error(err.Error())
		return allUsers
	}
	for _, columnToValue := range columnToValues {
		allUsers = append(allUsers, User{
			Id:       UuidCvt(columnToValue["id"].([16]uint8)),
			Username: columnToValue["username"].(string),
		})
	}
	return allUsers
}
