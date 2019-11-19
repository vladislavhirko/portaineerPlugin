package database

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

//Contains key-pair value
//Also contains container - chat (kv)

type DB struct{
	DB *leveldb.DB
}

//Добавляет или обновляет жлемент в базе данных
func (ccdb *DB) Put(key, value string) error {
	err := ccdb.DB.Put([]byte(key), []byte(value), nil)
	log.Println("\nContainer: ", key, "\nChat: ", value)
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
	return data
}

// Возвращает валуе по ключу
func (ccdb *DB) Get(key string) (string, error) {
	value, err := ccdb.DB.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

// удаляет запись по ключу
func (ccdb *DB) Delete(key string) error {
	err := ccdb.DB.Delete([]byte(key), nil)
	return err
}
