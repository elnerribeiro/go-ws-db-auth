package models

import (
	"database/sql"
	"time"

	"github.com/dgrijalva/jwt-go"
	db "github.com/elnerribeiro/go-mustache-db"
	u "github.com/elnerribeiro/go-ws-db-auth/utils"
)

//Token e um token JWT
type Token struct {
	UserID int
	jwt.StandardClaims
}

//User e um usuario no banco de dados
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

//ListUsers retorna todos os usuarios
func ListUsers() ([]User, error) {
	var dados db.Dados
	dados = make(db.Dados)
	val, err := db.SelectAll(nil, "buscarusuario", &[]User{}, dados)
	if err != nil {
		return nil, err
	}
	account := val.(*[]User)
	if account != nil {
		for _, u := range *account {
			u.Password = ""
		}
	}
	return *account, nil
}

//Login efetua o login do usuario
func Login(email, password string) map[string]interface{} {
	var dados db.Dados
	dados = make(db.Dados)
	dados["email"] = email
	var err error
	var val interface{}
	val, err = db.SelectOne(nil, "buscarusuario", &User{}, dados)
	if err != nil {
		if err == sql.ErrNoRows {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}
	account := val.(*User)
	if account.Password != password { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	expirationTime := time.Now().Add(12 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Token{
		UserID: account.ID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("softboxgopwd12@"))
	account.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

//GetUser buscar um usuario por id
func GetUser(u uint) *User {

	var dados db.Dados
	dados = make(db.Dados)
	dados["id"] = u
	var err error
	var val interface{}
	val, err = db.SelectOne(nil, "buscarusuario", &User{}, dados)
	if err != nil {
		return nil
	}
	account := val.(*User)
	if account.Email == "" { //User not found!
		return nil
	}
	account.Password = ""
	return account
}
