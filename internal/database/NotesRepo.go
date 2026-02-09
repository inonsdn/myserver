package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

const (
	REPO_NOTES_NAME  = "notes"
	NOTES_TABLE_NAME = "notes"
)

type NotesRepo struct {
	*BaseRepo
}

type Notes struct {
	id     uuid.UUID
	text   string
	userId uuid.UUID
}

func RegisterRepo_Notes(dh *DatabaseHandler) {
	dh.notes = &NotesRepo{
		BaseRepo: &BaseRepo{
			executor: dh.db,
		},
	}
}

func (n *NotesRepo) CreateNotes(userId uuid.UUID, text string) uuid.UUID {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := fmt.Sprintf("INSERT INTO %s (text, user) VALUES ($1, $2) RETURNING id", NOTES_TABLE_NAME)
	var notesId uuid.UUID
	err := n.pool.QueryRow(ctx, statement, text, userId).Scan(&notesId)
	if err != nil {
		slog.Error(err.Error())
		return uuid.Nil
	}
	return notesId
}

func (n *NotesRepo) GetAllNotes(userId uuid.UUID) []Notes {
	allNotes := []Notes{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := fmt.Sprintf("SELECT id, text FROM %s WHERE user = $1", NOTES_TABLE_NAME)
	rows, err := n.pool.Query(ctx, statement, userId)
	if err != nil {
		slog.Error(err.Error())
		return allNotes
	}
	defer rows.Close()

	for rows.Next() {
		var notes Notes
		notes.userId = userId
		if err := rows.Scan(&notes.id, &notes.text); err != nil {
			slog.Error(err.Error())
			return allNotes
		}
		allNotes = append(allNotes, notes)
	}
	return allNotes
}
