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
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Server struct{
	Config config.API
	LDB database.LevelDB
	pClient *portainer.ClientPortaineer
	JWTAuth types.JWTAuth
	Log *log.Entry
}

func NewServer(config config.API, ldb database.LevelDB, pClient *portainer.ClientPortaineer) Server{
	return Server{
		Config:  config,
		LDB:     ldb,
		pClient: pClient,
		JWTAuth: types.JWTAuth{
			SigningKey: []byte("das3f12A32f32a33efA3E32F32f3e2FW32f32e"),
			Claims:     types.MyClaims{},
		},
		Log:  log.WithFields(log.Fields{
			"Module": "Server",
		}),
	}
}

func (server Server) StartServer(){
	headersOk := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	//Function for checking token

	var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return server.JWTAuth.SigningKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	r := mux.NewRouter()
	r.HandleFunc("/pairs", server.Loger(jwtMiddleware.Handler(server.AddPairHandler())).ServeHTTP).Methods("POST")
	r.HandleFunc("/pairs/{container_name}", server.Loger(jwtMiddleware.Handler(server.DeletePairHandler())).ServeHTTP).Methods("DELETE")
	//r.HandleFunc("/pairs/{Container}", server.Loger(server.DeletePairHandler()).ServeHTTP).Methods("DELETE")
	r.HandleFunc("/pairs", server.Loger(jwtMiddleware.Handler(server.GetPairsHandler())).ServeHTTP).Methods("GET")
	//r.HandleFunc("/containers", Log(jwtMiddleware.Handler(server.GetContainersHandler())).ServeHTTP).Methods("GET")
	r.HandleFunc("/containers", server.GetContainersHandler()).Methods("GET")
	r.HandleFunc("/get_token", server.Loger(server.GetTokenHandler()).ServeHTTP).Methods("POST")
	http.Handle("/", r)
	server.Log.Info("Start")
	log.Fatal(http.ListenAndServe(":" + server.Config.Port, handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}

// Add pair of container-chat to database
func (server Server) AddPairHandler() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request){
		kv := types.KeyValue{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil{
			server.ErrorHandler(types.ErrorGroup{
				StatusCode:     400,
				Error:       err,
				ResponseWriter: w,
			})
			return
		}
		err = json.Unmarshal(body, &kv)
		if err != nil{
			server.ErrorHandler(types.ErrorGroup{
				StatusCode:     400,
				Error:       err,
				ResponseWriter: w,
			})
			return
		}
		err = server.LDB.DBContainerChat.Put(kv.Key, kv.Value)
		if err != nil{
			server.ErrorHandler(types.ErrorGroup{
				StatusCode:     400,
				Error:       err,
				ResponseWriter: w,
			})
			return
		}
		w.Write([]byte("Ok"))
	}
}

// Get containers list from database
func (server Server) GetContainersHandler() http.HandlerFunc {
	return func (w http.ResponseWriter, r * http.Request){
		err := server.LDB.DBContainerChat.Put("/crazy_volhard", "crazy")
		err = server.LDB.DBContainerChat.Put("/elated_brattain", "roman")
		if err != nil{
			server.ErrorHandler(types.ErrorGroup{
				StatusCode:     400,
				Error:       err,
				ResponseWriter: w,
			})
			return
		}
		err = server.LDB.DBAccounts.Put("admin", "adminadmin")
		if err != nil{
			server.ErrorHandler(types.ErrorGroup{
				StatusCode:     400,
				Error:       err,
				ResponseWriter: w,
			})
			return
		}
		containersJSON, err := json.Marshal(server.pClient.CurrentContainers)
		if err != nil{
			server.ErrorHandler(types.ErrorGroup{
				StatusCode:     400,
				Error:       err,
				ResponseWriter: w,
			})
			return
		}
		w.Write(containersJSON)
	}
}

//Get list of pair (container-chat) from database
func (server Server) GetPairsHandler() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
		pairs := server.LDB.DBContainerChat.GetAll()
		pairsJSON, err := json.Marshal(pairs)
		if err != nil{
			server.ErrorHandler(types.ErrorGroup{
				StatusCode:     400,
				Error:       err,
				ResponseWriter: w,
			})
			return
		}
		w.Write(pairsJSON)
	}
}

//Deletes pair (container-chat) from database
func (server Server) DeletePairHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//body, err := ioutil.ReadAll(r.Body)
		//if err != nil{
		//	server.ErrorHandler(types.ErrorGroup{
		//		StatusCode:     400,
		//		Error:       err,
		//		ResponseWriter: w,
		//	})
		//	return
		//}
		params := mux.Vars(r)
		container := params["container_name"]
		fmt.Println(container)
		err := server.LDB.DBContainerChat.Delete(container)
		if err != nil{
			server.ErrorHandler(types.ErrorGroup{
				StatusCode:     403,
				Error:       err,
				ResponseWriter: w,
			})
			return
		}
		w.Write([]byte("OK"))
	}
}


//-----------------------------------------------------------//

//Interlayer - logger
func (server Server) Loger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.Log.Info(r.URL, " / ", r.Method)
		h.ServeHTTP(w, r) //Calls handler h
	})
}

func (server Server) ErrorHandler(errGroup types.ErrorGroup){
	server.Log.Error(errGroup.Error)
	http.Error(errGroup.ResponseWriter, "", errGroup.StatusCode)
}