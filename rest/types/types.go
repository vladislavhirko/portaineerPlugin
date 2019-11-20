package types

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

type KeyValue struct{
	Key string `json:"key,name"`
	Value string `json:"value,pass"`
}

type JWTAuth struct{
	SigningKey []byte
	Claims MyClaims
}

type MyClaims struct {
	Name string `json:"name"`
	Admin bool `json:"admin"`
	jwt.StandardClaims
}

type ErrorGroup struct {
	StatusCode int
	Error error
	ResponseWriter http.ResponseWriter
}