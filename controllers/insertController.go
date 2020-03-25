package controllers

import (
	"net/http"

	"strconv"

	services "github.com/elnerribeiro/go-ws-db-auth/services"
	u "github.com/elnerribeiro/go-ws-db-auth/utils"
	"github.com/gorilla/mux"
)

//ClearInserts Clears the database
var ClearInserts = func(w http.ResponseWriter, r *http.Request) {
	err := services.ClearInserts()
	if err != nil {
		resp := u.Message(false, "Erro ao limpar batches")
		u.Respond(w, resp)
		return
	}
	u.Respond(w, u.Message(true, "success"))
}

//ListInsert Lists one insert batch
var ListInsert = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uID, _ := strconv.Atoi(vars["id"])
	data, err := services.ListInserts(uID)
	if err != nil {
		resp := u.Message(false, "Erro ao buscar batch")
		u.Respond(w, resp)
		return
	}
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

//InsertSync Inserts a batch of given quantity sync
var InsertSync = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qty, _ := strconv.Atoi(vars["qty"])
	data, err := services.InsertBatchSync(qty)
	if err != nil {
		resp := u.Message(false, "Erro ao inserir batch")
		u.Respond(w, resp)
		return
	}
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

//InsertASync Inserts a batch of given quantity async
var InsertASync = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qty, _ := strconv.Atoi(vars["qty"])
	data, err := services.InsertBatchASync(qty)
	if err != nil {
		resp := u.Message(false, "Erro ao inserir batch")
		u.Respond(w, resp)
		return
	}
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
