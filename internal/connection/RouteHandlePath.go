package connection

import (
	"fmt"
	"myserver/internal/database"
	"net/http"

	"github.com/google/uuid"
)

var routePath = []RoutePathHandler{
	{
		Method:  http.MethodGet,
		Path:    "/getAllUsers",
		Handler: GetUsers,
	},
	{
		Method:  http.MethodGet,
		Path:    "/myNotes/{userId}",
		Handler: GetMyNotes,
	},
	{
		Method:  http.MethodPost,
		Path:    "/notes",
		Handler: CreateNewNotes,
	},
	{
		Method:  http.MethodPut,
		Path:    "/notes/{id}",
		Handler: UpdateNotes,
	},
	{
		Method:  http.MethodPost,
		Path:    "/noteGroup",
		Handler: CreateNewNoteGroup,
	},
}

func getRoutes() map[string]RouteHandlerFunc {
	return map[string]RouteHandlerFunc{
		"/":         Home,
		"/getUser":  GetUsers,
		"/newNotes": CreateNewNotes,
	}
}

func Home(rh *RouteHandler) error {
	fmt.Println("HOME CALLED")
	rh.Response(http.StatusOK, "Hello")
	return nil
}

func GetUsers(rh *RouteHandler) error {
	fmt.Println("Get user")
	userCon := rh.dbHandler.GetUserConnection()
	allUsers := userCon.GetAllUser()
	fmt.Println("All users: ")
	fmt.Println(allUsers)

	rh.ResponseJSON(http.StatusOK, allUsers)

	return nil
}

func CreateNewNotes(rh *RouteHandler) error {
	fmt.Println("CreateNewNotes")
	req := database.CreateNotes{}

	if err := rh.GetJSON(&req); err != nil {
		rh.ResponseError(http.StatusBadRequest, "Invalid body")
		return nil
	}

	noteCon := rh.dbHandler.GetNotesConnection()

	notesId := noteCon.CreateNotes(req)

	rh.ResponseJSON(http.StatusOK, map[string]any{
		"id": notesId,
	})

	return nil
}

func CreateNewNoteGroup(rh *RouteHandler) error {
	fmt.Println("CreateNewNoteGroup")
	req := database.CreateGroupNote{}

	if err := rh.GetJSON(&req); err != nil {
		rh.ResponseError(http.StatusBadRequest, "Invalid body")
		return nil
	}

	noteCon := rh.dbHandler.GetNotesConnection()

	notesId := noteCon.CreateGroupNote(req)

	rh.ResponseJSON(http.StatusOK, map[string]any{
		"id": notesId,
	})

	return nil
}

func UpdateNotes(rh *RouteHandler) error {
	id, err := uuid.Parse(rh.r.PathValue("id"))
	if err != nil {
		rh.ResponseError(http.StatusBadRequest, "Invalid UUID")
		return err
	}
	req := database.UpdateNotes{}

	if err := rh.GetJSON(&req); err != nil {
		rh.ResponseError(http.StatusBadRequest, "Invalid body")
		return nil
	}

	fmt.Println(req.NoteGroup)

	noteCon := rh.dbHandler.GetNotesConnection()
	noteCon.UpdateNotes(id, req.Text, req.Title, req.NoteGroup)

	rh.ResponseJSON(http.StatusOK, map[string]any{
		"id": id,
	})

	return nil
}

func GetMyNotes(rh *RouteHandler) error {
	fmt.Println("Get user")
	userId, err := uuid.Parse(rh.r.PathValue("userId"))
	if err != nil {
		rh.ResponseError(http.StatusBadRequest, "Invalid UUID")
		return err
	}
	noteCon := rh.dbHandler.GetNotesConnection()
	notes := noteCon.GetAllNotes(userId)

	rh.ResponseJSON(http.StatusOK, notes)

	return nil
}
