package portainerAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/vladislavhirko/portaineerPlugin/portainerAPI/types"

)

//Структура клиента для портейнера
type ClientPortaineer struct {
	Jwt               string `json:"jwt"`
	Username          string
	Password          string
	CurrentContainers types.Containers
	LastContainers    types.Containers
	StopedContainers  types.Containers
}

//Создание нового клиента портейнера, выделяет память для слайсов контейнеров
func PClientNew() ClientPortaineer {
	return ClientPortaineer{
		Jwt:               "",
		Username:          "",
		Password:          "",
		CurrentContainers: make(types.Containers, 0),
		LastContainers:    make(types.Containers, 0),
		StopedContainers:  make(types.Containers, 0),
	}
}

//Аутентификация по логину и паролю, устанавливает JWT токен,
//должен передаваться во все последующие запросы на портейнер в заголовке
func (pClient *ClientPortaineer) Auth() error {
	pClient.Username = "admin"
	pClient.Password = "adminadmin"

	authObjJSON, err := json.Marshal(&pClient)
	if err != nil {
		return err
	}
	r := bytes.NewReader(authObjJSON)
	resp, err := http.Post("http://localhost:9000/api/auth", "application/json", r)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &pClient)
	if err != nil {
		return err
	}
	return nil
}

// Получает список контенеров и устанавливает их в переменную структуры
// так перед этим сохраняются все контенера работающие до обновления
// (список контейнеров которые работали Х секунд назад и список который работает сейчас)
func (pClient *ClientPortaineer) GetContainerrList() error {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET", "http://localhost:9000/api/endpoints/1/docker/containers/json?all=1", nil,
	)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", pClient.Jwt)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("HU((((R'")
	pClient.LastContainers = make(types.Containers, 0)
	pClient.LastContainers = append(pClient.LastContainers, pClient.CurrentContainers...)
	err = json.Unmarshal(body, &pClient.CurrentContainers)
	if err != nil {
		return err
	}
	//fmt.Println(10)
	//fmt.Println("Last running containers: ", pClient.LastContainers)
	//fmt.Println("Current running containers: ", pClient.CurrentContainers, "\n\n")

	return nil
}

// Находит в списках работающих сейчас контейнеров и работающих Х сек. назад два одинаковых, если такие имеются так же сравнивает статусы
// так как тут хранятся все контейнера и стопнутые тоже, сравнивает на то разные ли они? и если разные то проверяет на то какой сейчас статус
// если exited то добавляет в список упавших
func (pClient *ClientPortaineer) StopedTrigger() {
	for _, lastContainer := range pClient.LastContainers {
		isDroped := false
		for _, currentContainer := range pClient.CurrentContainers {
			if (lastContainer.Names == currentContainer.Names &&
				lastContainer.State != currentContainer.State &&
				currentContainer.State == "exited") {
				isDroped = true
				break
			}
		}
		if isDroped {
			pClient.StopedContainers = append(pClient.StopedContainers, lastContainer)
		}
	}
}

// Создает список упавших контейнеров которые не восстановились
func (pClient *ClientPortaineer) GetDropedContainer() (types.Containers, error) {
	stopedWithError := make(types.Containers, 0)
	for _, stopedContainer := range pClient.StopedContainers {
		stopedWithError = append(stopedWithError, stopedContainer)
		//fmt.Println("\n",stopedContainer.State,"\n")
	}
	pClient.StopedContainers = make(types.Containers, 0)
	return stopedWithError, nil
}
