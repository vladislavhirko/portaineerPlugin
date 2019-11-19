package database

import (
	"github.com/syndtr/goleveldb/leveldb"
	log "github.com/sirupsen/logrus"
)

type LevelDB struct {
	PathContainerChat string
	DBContainerChat   *DB
	PathAccounts string
	DBAccounts   *DB
}

//Открывает соединение с базой данных
func (ldb *LevelDB) Open() error {
	var err error
	ldb.DBContainerChat.DB, err = leveldb.OpenFile(ldb.PathContainerChat, nil)
	if err != nil{
		return err
	}
	ldb.DBAccounts.DB, err = leveldb.OpenFile(ldb.PathAccounts, nil)
	if err != nil{
		return err
	}
	return err
}

func NewLevelDB(path string) (*LevelDB, error) {
	ldb := &LevelDB{
		PathContainerChat: path + "Container-chat",
		DBContainerChat:   &DB{
			LogContext: log.WithFields(log.Fields{
				"Module": "Database",
				"Table": "Container-chat",
			})},
		PathAccounts: path + "Accounts",
		DBAccounts: &DB{LogContext: log.WithFields(log.Fields{
			"Module": "Database",
			"Table": "Accounts",
		})},
	}
	err := ldb.Open()
	if err != nil {
		return nil, err
	}
	return ldb, nil
}