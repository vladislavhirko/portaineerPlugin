package database

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

//Contains key-pair value
//Also contains container - chat (kv)
type LevelDB struct {
	Path string
	DB   *leveldb.DB
}

func NewLevelDB(path string) (*LevelDB, error) {
	ldb := &LevelDB{
		Path: path,
		DB:   nil,
	}
	err := ldb.Open()
	if err != nil {
		return nil, err
	}
	return ldb, nil
}

//Открывает соединение с базой данных
func (ldb *LevelDB) Open() error {
	var err error
	ldb.DB, err = leveldb.OpenFile(ldb.Path, nil)
	return err
}

//Добавляет или обновляет жлемент в базе данных
func (ldb *LevelDB) Put(key, value string) error {
	err := ldb.DB.Put([]byte(key), []byte(value), nil)
	log.Println("\nContainer: ", key, "\nChat: ", value)
	return err
}

//Получаем все значения  с базы данных
func (ldb *LevelDB) GetAll() map[string]string {
	data := map[string]string{}
	iterator := ldb.DB.NewIterator(nil, nil)
	for iterator.Next() {
		data[string(iterator.Key())] = string(iterator.Value())
	}
	iterator.Release()
	return data
}

// Возвращает валуе по ключу
func (ldb *LevelDB) Get(key string) (string, error) {
	value, err := ldb.DB.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

// удаляет запись по ключу
func (ldb *LevelDB) Delete(key string) error {
	err := ldb.DB.Delete([]byte(key), nil)
	return err
}
