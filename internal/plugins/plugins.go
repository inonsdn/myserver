package plugins

import (
	"myserver/internal/connection"
	"myserver/internal/plugins/notes"
	"myserver/internal/plugins/user"
)

func GetRoutePath() []connection.RoutePathHandler {
	allPath := []connection.RoutePathHandler{}
	allPath = append(allPath, user.UserRoutePath...)
	allPath = append(allPath, notes.NoteRoutePath...)
	return allPath
}
