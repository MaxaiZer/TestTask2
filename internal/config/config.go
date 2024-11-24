package config

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
)

type Config struct {
	Port               int    `mapstructure:"PORT"`
	Mode               string `mapstructure:"MODE"`
	DbConnectionString string `mapstructure:"DB_CONNECTION_STRING"`
}

const configPath = "./configs/config.env"

func Get() *Config {

	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func loadConfig(path string) (*Config, error) {

	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		path = envPath
	}

	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	viper.SetDefault("PORT", defaultPort)
	viper.SetDefault("MODE", defaultMode)

	if err := viper.ReadInConfig(); err != nil {
		currentDir, dirErr := os.Getwd()
		if dirErr != nil {
			return nil, dirErr
		}

		return nil, fmt.Errorf("error reading config file, %s. currect directory: %s", err, currentDir)
	}

	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	err := config.validate()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c Config) validate() error {

	var errs []error

	if c.DbConnectionString == "" {
		errs = append(errs, fmt.Errorf("missing variable DbConnectionString"))
	}

	if c.Port <= 0 {
		errs = append(errs, fmt.Errorf("invalid port: %d", c.Port))
	}

	if c.Mode != ReleaseMode && c.Mode != DebugMode {
		errs = append(errs, fmt.Errorf("invalid mode: %s", c.Mode))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred: %w", errors.Join(errs...))
	}

	return nil
}
