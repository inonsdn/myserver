package user

import (
	"fmt"
	"myserver/internal/connection"
	"net/http"
)

var userRoutePath = []connection.RoutePathHandler{
	{
		Method:      http.MethodGet,
		Path:        "/users",
		Handler:     GetUsers,
		RequireAuth: false,
	},
}

func GetUsers(rh *connection.RouteHandler) error {
	fmt.Println("Get user")
	userCon := rh.DbHandler.GetUserConnection()
	allUsers := userCon.GetAllUser()
	fmt.Println("All users: ")
	fmt.Println(allUsers)

	rh.ResponseJSON(http.StatusOK, allUsers)

	return nil
}
