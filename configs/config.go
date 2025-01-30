package configs

import (
	"log"
	"os"
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
func LoadConfig(path string, logger *log.Logger, errorLogger *log.Logger) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("gocommerce")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		errorLogger.Fatalf("Config file not found or cannot be read: %v", err)
	} else {
		logger.Println("Config file loaded successfully")
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		errorLogger.Fatalf("Error unmarshalling config: %v", err)
	}

	return
}

// GetMongoDBURI reads the value of environment variable MONGODB_URI
// and returns the value
func GetMongoDBURI() string {
	mongoURI := os.Getenv("MONGODB_URI")
	return mongoURI
}
