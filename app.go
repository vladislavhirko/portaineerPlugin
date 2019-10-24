package main

import (
	"fmt"
	"log"
	"github.com/vladislavhirko/portaineerPlugin/portainerAPI"
	"github.com/vladislavhirko/portaineerPlugin/portainerAPI/types"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	stopedContainerChan := make(chan types.Containers)
	go DockerChecker(stopedContainerChan, wg)
	go Sender(stopedContainerChan, wg)
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
		time.Sleep(time.Second * 10)
	}
}

// Функция которая будет рассылать список упавших контейнеров в метермост
func Sender(stopedContainerChan chan types.Containers, wg sync.WaitGroup) {
	defer wg.Done()
	for {
		fmt.Println("ALAAAAAAH AKBAR", <-stopedContainerChan)
	}
}
