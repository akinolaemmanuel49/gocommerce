package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	MongoDBURI     string `mapstructure:"MONGODB_URI"`
	MongoDBName    string `mapstructure:"MONGODB_NAME"`
	RabbitMQURI    string `mapstructure:"RABBITMQ_URI"`
	OrderQueueName string `mapstructure:"ORDER_QUEUE_NAME"`
	Port           string `mapstructure:"PORT"`
	DefaultTimeout string `mapstructure:"DEFAULT_TIMEOUT"`
	JWTSecretKey   string `mapstructure:"JWT_SECRET_KEY"`
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
