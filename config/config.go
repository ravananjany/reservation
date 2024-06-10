package config

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type config struct {
	Logfile  string
	Loglevel string
	Apiport  string
	Grpcport string
	Term     bool
}

var defaultLogger = "blog.log"

func LoadConfig() (*logrus.Logger, config) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	var configuration config
	logger := logrus.New()
	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.Fatalf("Unable to decode into struct in config, %v", err)
	}

	fmt.Print(configuration)

	level, err := logrus.ParseLevel(configuration.Loglevel)
	if err != nil {
		logger.Fatal("invalid logger level")
	}
	logrus.SetLevel(level)

	if len(configuration.Logfile) < 1 {
		return logger, configuration
	}

	file, err := os.OpenFile(configuration.Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		logger.Info("Failed to log to file, using default stderr")
	}
	return logger, configuration
}
