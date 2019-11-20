package main

import (
	"flag"
	"github.com/vladislavhirko/portaineerPlugin/config"
	"github.com/vladislavhirko/portaineerPlugin/database"
	"github.com/vladislavhirko/portaineerPlugin/mattermost"
	"github.com/vladislavhirko/portaineerPlugin/portainer"
	"github.com/vladislavhirko/portaineerPlugin/portainer/types"
	"github.com/vladislavhirko/portaineerPlugin/rest"
	"os/user"
	"sync"
	"time"
	log "github.com/sirupsen/logrus"
)

var stopedContainerChan = make(chan types.Containers) //Канал по которому передаются сообщение при падении контейнера

func main() {
	log.SetLevel(log.TraceLevel)
	log.Info("Run")
	usr, _ := user.Current()
	configPath := flag.String("config_path", usr.HomeDir + "/.portaineerPlugin/config.toml", "Path to file config")
	flag.Parse()
	systemConfig, err := config.GetConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	Starter(*systemConfig)
}

func Starter(config config.Config) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	levelDB := LDBStart(config.LevelDB)
	mClient := MattermostStart(levelDB, config.MClient)
	pClient := PortainerStart(config.PClient)
	server := rest.NewServer(config.API, levelDB, pClient)

	go Sender(mClient) //Функция слушающая канал и в случае попадания туда чего либо отправляющая в меттермост
	go DockerChecker(pClient)
	go server.StartServer() //В будущем тут будет рест
	wg.Wait()
}

//Функция для открытыя соедниенения с левелдб
func LDBStart(config config.Level) database.LevelDB {
	ldb, err := database.NewLevelDB(config.Path)
	if err != nil {
		log.Fatal(err)
	}
	return *ldb
}

//Начало работы с потейнером
func PortainerStart(config config.Portainer) *portainer.ClientPortaineer {
	pClient := portainer.NewPorteinerClient(config.Address, config.Port, config.CheckInterval, config.LogsAmount)
	err := pClient.Auth(config.Login, config.Password)
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
func DockerChecker(pClient *portainer.ClientPortaineer) {
	for {
		go func() {
			err := pClient.GetContainerrList()
			if err != nil {
				log.Error(err)
			}
			pClient.FinedDropedContainers()
			dropedContainers := types.Containers{}
			if len(pClient.StopedContainers) != 0 {
				dropedContainers, err = pClient.GetDropedContainer()
				//fmt.Println("LEEEEEENGTH", dropedContainers)
				if err != nil {
					log.Error(err)
				}
			}
			if len(dropedContainers) > 0 {
				stopedContainerChan <- dropedContainers
			}
		}()
		time.Sleep(time.Second * pClient.CheckInterval)
	}
}

// Функция которая рассылает список упавших контейнеров в метермост, пока что во все каналы
//Закоменчено потому что мы не можем сыпать логи в метермост  пока что
func Sender(mClient mattermost.MattermostClient) {
	for {
		dropedContainers := <-stopedContainerChan
		//fmt.Println("ALAAAAAAH AKBAR", dropedContainers)
		err := mClient.GetallChanels()
		if err != nil {
			log.Error(err)
		}
		err = mClient.SendMessage(dropedContainers, "")
		if err != nil {
			log.Error(err)
		}
	}
}
