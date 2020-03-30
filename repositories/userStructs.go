package repositories

import (
	"github.com/dgrijalva/jwt-go"
	db "github.com/elnerribeiro/go-mustache-db"
)

//Token JWT Token
type Token struct {
	UserID int
	Role   string
	jwt.StandardClaims
}

//ContextKey Key to use on a context
type ContextKey string

//UserRepository Repository for table usuario
type UserRepository interface {
	ListUsers() ([]User, error)
	Upsert(tx *db.Transacao) (*User, error)
	Delete(tx *db.Transacao)
	UserToDados(dados *db.Dados) *db.Dados
	GetUserByID() (*User, error)
	GetUserByEmail(password bool) (*User, error)
}

//User table usuario on database
type User struct {
	ID       int    `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
	Role     string `json:"role,omitempty"`
}
