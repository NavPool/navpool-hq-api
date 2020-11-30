package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Debug    bool         `yaml:"debug"`
	Server   ServerConfig `yaml:"server"`
	JWT      JWTConfig    `yaml:"jwt"`
	DB       DBConfig     `yaml:"db"`
	Explorer ApiConfig    `yaml:"explorer"`
	Pool     ApiConfig    `yaml:"pool"`
	Sentry   Sentry       `yaml:"sentry"`
	Network  string       `yaml:"network"`
}

type ServerConfig struct {
	Port int
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

type Sentry struct {
	Active bool   `yaml:"active"`
	DSN    string `yaml:"dsn"`
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.WithError(err).Fatal("Unable to init config")
	}
}

func Get() *Config {
	return &Config{
		Debug: getBool("DEBUG", false),
		Server: ServerConfig{
			Port: getInt("SERVER_PORT", 8080),
		},
		JWT: JWTConfig{
			Realm:       getString("JWT_REALM", "NavPool HQ"),
			Secret:      getString("JWT_SECRET", ""),
			IdentityKey: getString("JWT_IDENTITY_KEY", "id"),
		},
		DB: DBConfig{
			Dialect:  getString("DB_DIALECT", "postgres"),
			Host:     getString("DB_HOST", "localhost"),
			Port:     getInt("DB_PORT", 8432),
			Username: getString("DB_USERNAME", "navpool"),
			Password: getString("DB_PASSWORD", "navpool"),
			DbName:   getString("DB_NAME", "navpool"),
			SSLMode:  getString("DB_SSL_MODE", "disable"),
			LogMode:  getBool("DB_LOG_MODE", false),
		},
		Explorer: ApiConfig{
			Url: getString("EXPLORER_URL", "https://api.navexplorer.com"),
		},
		Pool: ApiConfig{
			Url: getString("POOL_URL", "http://api:8080"),
		},
		Sentry: Sentry{
			Active: getBool("SENTRY_ACTIVE", false),
			DSN:    getString("SENTRY_DSN", ""),
		},
		Network: getString("NETWORK", "mainnet"),
	}
}

func getString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func getInt(key string, defaultValue int) int {
	valStr := getString(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}

	return defaultValue
}

func getUint(key string, defaultValue uint) uint {
	return uint(getInt(key, int(defaultValue)))
}

func getUint64(key string, defaultValue uint) uint64 {
	return uint64(getInt(key, int(defaultValue)))
}

func getBool(key string, defaultValue bool) bool {
	valStr := getString(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultValue
}

func getSlice(key string, defaultVal []string, sep string) []string {
	valStr := getString(key, "")
	if valStr == "" {
		return defaultVal
	}

	return strings.Split(valStr, sep)
}
