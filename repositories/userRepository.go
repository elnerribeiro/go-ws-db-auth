package repositories

import (
	"database/sql"

	db "github.com/elnerribeiro/go-mustache-db"
	u "github.com/elnerribeiro/go-ws-db-auth/utils"
)

//ListUsers Lists all users
func (user *User) ListUsers() ([]User, error) {
	var dados db.Dados
	dados = make(db.Dados)
	val, err := db.SelectAll(nil, "getuser", &[]User{}, dados)
	if err != nil {
		u.Logger.Error("[ListUsers] Error listing users: %s", err)
		return nil, err
	}
	account := val.(*[]User)
	if account != nil {
		for i, u := range *account {
			u.Password = ""
			(*account)[i] = u
		}
		return *account, nil
	}
	u.Logger.Error("[ListUsers] No user found.")
	return *account, nil
}

//Upsert Inserts or updates an user
func (user *User) Upsert(tx *db.Transacao) (*User, error) {
	var dados db.Dados
	dados = make(db.Dados)
	user.UserToDados(&dados)
	if user.ID == 0 {
		res, err := db.InsertReturningPostgres(tx, "usuario", dados, "id", &User{})
		if err != nil {
			u.Logger.Error("[Upsert] Error inserting user: %s", err)
			return nil, err
		}
		newID := res.(*User)
		user.ID = newID.ID
		user.Password = ""
		return user, nil
	}

	var filter db.Dados
	filter = make(db.Dados)
	filter["id"] = user.ID
	_, err2 := db.Update(tx, "usuario", dados, filter)
	if err2 != nil {
		u.Logger.Error("[Upsert] Error updating user: %s", err2)
		return nil, err2
	}
	user.Password = ""
	return user, nil
}

//Delete Deletes an user
func (user *User) Delete(tx *db.Transacao) error {
	var filter db.Dados
	filter = make(db.Dados)
	filter["id"] = user.ID
	_, err := db.Delete(tx, "usuario", filter)
	return err
}

//UserToDados Fills a map containing the struct User
func (user *User) UserToDados(dados *db.Dados) *db.Dados {
	if user.ID != 0 {
		(*dados)["id"] = user.ID
	}
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
func (user *User) GetUserByID() (*User, error) {

	var dados db.Dados
	dados = make(db.Dados)
	dados["id"] = user.ID
	val, err := getUser(dados)
	if err != nil {
		u.Logger.Error("[GetUserByID] Error retrieving user: %s", err)
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
func (user *User) GetUserByEmail(password bool) (*User, error) {

	var dados db.Dados
	dados = make(db.Dados)
	dados["email"] = user.Email
	val, err := getUser(dados)
	if err != nil {
		u.Logger.Error("[GetUserByEmail] Error retrieving user: %s", err)
		return nil, err
	}
	account := val.(*User)
	if account.Email == "" { //User not found!
		u.Logger.Error("[GetUserByEmail] Error retrieving user: %s", err)
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
		u.Logger.Error("[getUser] Error selecting user: %s", err)
		return nil, err
	}
	return val, nil
}
