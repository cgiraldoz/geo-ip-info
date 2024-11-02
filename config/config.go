package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func LoadConfigurations() error {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	fixerApiKey := viper.GetString("FIXER_API_KEY")
	if fixerApiKey == "" {
		log.Fatal("FIXER_API_KEY is not set")
	}

	for _, key := range viper.AllKeys() {
		value := viper.GetString(key)
		viper.Set(key, strings.ReplaceAll(value, "FIXER_API_KEY", fixerApiKey))
	}
	return nil
}
