package config

import (
	"log"

	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/spf13/viper"
)

func Init(env string) *models.Config {
	var conf models.Config

	viper.SetConfigType("yaml")
	viper.SetConfigName(env)
	viper.AddConfigPath("./cmd")

	viper.SetEnvPrefix("UMF")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error on parsing configuration file: %+v", err)
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Fatal("Error on unmarshal config to struct")
	}
	log.Printf("Loaded config: %+v", conf)

	return &conf
}
