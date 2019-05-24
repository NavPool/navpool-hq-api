package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Config struct {
	Debug           bool         `yaml:"debug"`
	Server          ServerConfig `yaml:"server"`
	JWT             JWTConfig    `yaml:"jwt"`
	DB              DBConfig     `yaml:"db"`
	Explorer        ApiConfig    `yaml:"explorer"`
	Pool            ApiConfig    `yaml:"pool"`
	Sentry          SentryConfig `yaml:"sentry"`
	SelectedNetwork string       `yaml:"selectedNetwork"`
}

type ServerConfig struct {
	Port   string `yaml:"port"`
	Domain string `yaml:"domain"`
}

type JWTConfig struct {
	Realm       string `yaml:"realm"`
	Secret      string `yaml:"secret"`
	IdentityKey string `yaml:"identityKey"`
}

type DBConfig struct {
	Dialect  string `yaml:"dialect"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbName"`
	SSLMode  string `yaml:"sslMode"`
	LogMode  bool   `yaml:"logMode"`
}

type ApiConfig struct {
	Url string `yaml:"url"`
}

type SentryConfig struct {
	Active bool   `yaml:"active"`
	DSN    string `yaml:"dsn"`
}

var instance *Config
var once sync.Once

func Get() *Config {
	once.Do(func() {
		filePath := fmt.Sprintf("./config.%s.yaml", env())
		log.Printf("ConfigFile: %s", filePath)

		configFile, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		instance = &Config{}
		err = yaml.Unmarshal(configFile, instance)
		if err != nil {
			log.Fatal(err)
		}

		if instance.Debug {
			configJson, _ := json.Marshal(instance)
			log.Printf("Config: %s", string(configJson))
		}
	})
	return instance
}

func env() string {
	var env = "prod"
	if len(os.Args) > 1 {
		env = os.Args[1]
	}
	log.Print("Environment: " + env)

	return env
}
