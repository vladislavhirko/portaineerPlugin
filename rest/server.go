package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/vladislavhirko/portaineerPlugin/config"
	"github.com/vladislavhirko/portaineerPlugin/database"
	"github.com/vladislavhirko/portaineerPlugin/rest/types"
	"io/ioutil"
	"log"
	"net/http"
)

func StartServer(ldb database.LevelDB, api config.API){
	r := mux.NewRouter()
	r.HandleFunc("/pairs", AddPairHandler(ldb)).Methods("POST")
	r.HandleFunc("/pairs", GetPairsHandler(ldb)).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":" + api.Port, nil))
}

// Добавляет в базу данных новый ключ-значение
func AddPairHandler(ldb database.LevelDB) func(w http.ResponseWriter, r *http.Request){
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
		err = ldb.Put(kv.Key, kv.Value)
		if err != nil{
			http.Error(w, "some troubles with database", 400)
		}
		w.Write([]byte("Ok"))
	}
}

//возвращает список ключ-значение (контейнер - чат)
func GetPairsHandler(ldb database.LevelDB) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request) {
		pairs := ldb.GetAll()
		pairsJSON, _ := json.Marshal(pairs)
		w.Write(pairsJSON)
	}
}