package config

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var conf *Config

type Config struct {
	AppPort            int    `mapstructure:"APP_PORT"`
	RunMode            string `mapstructure:"RUN_MODE"`
	DbType             string `mapstructure:"DB_TYPE"`
	DbHost             string `mapstructure:"DB_HOST"`
	DbPort             int    `mapstructure:"DB_PORT"`
	DbUser             string `mapstructure:"DB_USER"`
	DbPass             string `mapstructure:"DB_PASS"`
	DbName             string `mapstructure:"DB_NAME"`
	BrokerType         string `mapstructure:"BROKER_TYPE"`
	KafkaUrl           string `mapstructure:"KAFKA_URL"`
	KafkaConsumerGroup string `mapstructure:"KAFKA_CONSUMER_GROUP"`
	KafkaTopic         string `mapstructure:"KAFKA_TOPIC"`
	CacheType          string `mapstructure:"CACHE_TYPE"`
	RedisHost          string `mapstructure:"REDIS_HOST"`
	RedisPort          int    `mapstructure:"REDIS_PORT"`
	RedisPass          string `mapstructure:"REDIS_PASS"`
	RedisDatabase      int    `mapstructure:"REDIS_DATABASE"`
}

func (c *Config) Init(_ chan error) error {
	v := viper.NewWithOptions(viper.ExperimentalBindStruct())
	v.AddConfigPath(".")
	v.SetConfigName(".env")
	v.SetConfigType("env")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if os.Getenv("APP_PORT") == "" {
			return fmt.Errorf("error reading config file, %v", err)
		}
	}

	if err := v.Unmarshal(c); err != nil {
		return fmt.Errorf("unable to decode into struct: %v", err)
	}
	conf = c
	return nil
}

func (c *Config) SuccessfulMessage() string {
	return "Config successfully initialized"
}

func (c *Config) Shutdown(_ context.Context) error {
	return nil
}

func GetConfig() *Config {
	return conf
}
