package rest

import (
	"fmt"
	"github.com/vladislavhirko/portaineerPlugin/database"
	"sync"
	"time"
)

//TODO тут будем запускать рест
func RunServer(ldb database.LevelDB, wg sync.WaitGroup){
	defer wg.Done()
	ldb.Put("FirstValue", "SecondValue")
	fmt.Println(ldb.GetAll())
	for{
		time.Sleep(time.Second * 10)
	}
}
