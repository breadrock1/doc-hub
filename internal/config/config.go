package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"docs-hub/internal/cloud"
	"docs-hub/internal/server"
	"github.com/lpernett/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Cloud  cloud.CloudConfig
	Server server.Config
}

func FromFile(filePath string) (*Config, error) {
	config := &Config{}

	viperInstance := viper.New()
	viperInstance.AutomaticEnv()
	viperInstance.SetConfigFile(filePath)

	viperInstance.SetDefault("server.Address", "0.0.0.0:2866")
	viperInstance.SetDefault("server.LoggerLevel", "INFO")

	viperInstance.SetDefault("cloud.Address", "localhost:9000")
	viperInstance.SetDefault("cloud.Username", "minio-root")
	viperInstance.SetDefault("cloud.Password", "minio-root")
	viperInstance.SetDefault("cloud.EnableSSL", false)

	if err := viperInstance.ReadInConfig(); err != nil {
		confErr := fmt.Errorf("failed while reading config file %s: %w", filePath, err)
		return config, confErr
	}

	if err := viperInstance.Unmarshal(config); err != nil {
		confErr := fmt.Errorf("failed while unmarshaling config file %s: %w", filePath, err)
		return config, confErr
	}

	return config, nil
}

func LoadEnv(enableDotenv bool) (*Config, error) {
	if enableDotenv {
		_ = godotenv.Load()
	}

	servAddr := loadString("NEWS_WEEDER_SERVER_ADDRESS")
	servLogger := loadString("NEWS_WEEDER_SERVER_LOGGER_LEVEL")
	serverConfig := server.Config{Address: servAddr, LoggerLevel: servLogger}

	cloudAddr := loadString("DOCS_HUB_CLOUD_ADDRESS")
	cloudUser := loadString("DOCS_HUB_CLOUD_USERNAME")
	cloudPasswd := loadString("DOCS_HUB_CLOUD_PASSWORD")
	cloudEnableSSL := loadBool("DOCS_HUB_CLOUD_ENABLE_SSL")
	cloudConfig := cloud.CloudConfig{
		Address:   cloudAddr,
		Username:  cloudUser,
		Password:  cloudPasswd,
		EnableSSL: cloudEnableSSL,
	}

	return &Config{
		Cloud:  cloudConfig,
		Server: serverConfig,
	}, nil
}

func loadString(envName string) string {
	value, exists := os.LookupEnv(envName)
	if !exists {
		msg := fmt.Sprintf("failed to extract %s env var: %s", envName, value)
		log.Println(msg)
		return ""
	}
	return value
}

func loadNumber(envName string, bitSize int) int {
	value, exists := os.LookupEnv(envName)
	if !exists {
		msg := fmt.Sprintf("failed to extract %s env var: %s", envName, value)
		log.Println(msg)
		return 0
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		msg := fmt.Sprintf("failed to convert %s env var: %s", envName, value)
		log.Println(msg)
		return 0
	}

	return int(number)
}

func loadBool(envName string) bool {
	value, exists := os.LookupEnv(envName)
	if !exists {
		msg := fmt.Sprintf("faile to extract %s env var: %s", envName, value)
		log.Println(msg)
		return false
	}

	boolean, err := strconv.ParseBool(value)
	if err != nil {
		msg := fmt.Sprintf("faile to convert %s env var: %s", envName, value)
		log.Println(msg)
		return false
	}

	return boolean
}
