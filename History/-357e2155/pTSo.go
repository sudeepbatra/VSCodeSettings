package config

import (
	"fmt"
	"log"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type ConfigStructure struct {
	Version string `json:"version"`

	BrokerConfig struct {
		Breeze struct {
			Name   string `mapstructure:"app_name"`
			Key    string `mapstructure:"app_key"`
			Secret string `mapstructure:"secret_key"`
		} `mapstructure:"breeze"`
		SmartApi struct {
			ClientCode string `mapstructure:"client_code"`
			Key        string `mapstructure:"app_key"`
			Password   string `mapstructure:"password"`
		} `mapstructure:"smart_api"`
		SmartAPIHistorical struct {
			ClientCode string `mapstructure:"client_code"`
			APIKey     string `mapstructure:"api_key"`
			SecretKey  string `mapstructure:"secret_key"`
			Password   string `mapstructure:"password"`
		} `mapstructure:"smart_api_historical"`
		SmartAPITrading struct {
			ClientCode string `mapstructure:"client_code"`
			APIKey     string `mapstructure:"api_key"`
			SecretKey  string `mapstructure:"secret_key"`
			Password   string `mapstructure:"password"`
		} `mapstructure:"smart_api_trading"`
	} `mapstructure:"broker"`

	LoggingConfig struct {
		LogLevel zerolog.Level `mapstructure:"logLevel"`
	} `mapstructure:"logging"`

	DatabaseConfig struct {
		MainTsdbUri string `mapstructure:"main_tsdb_uri"`
	} `mapstructure:"database_config"`
}

var Config ConfigStructure

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("settings")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("Unable to find the secrets file.")
		} else {
			log.Fatalf("unable to decode into struct, %v", err)
		}
	}

	fmt.Println(viper.AllSettings())
	err := viper.Unmarshal(&Config)

	if err != nil {
		log.Fatalf("unable to decode configs into struct, %v", err)
	}

	fmt.Println(Config)
}
