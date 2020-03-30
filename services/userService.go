package services

import (
	"database/sql"
	"time"

	"github.com/dgrijalva/jwt-go"
	db "github.com/elnerribeiro/go-mustache-db"
	repo "github.com/elnerribeiro/go-ws-db-auth/repositories"
	u "github.com/elnerribeiro/go-ws-db-auth/utils"
)

//ListUsers Lists all users
func ListUsers(user *repo.User) ([]repo.User, error) {
	return user.ListUsers()
}

//Login Authenticates an user
func Login(user *repo.User, password string) map[string]interface{} {
	account, err := user.GetUserByEmail(true)
	if err != nil {
		if err == sql.ErrNoRows {
			u.Logger.Error("[Login] Email not found.")
			return u.Message(false, "Email address not found")
		}
		u.Logger.Error("[Login] Connection error. Please retry: %s", err)
		return u.Message(false, "Connection error. Please retry")
	}
	if account.Password != password { //Password does not match!
		u.Logger.Error("[Login] Invalid login credentials. Please try again: %s != %s", account.Password, password)
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	expirationTime := time.Now().Add(12 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &repo.Token{
		UserID: account.ID,
		Role:   account.Role,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("JWTpassword123@"))
	account.Token = tokenString //Store the token in the response

	u.Logger.Info("[Login] User Logged In: %d", account.ID)

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

//GetUserByID Gets an user by ID
func GetUserByID(user *repo.User) (*repo.User, error) {

	val, err := user.GetUserByID()
	if err != nil {
		u.Logger.Error("[GetUserByID] Error : %s", err)
		return nil, err
	}
	return val, nil
}

//Upsert Inserts or updates an user
func Upsert(user *repo.User) (*repo.User, error) {
	tx, err := db.GetTransaction()
	defer db.Rollback(tx)
	if err != nil {
		u.Logger.Error("[Upsert] Error starting transaction: %s", err)
		return nil, err
	}
	val, err := user.Upsert(tx)
	if err != nil {
		u.Logger.Error("[Upsert] Error executing upsert: %s", err)
		return nil, err
	}
	db.Commit(tx)
	return val, nil
}

//Delete Deletes an user
func Delete(user *repo.User) error {
	tx, err := db.GetTransaction()
	defer db.Rollback(tx)
	if err != nil {
		u.Logger.Error("[Delete] Error starting transaction: %s", err)
		return err
	}
	if err2 := user.Delete(tx); err2 != nil {
		u.Logger.Error("[Delete] Error executing delete: %s", err)
		return err2
	}
	db.Commit(tx)
	return nil
}
