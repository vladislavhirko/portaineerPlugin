package rest

import (
	"encoding/json"
	"fmt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/vladislavhirko/portaineerPlugin/config"
	"github.com/vladislavhirko/portaineerPlugin/database"
	"github.com/vladislavhirko/portaineerPlugin/portainer"
	"github.com/vladislavhirko/portaineerPlugin/rest/types"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	//"time"
)

type Server struct{
	Config config.API
	LDB database.LevelDB
	pClient *portainer.ClientPortaineer
}

func NewServer(config config.API, ldb database.LevelDB, pClient *portainer.ClientPortaineer) Server{
	return Server{
		Config:  config,
		LDB:     ldb,
		pClient: pClient,
	}
}

func (server Server) StartServer(){
	headersOk := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	r := mux.NewRouter()
	r.HandleFunc("/pairs", jwtMiddleware.Handler(server.AddPairHandler()).ServeHTTP).Methods("POST")
	r.HandleFunc("/pairs", jwtMiddleware.Handler(server.DeletePairHandler()).ServeHTTP).Methods("DELETE")
	r.HandleFunc("/pairs", jwtMiddleware.Handler(server.GetPairsHandler()).ServeHTTP).Methods("GET")
	r.HandleFunc("/containers", jwtMiddleware.Handler(server.GetContainersHandler()).ServeHTTP).Methods("GET")
	r.HandleFunc("/get_token", http.HandlerFunc(GetTokenHandler)).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":" + server.Config.Port, handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}

// Добавляет в базу данных новый ключ-значение
func (server Server) AddPairHandler() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request){

		kv := types.KeyValue{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil{
			http.Error(w, "something wrong", 400)
		}
		err = json.Unmarshal(body, &kv)
		if err != nil{
			http.Error(w, "uncorrect json format", 400)
		}
		err = server.LDB.Put(kv.Key, kv.Value)
		if err != nil{
			http.Error(w, "some troubles with database", 400)
		}
		w.Write([]byte("Ok"))
	}
}

func (server Server) GetContainersHandler() http.HandlerFunc {
	return func (w http.ResponseWriter, r * http.Request){
		fmt.Println(r.Header)
		err := server.LDB.Put("/crazy_volhard", "crazy")
		if err != nil{
			http.Error(w, err.Error(), 400)
		}
		containersJSON, _ := json.Marshal(server.pClient.CurrentContainers)
		w.Write(containersJSON)
	}
}

//возвращает список ключ-значение (контейнер - чат)
func (server Server) GetPairsHandler() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
		pairs := server.LDB.GetAll()
		pairsJSON, _ := json.Marshal(pairs)
		w.Write(pairsJSON)
	}
}

func (server Server) DeletePairHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil{
			http.Error(w, err.Error(), 400)
		}
		err = server.LDB.Delete(string(body))
		if err != nil{
			http.Error(w, err.Error(), 400)
		}
		w.Write([]byte("OK"))
	}
}


//-----------------------------------------------------------//

