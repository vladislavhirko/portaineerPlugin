package types

import "github.com/dgrijalva/jwt-go"

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