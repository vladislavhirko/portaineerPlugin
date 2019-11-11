package mattermost

import (
	"encoding/json"
	"github.com/mattermost/mattermost-server/model"
	"github.com/vladislavhirko/portaineerPlugin/database"
	"github.com/vladislavhirko/portaineerPlugin/mattermost/types"
	pTypes "github.com/vladislavhirko/portaineerPlugin/portainer/types"
	"log"
)

type MattermostClient struct {
	Client  *model.Client4
	User    *model.User
	Chanels types.Chanels
	DB      database.LevelDB
}

func NewMattermostClient(ldb database.LevelDB, address, port string) MattermostClient {
	url := "http://" + address + ":" + port
	return MattermostClient{
		Client:  model.NewAPIv4Client(url),
		User:    &model.User{},
		Chanels: make(types.Chanels, 0),
		DB:      ldb,
	}
}

//Логинит пользователя в систему
func (mClient *MattermostClient) Auth(login, password string) error {
	user, resp := mClient.Client.Login(login, password)
	if resp.Error != nil {
		return resp.Error
	}
	mClient.User = user
	return nil
}

//Получает список всех каналов и кладет в струкутруу
//
func (mClient *MattermostClient) GetallChanels() error {
	chanels, resp := mClient.Client.GetAllChannels(0, 100, "")
	if resp.Error != nil {
		return resp.Error
	}
	err := json.Unmarshal([]byte(chanels.ToJson()), &mClient.Chanels)
	if err != nil {
		return err
	}
	return nil
}

//Отправляет сообщение всем каналам
// Делается запрос с бд для каждого упавшего контейнера, достаем название канала, далее бегаем по всем каналам и в нужный нам отправляем инфу
func (mClient *MattermostClient) SendMessage(containers pTypes.Containers, patternChanel string) error {
	// Листает список всех каналов и когда находит тот который в бд, отправляет туда
	for _, container := range containers {
		chanelName, err := mClient.DB.Get(container.Names[0])
		if err != nil {
			log.Println("No chanel for ", container)
			continue
		}
		for _, chanel := range mClient.Chanels {
			if chanelName == chanel.Name {
				post := &model.Post{}
				post.ChannelId = chanel.ID
				post.Message = "Fault: " + container.Names[0] + "\n" + container.Logs
				_, resp := mClient.Client.CreatePost(post)
				if resp.Error != nil {
					return resp.Error
				}

			}
		}
	}
	return nil
}