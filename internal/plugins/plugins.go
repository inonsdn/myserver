package plugins

import (
	"myserver/internal/connection"
	"myserver/internal/plugins/notes"
)

func GetRoutePath() []connection.RoutePathHandler {
	return notes.NoteRoutePath
}
