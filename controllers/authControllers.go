package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/elnerribeiro/gosbxauth/models"
	u "github.com/elnerribeiro/gosbxauth/utils"
)

//Authenticate autentica o usuario
var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.User{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}

//ListUsers Lista todos os usuarios
var ListUsers = func(w http.ResponseWriter, r *http.Request) {

	//id := r.Context().Value("user").(int)
	data, err := models.ListUsers()
	if err != nil {
		resp := u.Message(false, "Erro ao buscar usuarios")
		u.Respond(w, resp)
		return
	}
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
