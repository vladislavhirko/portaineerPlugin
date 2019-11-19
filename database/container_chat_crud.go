package database

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

//Contains key-pair value
//Also contains container - chat (kv)
type LevelDB struct {
	PathContainerChat string
	DBContainerChat   *leveldb.DB

	PathAccounts string
	DBAccounts   *leveldb.DB
}

func NewLevelDB(path string) (*LevelDB, error) {
	ldb := &LevelDB{
		PathContainerChat: path + "storage",
		DBContainerChat:   nil,
		PathAccounts: path + "accounts",
		DBAccounts: nil,
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
	ldb.DBContainerChat, err = leveldb.OpenFile(ldb.PathContainerChat, nil)
	return err
}

//Добавляет или обновляет жлемент в базе данных
func (ldb *LevelDB) Put(key, value string) error {
	err := ldb.DBContainerChat.Put([]byte(key), []byte(value), nil)
	log.Println("\nContainer: ", key, "\nChat: ", value)
	return err
}

//Получаем все значения  с базы данных
func (ldb *LevelDB) GetAll() map[string]string {
	data := map[string]string{}
	iterator := ldb.DBContainerChat.NewIterator(nil, nil)
	for iterator.Next() {
		data[string(iterator.Key())] = string(iterator.Value())
	}
	iterator.Release()
	return data
}

// Возвращает валуе по ключу
func (ldb *LevelDB) Get(key string) (string, error) {
	value, err := ldb.DBContainerChat.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

// удаляет запись по ключу
func (ldb *LevelDB) Delete(key string) error {
	err := ldb.DBContainerChat.Delete([]byte(key), nil)
	return err
}
