package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	User       string
	Password   string
	Dbname     string
	Sqlhost    string
	Sqlport    string
	Webserver  string
	Webport    string
	Sessionkey string
}

var Cfg *Config
var Fullurl string

func LoadConfig() *Config {
	raw, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	configuration := Config{}
	err = json.Unmarshal(raw, &configuration)
	return &configuration
}

func init() {
	Cfg = LoadConfig()
	if k := os.Getenv("GOVWA_SESSIONKEY"); k != "" {
		Cfg.Sessionkey = k
	}
	if len(Cfg.Sessionkey) < 32 {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			panic(err)
		}
		Cfg.Sessionkey = hex.EncodeToString(b)
	}
	url := fmt.Sprintf("%s:%s/", Cfg.Webserver, Cfg.Webport)
	Fullurl = url
}
