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
func ListUsers() ([]repo.User, error) {
	return repo.ListUsers()
}

//Login Authenticates an user
func Login(email, password string) map[string]interface{} {
	account, err := repo.GetUserByEmail(email, true)
	if err != nil {
		if err == sql.ErrNoRows {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}
	if account.Password != password { //Password does not match!
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

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

//GetUserByID Gets an user by ID
func GetUserByID(u int) (*repo.User, error) {

	val, err := repo.GetUserByID(u)
	if err != nil {
		return nil, err
	}
	return val, nil
}

//Upsert Inserts or updates an user
func Upsert(user *repo.User) (*repo.User, error) {
	tx, err := db.GetTransaction()
	defer db.Rollback(tx)
	if err != nil {
		return nil, err
	}
	val, err := repo.Upsert(tx, user)
	if err != nil {
		return nil, err
	}
	db.Commit(tx)
	return val, nil
}

//Delete Deletes an user
func Delete(uid int) error {
	tx, err := db.GetTransaction()
	defer db.Rollback(tx)
	if err != nil {
		return err
	}
	err2 := repo.Delete(tx, uid)
	if err2 != nil {
		return err2
	}
	db.Commit(tx)
	return nil
}
