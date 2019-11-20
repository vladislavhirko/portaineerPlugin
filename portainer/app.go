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

//Struct of client for portainer
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

//Creates a new client for work with portainer, and allocate memory for slices of containers
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

//Auth by login and password, setup JWT. All next requests to portainer must contain JWT token
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

// Takes all containers list and set it to struct variable
//  thus before this saves all containers which worked before refreshing
//(containers list which worked 'X' second ago and list which works now)
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
	return nil
}

//TODO заменить два слайса на 2 мапы
//Finds same container in 2 lists (1: which worked 'X' second ago; 2: which workes now).
//After that compares theis states, if states are different
//Checks for exited state in last containers list
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

// Create list of stoped containerm which didn't recovery
func (pClient *ClientPortaineer) GetDropedContainer() (types.Containers, error) {
	stopedWithError := make(types.Containers, 0)
	err := pClient.GetDropedLogs()
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
