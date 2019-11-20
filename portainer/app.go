package portainer

import (
	"bytes"
	"encoding/json"
	"github.com/vladislavhirko/portaineerPlugin/portainer/types"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	log "github.com/sirupsen/logrus"

)

//Структура клиента для портейнера
type ClientPortaineer struct {
	Jwt               string `json:"jwt"`
	Username          string
	Password          string
	Address           string
	Port              string
	CheckInterval     time.Duration
	CurrentContainers types.Containers
	LastContainers    types.Containers
	StopedContainers  types.Containers
	ContainerLogsAmount        string
	Log *log.Entry `json:"-"`
}

//Создание нового клиента портейнера, выделяет память для слайсов контейнеров
func NewPorteinerClient(address, port, checkInterval string, logsAmount string) *ClientPortaineer {
	checkIntervalInt, err := strconv.Atoi(checkInterval)
	if err != nil {
		log.Fatal(err)
	}
	checkIntervalDuration := time.Duration(checkIntervalInt)
	return &ClientPortaineer{
		Jwt:               "",
		Username:          "",
		Password:          "",
		Address:           address,
		Port:              port,
		CheckInterval:     checkIntervalDuration,
		CurrentContainers: make(types.Containers, 0),
		LastContainers:    make(types.Containers, 0),
		StopedContainers:  make(types.Containers, 0),
		ContainerLogsAmount:        logsAmount,
		Log: log.WithFields(log.Fields{
				"Module": "Portainer",
			}),
	}
}

//Аутентификация по логину и паролю, устанавливает JWT токен,
//должен передаваться во все последующие запросы на портейнер в заголовке
func (pClient *ClientPortaineer) Auth(login, password string) error {
	pClient.Username = login
	pClient.Password = password

	authObjJSON, err := json.Marshal(&pClient)
	if err != nil {
		return err
	}
	r := bytes.NewReader(authObjJSON)
	resp, err := http.Post("http://"+pClient.Address+":"+pClient.Port+"/api/auth", "application/json", r)
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
		"GET", "http://"+pClient.Address+":"+pClient.Port+"/api/endpoints/1/docker/containers/json?all=1", nil,
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
	pClient.Log.Trace("Tick interval: ", pClient.CheckInterval, ". Time: ", time.Now())
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
//TODO заменить два слайса на 2 мапы
func (pClient *ClientPortaineer) FinedDropedContainers() {
	for _, lastContainer := range pClient.LastContainers {
		for _, currentContainer := range pClient.CurrentContainers {
			if lastContainer.Names == currentContainer.Names &&
				lastContainer.State != currentContainer.State &&
				currentContainer.State == "exited" {
				pClient.StopedContainers = append(pClient.StopedContainers, lastContainer)
				break
			}
		}
	}
}

// Создает список упавших контейнеров которые не восстановились
func (pClient *ClientPortaineer) GetDropedContainer() (types.Containers, error) {
	stopedWithError := make(types.Containers, 0)
	err := pClient.GetDropedLogs() //в будущем тут будут доставаться логи упавших контейнеров
	if err != nil {
		return nil, err
	}
	for _, stopedContainer := range pClient.StopedContainers {
		stopedWithError = append(stopedWithError, stopedContainer)
	}
	pClient.StopedContainers = make(types.Containers, 0)
	return stopedWithError, nil
}

func (pClient *ClientPortaineer) GetDropedLogs() error {
	for i := 0; i < len(pClient.StopedContainers); i++ {
		client := &http.Client{}
		req, err := http.NewRequest(
			"GET",
			"http://"+pClient.Address+":"+pClient.Port+"/api/endpoints/1/docker/containers/"+
				pClient.StopedContainers[i].Id[:12]+"/logs?stderr=1&stdout=1&follow=1&tail="+pClient.ContainerLogsAmount,
			nil,
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
		//fmt.Println(string(body))
		pClient.StopedContainers[i].Logs = string(body)
	}
	return nil
}
