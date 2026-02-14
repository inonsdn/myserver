package connection

import (
	"fmt"
	"net/http"
)

func getRoutes() map[string]RouteHandlerFunc {
	return map[string]RouteHandlerFunc{
		"/":        Home,
		"/getUser": GetUsers,
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
