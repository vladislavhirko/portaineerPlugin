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

//Chanel by which sending message about stopped containers
var stopedContainerChan = make(chan types.Containers)

func main() {
	log.Info("Run")
	usr, _ := user.Current()
	configPath := flag.String("config_path", usr.HomeDir + "/.portaineerPlugin/config.toml", "Path to file config")
	flag.Parse()
	systemConfig, err := config.GetConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	SetLogLevel(systemConfig.LogLevel.Level)
	Starter(*systemConfig)
}

//Create objects for working with leveldb, mattermoost, portainer, server
//Run three goroutine: sender (send message), checker (check containers state), server (api)
func Starter(config config.Config) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	levelDB := LDBStart(config.LevelDB)
	mClient := MattermostStart(levelDB, config.MClient)
	pClient := PortainerStart(config.PClient)
	server := rest.NewServer(config.API, levelDB, pClient)

	go Sender(mClient) //Function that listens to the channel and if something gets there, it sends to the mattermost
	go DockerChecker(pClient)
	go server.StartServer()
	wg.Wait()
}

//Function which opens connection with leveldb
func LDBStart(config config.Level) database.LevelDB {
	ldb, err := database.NewLevelDB(config.Path)
	if err != nil {
		log.Fatal(err)
	}
	return *ldb
}

//Starts work with portainer
func PortainerStart(config config.Portainer) *portainer.ClientPortaineer {
	pClient := portainer.NewPorteinerClient(config.Address, config.Port, config.CheckInterval, config.LogsAmount)
	err := pClient.Auth(config.Login, config.Password)
	if err != nil {
		log.Fatal(err)
	}
	return pClient
}

//Start work with mattermost, and login client
func MattermostStart(ldb database.LevelDB, config config.Mattermost) mattermost.MattermostClient {
	mClient := mattermost.NewMattermostClient(ldb, config.Address, config.Port)
	err := mClient.Auth(config.Email, config.Password)
	if err != nil {
		log.Fatal(err)
	}
	return mClient
}

//Function which setup pauses for timestamps every 'X' second and calls function wich
//takes containers list, after that it trieing to find stoped containers, adds to list of stoped, creats new list
//for saving these (list of stoped will clean after that)
//and sends this list by chanel to Sender() function
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

//Which which sends list of stopped containers to mattermost
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

func SetLogLevel(level string){
	logLevel, err := log.ParseLevel(level)
	if err != nil{
		log.Fatal(err)
	}
	log.SetLevel(logLevel)
}