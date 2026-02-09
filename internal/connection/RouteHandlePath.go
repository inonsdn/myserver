package connection

import "fmt"

func getRoutes() map[string]RouteHandlerFunc {
	return map[string]RouteHandlerFunc{
		"/":        Home,
		"/getUser": Home,
	}
}

func Home(rh *RouteHandler) {
	fmt.Println("HOME CALLED")
	rh.sendResponse()
}
