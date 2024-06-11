package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"time"
)

type Interrupter struct {
	InitialDelay time.Duration `mapstructure:"initial-delay"`
}

type Server struct {
	Port int `mapstructure:"port"`
}

type Client struct {
	MaxIdleConns        int           `mapstructure:"max-idle-conns"`
	MaxIdleConnsPerHost int           `mapstructure:"max-idle-conns-per-host"`
	IdleConnTimeout     time.Duration `mapstructure:"idle-conn-timeout"`
	Timeout             time.Duration `mapstructure:"timeout"`
	RetryCount          int           `mapstructure:"retry-count"`
	RetryWaitTime       time.Duration `mapstructure:"retry-wait-time"`
	RetryMaxWaitTime    time.Duration `mapstructure:"retry-max-wait-time"`
}

type App struct {
	Server      Server
	Client      Client
	Interrupter Interrupter
}

func LoadAppConfig() *App {
	v := viper.New()
	v.SetConfigName("application")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("failed to read viper config. %v", err)
	}

	for _, k := range v.AllKeys() {
		value := v.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			envVarName, defaultValue := getEnvVarNameWithDefaultValue(value)
			v.Set(k, getEnv(envVarName, defaultValue))
		}
	}

	var app App
	if err := v.Unmarshal(&app); err != nil {
		log.Fatalf("failed to decode viper config into struct. %v", err)
	}
	return &app
}

func getEnvVarNameWithDefaultValue(stringTemplate string) (string, string) {
	envVarName := strings.TrimSuffix(strings.TrimPrefix(stringTemplate, "${"), "}")
	if strings.Contains(envVarName, ":") {
		split := strings.SplitN(envVarName, ":", 2)
		return split[0], split[1]
	}
	return envVarName, ""
}

func getEnv(envVarName, defaultValue string) string {
	value, found := os.LookupEnv(envVarName)
	if found {
		return value
	}
	if !found && defaultValue != "" {
		return defaultValue
	}
	panic(fmt.Sprintf("Missing required environment variable %s", envVarName))
}
