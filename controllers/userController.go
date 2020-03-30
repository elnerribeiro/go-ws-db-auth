package controllers

import (
	"encoding/json"
	"net/http"

	"strconv"

	repo "github.com/elnerribeiro/go-ws-db-auth/repositories"
	services "github.com/elnerribeiro/go-ws-db-auth/services"
	u "github.com/elnerribeiro/go-ws-db-auth/utils"
	"github.com/gorilla/mux"
)

//ListUsers Lists all users
var ListUsers = func(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(repo.ContextKey("role")).(string)
	if role != "admin" {
		resp := u.Message(false, "Usuário sem permissão")
		u.Respond(w, resp)
		return
	}
	account := &repo.User{}
	data, err := services.ListUsers(account)
	if err != nil {
		resp := u.Message(false, "Erro ao buscar usuarios")
		u.Respond(w, resp)
		return
	}
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

//GetUserByID Get an user by ID
var GetUserByID = func(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(repo.ContextKey("role")).(string)
	if role != "admin" {
		resp := u.Message(false, "Usuário sem permissão")
		u.Respond(w, resp)
		return
	}
	vars := mux.Vars(r)
	uID, _ := strconv.Atoi(vars["id"])
	account := &repo.User{}
	account.ID = uID
	data, err := services.GetUserByID(account)
	if err != nil {
		resp := u.Message(false, "Erro ao buscar usuario")
		u.Respond(w, resp)
		return
	}
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

//Upsert Inserts or updates an user
var Upsert = func(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(repo.ContextKey("role")).(string)
	if role != "admin" {
		resp := u.Message(false, "Usuário sem permissão")
		u.Respond(w, resp)
		return
	}
	account := &repo.User{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	us, err2 := services.Upsert(account)
	if err2 != nil {
		resp := u.Message(false, "Erro ao atualizar usuario")
		u.Respond(w, resp)
		return
	}
	resp := u.Message(true, "success")
	resp["data"] = us
	u.Respond(w, resp)
}

//Delete Deletes an user by ID
var Delete = func(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(repo.ContextKey("role")).(string)
	if role != "admin" {
		resp := u.Message(false, "Usuário sem permissão")
		u.Respond(w, resp)
		return
	}
	vars := mux.Vars(r)
	uID, _ := strconv.Atoi(vars["id"])
	account := &repo.User{}
	account.ID = uID
	if err := services.Delete(account); err != nil {
		resp := u.Message(false, "Erro ao remover usuario")
		u.Respond(w, resp)
		return
	}
	u.Respond(w, u.Message(true, "success"))
}
