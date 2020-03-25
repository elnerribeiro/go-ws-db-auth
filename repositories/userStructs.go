package models

import (
	"github.com/dgrijalva/jwt-go"
)

//Token JWT Token
type Token struct {
	UserID int
	Role   string
	jwt.StandardClaims
}

//User table usuario on database
type User struct {
	ID       int    `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
	Role     string `json:"role,omitempty"`
}

//ContextKey Key to use on a context
type ContextKey string
