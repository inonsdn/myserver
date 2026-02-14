package connection

import (
	"encoding/json"
	"myserver/internal/database"
	"net/http"
)

// type of function to received route handler which is wrapper of connection
type RouteHandlerFunc func(*RouteHandler) error

type RouteHandler struct {
	w         http.ResponseWriter
	r         *http.Request
	dbHandler *database.DatabaseHandler
}

// Create handler function for serve http
// by wrapping function
// function must receive argument of route handler
func makeHandler(f RouteHandlerFunc, dbHandler *database.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// new handler for this response
		rh := &RouteHandler{
			w:         w,
			r:         r,
			dbHandler: dbHandler,
		}

		// execute function
		if err := f(rh); err != nil {
			rh.Response(http.StatusInternalServerError, err.Error())
		}
	}
}

// send response back in JSON format
func (rh *RouteHandler) ResponseJSON(status int, data any) {
	// set header of response to be type of json
	rh.w.Header().Set("Content-Type", "application/json")
	rh.w.WriteHeader(status)

	// encoded data
	if err := json.NewEncoder(rh.w).Encode(data); err != nil {
		http.Error(rh.w, err.Error(), http.StatusInternalServerError)
	}
}

// response function to just response plain text
// with http status
func (rh *RouteHandler) Response(status int, body string) {
	// set header of response
	rh.w.Header().Set("Content-Type", "text/plain")
	rh.w.WriteHeader(status)
	rh.w.Write([]byte(body))
}
