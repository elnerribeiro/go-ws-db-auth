package app

import (
	"context"
	"net/http"
	"strings"

	repo "github.com/elnerribeiro/go-ws-db-auth/repositories"
	u "github.com/elnerribeiro/go-ws-db-auth/utils"

	"github.com/dgrijalva/jwt-go"
)

//JwtAuthentication Auth with JWT
var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path         //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			u.Logger.Info("[JwtAuthentication] 403 Token Not found!")
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			u.Logger.Info("[JwtAuthentication] 403 Invalid/Malformed auth token!")
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &repo.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte("JWTpassword123@"), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			u.Logger.Info("[JwtAuthentication] 403 Invalid, expired or malformed token!")
			response = u.Message(false, "Invalid, expired or malformed token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			u.Logger.Info("[JwtAuthentication] 403 Token not valid on this server!")
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		u.Logger.Debug("User %d just logged in", tk.UserID) //Useful for monitoring
		var ctx = context.WithValue(r.Context(), repo.ContextKey("user"), tk.UserID)
		ctx = context.WithValue(ctx, repo.ContextKey("role"), tk.Role)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
