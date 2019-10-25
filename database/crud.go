package database

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	)

type LevelDB struct{
	DB *leveldb.DB
}

func NewLevelDB() LevelDB{
	return LevelDB{
		DB:nil,
	}
}

//Открывает соединение с базой данных
func (ldb *LevelDB) Open(){
	var err error
	ldb.DB, err = leveldb.OpenFile("storage/storage", nil)
	if err != nil{
		fmt.Println(err)
	}
}

//Добавляет или обновляет жлемент в базе данных
func (ldb *LevelDB) Put(key, value string){
	err := ldb.DB.Put([]byte(key), []byte(value), nil)
	if err != nil{
		fmt.Println(err)
	}
}

//Получаем все значения  с базы данных
func (ldb *LevelDB) GetAll() (map[string]string){
	data := map[string]string{}
	iterator := ldb.DB.NewIterator(nil, nil)
	for iterator.Next(){
		data[string(iterator.Key())] = string(iterator.Value())
	}
	iterator.Release()
	return data
}

//TODO сделать метод удаления элдементов с базы данных