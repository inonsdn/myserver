package notes

import "myserver/internal/database"

type NotesHandler struct {
	dbHandler *database.DatabaseHandler
}
