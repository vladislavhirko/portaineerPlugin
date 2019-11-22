package database

import (
	"github.com/syndtr/goleveldb/leveldb"
	log "github.com/sirupsen/logrus"
	"os"
)

type LevelDB struct {
	PathContainerChat string
	DBContainerChat   *DB
	PathAccounts string
	DBAccounts   *DB
}

//Open connection with database
func (ldb *LevelDB) Open() error {
	var err error

	homeDir, err := os.UserHomeDir()
	if err != nil{
		return err
	}
	ldb.DBContainerChat.DB, err = leveldb.OpenFile(homeDir + ldb.PathContainerChat, nil)
	if err != nil{
		return err
	}
	ldb.DBAccounts.DB, err = leveldb.OpenFile(homeDir + ldb.PathAccounts, nil)
	if err != nil{
		return err
	}
	return err
}

//Setup paths to storage, setup default logs
func NewLevelDB(path string) (*LevelDB, error) {
	ldb := &LevelDB{
		PathContainerChat: path + "/container-chat",
		DBContainerChat:   &DB{
			LogContext: log.WithFields(log.Fields{
				"Module": "Database",
				"Table": "Container-chat",
			})},
		PathAccounts: path + "/accounts",
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