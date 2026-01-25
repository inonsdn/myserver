package connection

func getRoutes() map[string]RouteHandlerFunc {
	return map[string]RouteHandlerFunc{
		"/": Home,
	}
}

func Home(rh *RouteHandler) {
	rh.sendResponse()
}
