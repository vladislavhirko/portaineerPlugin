package main

import (
	"fmt"
	"log"
	"github.com/vladislavhirko/portaineerPlugin/database"
	"github.com/vladislavhirko/portaineerPlugin/mattermost"
	"github.com/vladislavhirko/portaineerPlugin/portainerAPI"
	"github.com/vladislavhirko/portaineerPlugin/portainerAPI/types"
	"github.com/vladislavhirko/portaineerPlugin/rest"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(3)
	stopedContainerChan := make(chan types.Containers)
	go DockerChecker(stopedContainerChan, wg) //Раз в Х сеунд делает запросы в апи докера
	go Sender(stopedContainerChan, wg) //Функция слушающая канал и в случае попадания туда чего либо отправляющая в меттермост
	ldb := database.NewLevelDB()
	ldb.Open() //открывает конекшн с базой данных
	go rest.RunServer(ldb, wg) //В будущем тут будет рест
	wg.Wait()
}

//Функция где устанавливает время для таймстемпов и и раз в Х секунд вызываются функции которые
// тянут список контейнеров, далее ищет те которые упали
// далее добавялют в список упавших, создает новый список в котором они хранятся (список упавших после этого чистится)
// затем отправляет по каналу в функию Sender()
func DockerChecker(stopedContainerChan chan types.Containers, wg sync.WaitGroup) {
	defer wg.Done()
	pClient := portainerAPI.PClientNew()
	err := pClient.Auth()
	if err != nil {
		log.Fatal(err)
	}
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
// TODO рассылка в каналы по названию
func Sender(stopedContainerChan chan types.Containers, wg sync.WaitGroup) {
	defer wg.Done()
	//mClient := MattermostStart()
	for {
		dropedContainers := <-stopedContainerChan
		fmt.Println("ALAAAAAAH AKBAR", dropedContainers)
		//err := mClient.GetallChanels()
		//if err != nil{
		//	fmt.Println(err)
		//}
		//err = mClient.SendMessage(dropedContainers, "")
		//if err != nil{
		//	fmt.Println(err)
		//}
	}
}

// Создает клиента для работы с меттермостом, логинит пользователя
func MattermostStart() mattermost.MattermostClient{
	mClient := mattermost.NewMattermostClient()
	mClient.CreateClient("http://192.168.88.62:8065")
	err := mClient.Login("qwerty@gmail.com", "qwerty")
	if err != nil{
		fmt.Println(err)
	}

	return mClient
}