package models

import (
	"database/sql"

	db "github.com/elnerribeiro/go-mustache-db"
)

//ListUsers Lists all users
func ListUsers() ([]User, error) {
	var dados db.Dados
	dados = make(db.Dados)
	val, err := db.SelectAll(nil, "getuser", &[]User{}, dados)
	if err != nil {
		return nil, err
	}
	account := val.(*[]User)
	if account != nil {
		//finalAccounts := make([]User, len(*account))
		for i, u := range *account {
			u.Password = ""
			(*account)[i] = u
		}
		return *account, nil
	}
	return *account, nil
}

//Upsert Inserts or updates an user
func Upsert(tx *db.Transacao, user *User) (*User, error) {
	var dados db.Dados
	dados = make(db.Dados)
	userToDados(user, &dados)
	if user.ID == 0 {
		res, err := db.InsertReturningPostgres(tx, "usuario", dados, "id", &User{})
		if err != nil {
			return nil, err
		}
		newID := res.(*User)
		user.ID = newID.ID
		return user, nil
	}

	var filter db.Dados
	filter = make(db.Dados)
	filter["id"] = user.ID
	_, err2 := db.Update(tx, "usuario", dados, filter)
	if err2 != nil {
		return nil, err2
	}
	return user, nil
}

//Delete Deletes an user
func Delete(tx *db.Transacao, uid int) error {
	var filter db.Dados
	filter = make(db.Dados)
	filter["id"] = uid
	_, err := db.Delete(tx, "usuario", filter)
	return err
}

func userToDados(user *User, dados *db.Dados) *db.Dados {
	if user.Email != "" {
		(*dados)["email"] = user.Email
	}
	if user.Role != "" {
		(*dados)["role"] = user.Role
	}
	if user.Password != "" {
		(*dados)["password"] = user.Password
	}
	return dados
}

//GetUserByID Get an user by ID
func GetUserByID(u int) (*User, error) {

	var dados db.Dados
	dados = make(db.Dados)
	dados["id"] = u
	val, err := getUser(dados)
	if err != nil {
		return nil, err
	}
	account := val.(*User)
	if account.Email == "" { //User not found!
		return nil, sql.ErrNoRows
	}
	account.Password = ""
	return account, nil
}

//GetUserByEmail Get an user by email
func GetUserByEmail(email string, password bool) (*User, error) {

	var dados db.Dados
	dados = make(db.Dados)
	dados["email"] = email
	val, err := getUser(dados)
	if err != nil {
		return nil, err
	}
	account := val.(*User)
	if account.Email == "" { //User not found!
		return nil, sql.ErrNoRows
	}

	if !password {
		account.Password = ""
	}
	return account, nil
}

func getUser(dados db.Dados) (interface{}, error) {
	val, err := db.SelectOne(nil, "getuser", &User{}, dados)
	if err != nil {
		return nil, err
	}
	return val, nil
}
