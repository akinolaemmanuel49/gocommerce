package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	MongoDBURI     string `mapstructure:"MongoDB_URI"`
	MongoDBName    string `mapstructure:"MongoDB_Name"`
	RedisAddr      string
	RabbitMQURI    string `mapstructure:"RabbitMQ_URI"`
	OrderQueueName string `mapstructure:"Order_Queue_Name"`
	Port           string
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("gocommerce")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
