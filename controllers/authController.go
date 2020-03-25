package controllers

import (
	"encoding/json"
	"net/http"

	repo "github.com/elnerribeiro/go-ws-db-auth/repositories"
	services "github.com/elnerribeiro/go-ws-db-auth/services"
	u "github.com/elnerribeiro/go-ws-db-auth/utils"
)

//Authenticate do user authentication
var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &repo.User{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := services.Login(account.Email, account.Password)
	u.Respond(w, resp)
}

//Validate do user validation - gets ID
var Validate = func(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value(repo.ContextKey("user")).(int)
	role := r.Context().Value(repo.ContextKey("role")).(string)
	resp := u.Message(true, "success")
	resp["userId"] = id
	resp["role"] = role
	u.Respond(w, resp)
}
