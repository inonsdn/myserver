package connection

import (
	"net/http"
)

// type of function to received route handler which is wrapper of connection
type RouteHandlerFunc func(*RouteHandler)

type RouteHandler struct {
	w http.ResponseWriter
	r *http.Request
}

// Create handler function for serve http
// by wrapping function
// function must receive argument of route handler
func makeHandler(f RouteHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// new handler for this response
		rh := &RouteHandler{
			w: w,
			r: r,
		}

		// execute function
		f(rh)
	}
}

func (rh *RouteHandler) sendResponse() {
	rh.w.Header()
}

func LandingPage(rh *RouteHandler) {
	rh.sendResponse()
}
