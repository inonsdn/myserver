package connection

import (
	"fmt"
	"myserver/internal/database"
	"net/http"
)

var routePath = []RoutePathHandler{
	{
		Method:  http.MethodGet,
		Path:    "/getAllUsers",
		Handler: GetUsers,
	},
	{
		Method:  http.MethodPost,
		Path:    "/notes",
		Handler: CreateNewNotes,
	},
	{
		Method:  http.MethodPut,
		Path:    "/notes/:id",
		Handler: CreateNewNotes,
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

	notesId := noteCon.CreateNotes(req.UserId, req.Text)

	rh.ResponseJSON(http.StatusOK, map[string]any{
		"id": notesId,
	})

	return nil
}

// func UpdateNotes(rh *RouteHandler) error {
// 	id := rh.GetQuery("id")

// 	return nil
// }
