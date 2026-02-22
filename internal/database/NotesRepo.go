package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

const (
	REPO_NOTES_NAME        = "notes"
	NOTES_TABLE_NAME       = "notes"
	NOTES_GROUP_TABLE_NAME = "noteGroup"
)

type NotesRepo struct {
	*BaseRepo
}

type CreateNotes struct {
	Title     string `json:"title"`
	Text      string `json:"text"`
	UserId    string `json:"user_id"`
	NoteGroup string `json:"note_group"`
}

type UpdateNotes struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	UserId    uuid.UUID `json:"user_id"`
	NoteGroup uuid.UUID `json:"note_group"`
}

type CreateGroupNote struct {
	ParentId SQLValue[uuid.UUID] `json:"parent_id"`
	Name     SQLValue[string]    `json:"name"`
	UserId   uuid.UUID           `json:"user_id"`
}

type UpdateGroupNote struct {
	Id       uuid.UUID        `json:"id"`
	ParentId SQLValue[string] `json:"parent_id"`
	Name     SQLValue[string] `json:"name"`
	UserId   uuid.UUID        `json:"user_id"`
}

type Notes struct {
	Id        uuid.UUID
	Text      string
	UserId    uuid.UUID
	Title     string
	NoteGroup uuid.UUID
}

func RegisterRepo_Notes(dh *DatabaseHandler) {
	dh.notes = &NotesRepo{
		BaseRepo: &BaseRepo{
			executor: dh.db,
		},
	}
}

func (n *NotesRepo) CreateGroupNote(gn CreateGroupNote) uuid.UUID {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sqlStatementInfo := SQLStructExtraction(gn)
	colStr := sqlStatementInfo.GetStatementCols()
	argsStr := sqlStatementInfo.GetStatementArgs()

	statement := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", NOTES_GROUP_TABLE_NAME, colStr, argsStr)
	var notesId uuid.UUID
	err := n.executor.QueryRow(ctx, []any{&notesId}, statement, sqlStatementInfo.vals...)
	if err != nil {
		fmt.Println(err.Error())
		slog.Error(err.Error())
		return uuid.Nil
	}
	return notesId
}

func (n *NotesRepo) CreateNotes(cn CreateNotes) uuid.UUID {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sqlStatementInfo := SQLStructExtraction(cn)
	colStr := sqlStatementInfo.GetStatementCols()
	argsStr := sqlStatementInfo.GetStatementArgs()

	statement := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", NOTES_TABLE_NAME, colStr, argsStr)
	var notesId uuid.UUID
	err := n.executor.QueryRow(ctx, []any{&notesId}, statement, sqlStatementInfo.vals...)
	if err != nil {
		fmt.Println(err.Error())
		slog.Error(err.Error())
		return uuid.Nil
	}
	return notesId
}

func (n *NotesRepo) UpdateNotes(notesId uuid.UUID, text string, title string, noteGroup uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var statement string
	var args []any
	if noteGroup == uuid.Nil {
		args = []any{text, title, notesId.String()}
		statement = fmt.Sprintf("UPDATE %s SET (text, title) = ($1, $2) WHERE id = $3", NOTES_TABLE_NAME)
	} else {
		args = []any{text, title, noteGroup.String(), notesId.String()}
		statement = fmt.Sprintf("UPDATE %s SET (text, title, noteGroup) = ($1, $2, $3) WHERE id = $4", NOTES_TABLE_NAME)
	}
	rowAffect, err := n.executor.Execute(ctx, statement, args...)
	if err != nil {
		fmt.Println(err.Error())
		slog.Error(err.Error())
		return err
	}
	fmt.Println("Row affect", rowAffect)
	return nil
}

func (n *NotesRepo) GetAllNotes(userId uuid.UUID) []Notes {
	allNotes := []Notes{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := fmt.Sprintf("SELECT id, text, title, note_group FROM %s WHERE user_id = $1", NOTES_TABLE_NAME)
	columnToValues, err := n.executor.Query(ctx, statement, userId)
	if err != nil {
		slog.Error(err.Error())
		return allNotes
	}

	for _, columnToValue := range columnToValues {
		allNotes = append(allNotes, Notes{
			Id:        UuidCvtFromDb(columnToValue["id"]),
			Text:      columnToValue["text"].(string),
			UserId:    userId,
			Title:     columnToValue["title"].(string),
			NoteGroup: UuidCvtFromDb(columnToValue["note_group"]),
		})
	}
	return allNotes
}
