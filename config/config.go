package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"sync"
)

type Config struct {
	Debug bool

	Server struct {
		Port   string
		Domain string
	}

	JWT struct {
		Realm       string
		Secret      string
		IdentityKey string
	}

	DB struct {
		Dialect  string
		Host     string
		Port     int
		Username string
		Password string
		DbName   string
		SSLMode  string
	}

	ElasticSearch struct {
		Urls        string
		Sniff       bool
		HealthCheck bool
	}

	Networks []Network

	SelectedNetwork string
}

type Network struct {
	Name string

	Host     string
	Port     int
	Username string
	Password string
}

var instance *Config
var once sync.Once

func Get() *Config {
	once.Do(func() {
		log.Println("Creating Config")
		var env = "prod"
		if len(os.Args) > 1 {
			env = os.Args[1]
		}

		viper.SetConfigName("config." + env)
		viper.AddConfigPath(".")

		instance = &Config{}

		if err := viper.ReadInConfig(); err != nil {
			log.Fatal(err)
		}

		if err := viper.Unmarshal(instance); err != nil {
			log.Fatal(err)
		}
	})

	return instance
}
