package app

import (
	"net/http"

	u "github.com/elnerribeiro/go-ws-db-auth/utils"
)

//NotFoundHandler Handles 404
var NotFoundHandler = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u.Logger.Info("[NotFoundHandler] Not found!")
		w.WriteHeader(http.StatusNotFound)
		u.Respond(w, u.Message(false, "This resources was not found on our server"))
		next.ServeHTTP(w, r)
	})
}
