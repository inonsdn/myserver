package connection

import (
	"net/http"
)

type RouteHandlerFunc func(*RouteHandler)

// Create handler function for serve http
// by wrapping function
// function must receive argument of route handler
func makeHandler(f RouteHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// new handler for this response
		rh := newRouteHandler(w, r)

		// execute function
		f(rh)
	}
}

type RouteHandler struct {
	w http.ResponseWriter
	r *http.Request
}

func newRouteHandler(w http.ResponseWriter, r *http.Request) *RouteHandler {
	return &RouteHandler{
		w: w,
		r: r,
	}
}

func (rh *RouteHandler) sendResponse() {
	rh.w.Header()
}

func LandingPage(rh *RouteHandler) {
	rh.sendResponse()
}
