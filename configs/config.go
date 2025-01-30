package configs

import (
	"time"

	"github.com/spf13/viper"
)

// Config defines fields for MongoDB and RabbitMQ connection details, port number, and JWT
// secret key.
type Config struct {
	MongoDBURI     string        `mapstructure:"MONGODB_URI"`
	MongoDBName    string        `mapstructure:"MONGODB_NAME"`
	RabbitMQURI    string        `mapstructure:"RABBITMQ_URI"`
	OrderQueueName string        `mapstructure:"ORDER_QUEUE_NAME"`
	Port           string        `mapstructure:"PORT"`
	DefaultTimeout time.Duration `mapstructure:"DEFAULT_TIMEOUT"`
	JWTSecretKey   string        `mapstructure:"JWT_SECRET_KEY"`
}

// LoadConfig loads environment variable into the Config struct by
// unmarshalling them
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("gocommerce")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	config.DefaultTimeout = viper.GetDuration("DEFAULT_TIMEOUT")

	err = viper.Unmarshal(&config)
	return
}

func GetConfig() (*Config, error) {
	config, err := LoadConfig(".")
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func GetMongoDBURI() string {
	config, _ := GetConfig()
	return config.MongoDBURI
}
