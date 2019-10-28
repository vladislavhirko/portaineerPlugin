package main

import (
	"fmt"
	"github.com/vladislavhirko/portaineerPlugin/config"
	"github.com/vladislavhirko/portaineerPlugin/database"
	"github.com/vladislavhirko/portaineerPlugin/mattermost"
	"github.com/vladislavhirko/portaineerPlugin/portainerAPI"
	"github.com/vladislavhirko/portaineerPlugin/portainerAPI/types"
	"github.com/vladislavhirko/portaineerPlugin/rest"
	"log"
	"sync"
	"time"
)

var stopedContainerChan = make(chan types.Containers) //Канал по которому передаются сообщение при падении контейнера
var wg = sync.WaitGroup{}

func main() {
	systemConfig := config.GetConfig()
	Starter(systemConfig)
}

func Starter(config config.Config) {
	wg.Add(3)
	levelDB := LDBStart(config.LevelDB)
	mClient := MattermostStart(levelDB, config.MClient)
	pClient := PortainerStart(config.PClient)

	go Sender(mClient) //Функция слушающая канал и в случае попадания туда чего либо отправляющая в меттермост
	go DockerChecker(pClient)
	go rest.RunServer(levelDB, wg) //В будущем тут будет рест
	wg.Wait()
}

func LDBStart(config config.Level) database.LevelDB {
	ldb, err := database.NewLevelDB(config.Path)
	if err != nil {
		log.Fatal(err)
	}
	return *ldb
}

func PortainerStart(config config.Portainer) portainerAPI.ClientPortaineer {
	pClient := portainerAPI.NewPorteinerClient(config.Address, config.Port)
	err := pClient.Auth(config.Email, config.Password)
	if err != nil {
		log.Fatal(err)
	}
	return pClient
}

// Создает клиента для работы с меттермостом, логинит пользователя
func MattermostStart(ldb database.LevelDB, config config.Mattermost) mattermost.MattermostClient {
	mClient := mattermost.NewMattermostClient(ldb, config.Address, config.Port)
	err := mClient.Auth(config.Email, config.Password)
	if err != nil {
		log.Fatal(err)
	}
	return mClient
}

//Функция где устанавливает время для таймстемпов и и раз в Х секунд вызываются функции которые
// тянут список контейнеров, далее ищет те которые упали
// далее добавялют в список упавших, создает новый список в котором они хранятся (список упавших после этого чистится)
// затем отправляет по каналу в функию Sender()
func DockerChecker(pClient portainerAPI.ClientPortaineer) {
	defer wg.Done()
	for {
		err := pClient.GetContainerrList()
		if err != nil {
			log.Fatal(err)
		}
		pClient.StopedTrigger()
		dropedContainers := types.Containers{}
		if len(pClient.StopedContainers) != 0 {
			dropedContainers, err = pClient.GetDropedContainer()
			//fmt.Println("LEEEEEENGTH", dropedContainers)
			if err != nil {
				log.Fatal(err)
			}
		}
		if len(dropedContainers) > 0 {
			stopedContainerChan <- dropedContainers
		}
		time.Sleep(time.Second * 3)
	}
}

// Функция которая рассылает список упавших контейнеров в метермост, пока что во все каналы
//Закоменчено потому что мы не можем сыпать логи в метермост  пока что
func Sender(mClient mattermost.MattermostClient) {
	defer wg.Done()
	for {
		dropedContainers := <-stopedContainerChan
		//fmt.Println("ALAAAAAAH AKBAR", dropedContainers)
		err := mClient.GetallChanels()
		if err != nil {
			fmt.Println(err)
		}
		err = mClient.SendMessage(dropedContainers, "")
		if err != nil {
			fmt.Println(err)
		}
	}
}
