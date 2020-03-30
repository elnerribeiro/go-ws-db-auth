package controllers

import (
	"net/http"

	"strconv"

	repo "github.com/elnerribeiro/go-ws-db-auth/repositories"
	services "github.com/elnerribeiro/go-ws-db-auth/services"
	u "github.com/elnerribeiro/go-ws-db-auth/utils"
	"github.com/gorilla/mux"
)

//ClearInserts Clears the database
var ClearInserts = func(w http.ResponseWriter, r *http.Request) {
	insert := &repo.Insert{}
	err := services.ClearInserts(insert)
	if err != nil {
		u.Logger.Error("[ClearInserts] Error while cleaning inserts: %s", err)
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
	insert := &repo.Insert{}
	insert.ID = uID
	data, err := services.ListInserts(insert)
	if err != nil {
		u.Logger.Error("[ListInsert] Error listing inserts: %s", err)
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
	insert := &repo.Insert{}
	insert.Quantity = qty
	data, err := services.InsertBatchSync(insert)
	if err != nil {
		u.Logger.Error("[InsertSync] Error inserting batch synchronous: %s", err)
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
	insert := &repo.Insert{}
	insert.Quantity = qty
	data, err := services.InsertBatchASync(insert)
	if err != nil {
		u.Logger.Error("[InsertASync] Error inserting batch asynchronous: %s", err)
		resp := u.Message(false, "Erro ao inserir batch")
		u.Respond(w, resp)
		return
	}
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
