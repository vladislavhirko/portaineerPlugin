package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	API API `toml:"api"`
	LevelDB Level      `toml:"level"`
	MClient Mattermost `toml:"mattermost"`
	PClient Portainer  `toml:"portainer"`
}

type API struct{
	Port string `toml:"port"`
}

type Level struct {
	Path string `toml:"path"`
}

type Mattermost struct {
	Address  string `toml:"address"`
	Port     string `toml:"port"`
	Email    string `toml:"email"`
	Password string `toml:"password"`
}

type Portainer struct {
	Login         string `toml:"login"`
	Password      string `toml:"password"`
	Address       string `toml:"address"`
	Port          string `toml:"port"`
	CheckInterval string `toml:"check_interval"`
	LogsAmount string `toml:"logs_amount"`
}

//Parse config file
func GetConfig(path string) (*Config, error) {
	//fmt.Println(path)
	//err := CreateEnvironment()
	config := Config{
		API: *new(API),
		LevelDB: *new(Level),
		MClient: *new(Mattermost),
		PClient: *new(Portainer),
	}
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

//Fuction which will create folder with config and storage
//func CreateEnvironment() error{
//	usr, _ := user.Current()
//	_, err := os.Stat(usr.HomeDir + "/.portaineerPlugin1")
//	if os.IsNotExist(err){
//		err := os.Mkdir(usr.HomeDir + "/.portaineerPlugin1", 0777)
//		if err != nil{
//			return err
//		}
//	}
//	return nil
//}
