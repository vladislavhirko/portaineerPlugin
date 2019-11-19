package rest

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var mySigningKey = []byte("das3f12A32f32a33efA3E32F32f3e2FW32f32e")

type MyClaims struct {
	Name string `json:"name"`
	Admin bool `json:"admin"`
	jwt.StandardClaims
}

//Генерит и возвращает токен клиенту
func GetTokenHandler(w http.ResponseWriter, r *http.Request){
	expirationTime := time.Now().Add(168 * time.Hour)
	// Устанавливаем набор параметров для токена
	claims := MyClaims{
		"john_lenon",
		true,
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Подписываем токен нашим секретным ключем
	tokenString, _ := token.SignedString(mySigningKey)
	// Отдаем токен клиенту
	w.Write([]byte(tokenString))
}




//func JWTMiddlewear(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request){
//	return func(w http.ResponseWriter, r *http.Request){
//		jwtToken := r.Header.Get("Authentication")
//		token, _ := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
//			return mySigningKey, nil
//		})
//		if !token.Valid{
//			http.Error(w, "", 403)
//		}
//	}
//}

//var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
//	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
//		return mySigningKey, nil
//	},
//	SigningMethod: jwt.SigningMethodHS256,
//})

//func ParseClaim() {
//	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"
//
//	type MyClaims struct {
//		Name string `json:"name"`
//		Admin bool `json:"admin"`
//		jwt.StandardClaims
//	}
//
//	// sample token is expired.  override time so it parses as valid
//	at(time.Unix(0, 0), func() {
//		token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
//			return []byte("AllYourBase"), nil
//		})
//
//		if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
//			fmt.Printf("%v %v %v", claims.Name, claims.Admin,claims.StandardClaims.ExpiresAt)
//		} else {
//			fmt.Println(err)
//		}
//	})
//	// Output: bar 15000
//}
//
////Обновляет время жизни токена
//func at(t time.Time, f func()) {
//	jwt.TimeFunc = func() time.Time {
//		return t
//	}
//	f()
//	jwt.TimeFunc = time.Now
//}
//
