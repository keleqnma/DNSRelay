package common

import (
	"log"

	"github.com/spf13/viper"
)

// LoadConfig 载入配置
func LoadConfig(in string, paths ...string) {
	viper.SetConfigName(in)
	viper.AddConfigPath("config/")
	for _, configPath := range paths {
		viper.AddConfigPath(configPath)
	}
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal("fail to load config file:", err)
	}
}
