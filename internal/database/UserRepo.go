package database

import (
	"context"
	"log/slog"
	"time"
)

const (
	REPO_USER_NAME = "user"
)

type UserRepo struct {
	*BaseRepo
}

type User struct {
	id       string
	username string
}

func RegisterRepo_User(dh *DatabaseHandler) {
	dh.user = &UserRepo{
		BaseRepo: &BaseRepo{
			executor: dh.db,
		},
	}
}

func (u *UserRepo) CreateNewUser(username string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	statement := "INSERT INTO users (username) VALUES ($1) RETURNING id"
	var userId string
	err := u.pool.QueryRow(ctx, statement, username).Scan(&userId)
	if err != nil {
		slog.Error(err.Error())
		return ""
	}
	return userId
}

func (u *UserRepo) GetAllUser() []User {
	allUsers := []User{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	statement := "SELECT id, username FROM users"
	rows, err := u.pool.Query(ctx, statement)
	if err != nil {
		slog.Error(err.Error())
		return allUsers
	}
	defer rows.Close()
	for rows.Next() {
		var userRow User
		if err := rows.Scan(&userRow.id, &userRow.username); err != nil {
			slog.Error(err.Error())
			return allUsers
		}
		allUsers = append(allUsers, userRow)
	}
	return allUsers
}
