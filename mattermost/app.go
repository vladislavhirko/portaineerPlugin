package mattermost

import (
	"encoding/json"
	"github.com/mattermost/mattermost-server/model"
	"github.com/vladislavhirko/portaineerPlugin/database"
	"github.com/vladislavhirko/portaineerPlugin/mattermost/types"
	pTypes "github.com/vladislavhirko/portaineerPlugin/portainer/types"
	log "github.com/sirupsen/logrus"


)

type MattermostClient struct {
	Client  *model.Client4
	User    *model.User
	Chanels types.Chanels
	DB      database.LevelDB
	LogContext *log.Entry
}

func NewMattermostClient(ldb database.LevelDB, address, port string) MattermostClient {
	url := "http://" + address + ":" + port
	return MattermostClient{
		Client:  model.NewAPIv4Client(url),
		User:    &model.User{},
		Chanels: make(types.Chanels, 0),
		DB:      ldb,
		LogContext: log.WithFields(log.Fields{
			"Module": "Mattermost",
		}),
	}
}

//Login user to system
func (mClient *MattermostClient) Auth(login, password string) error {
	user, resp := mClient.Client.Login(login, password)
	if resp.Error != nil {
		return resp.Error
	}
	mClient.User = user
	return nil
}

//Gets list of all chanels and put it to struct
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


//Creates request to database for each stopped container, takes chanel name for each, after that we look through all channels and send information to linked chanel with container
func (mClient *MattermostClient) SendMessage(containers pTypes.Containers, patternChanel string) error {
	// Листает список всех каналов и когда находит тот который в бд, отправляет туда
	// Look through all channels and after finding chanel in database, send info to it
	for _, container := range containers {
		mClient.LogContext.Warn("Container fault ", container)
		chanelName, err := mClient.DB.DBContainerChat.Get(container.Names[0])
		if err != nil {
			log.Error(err)
		}
		for _, chanel := range mClient.Chanels {
			if chanelName == chanel.Name {
				post := &model.Post{}
				post.ChannelId = chanel.ID
				post.Message = "Fault: **_" + container.Names[0] + "_**\n```" + container.Logs[:len(container.Logs) - 1] + "```"
				_, resp := mClient.Client.CreatePost(post)
				if resp.Error != nil {
					return resp.Error
				}
			}
		}
	}
	return nil
}