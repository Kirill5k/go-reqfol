package config

import (
	"github.com/spf13/viper"
	"log"
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

	var app App
	if err := v.Unmarshal(&app); err != nil {
		log.Fatalf("failed to decode viper config into struct. %v", err)
	}
	return &app
}
