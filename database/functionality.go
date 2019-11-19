package database

import (
	"github.com/syndtr/goleveldb/leveldb"
	log "github.com/sirupsen/logrus"

)

//Contains key-pair value
//Also contains container - chat (kv)

type DB struct{
	DB *leveldb.DB
	LogContext *log.Entry
}

//Добавляет или обновляет жлемент в базе данных
func (ccdb *DB) Put(key, value string) error {
	err := ccdb.DB.Put([]byte(key), []byte(value), nil)
	ccdb.LogContext.Info("Put: Key: \"", key, "\" -> Value: \"", value + "\"")
	return err
}

//Получаем все значения  с базы данных
func (ccdb *DB) GetAll() map[string]string {
	data := map[string]string{}
	iterator := ccdb.DB.NewIterator(nil, nil)
	for iterator.Next() {
		data[string(iterator.Key())] = string(iterator.Value())
	}
	iterator.Release()
	ccdb.LogContext.Info("Get all")
	return data
}

// Возвращает валуе по ключу
func (ccdb *DB) Get(key string) (string, error) {
	value, err := ccdb.DB.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	ccdb.LogContext.Info("Get: ", key)
	return string(value), nil
}

// удаляет запись по ключу
func (ccdb *DB) Delete(key string) error {
	err := ccdb.DB.Delete([]byte(key), nil)
	ccdb.LogContext.Info("Delete: ", key)
	return err
}
