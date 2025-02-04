package configs

import (
	"fmt"
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

func GetEnvironmentVariable(config *Config) {
	config.MongoDBURI = os.Getenv("MONGODB_URI")
	config.MongoDBName = os.Getenv("MONGODB_NAME")
	config.RabbitMQURI = os.Getenv("RABBITMQ_URI")
	config.OrderQueueName = os.Getenv("ORDER_QUEUE_NAME")
	config.Port = os.Getenv("PORT")
	timeDuration, err := time.ParseDuration(os.Getenv("DEFAULT_TIMEOUT"))
	if err != nil {
		timeDuration = time.Duration(time.Second * 30)
	}
	config.DefaultTimeout = timeDuration
	config.JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
}

// LoadConfig loads environment variable into the Config struct by
// unmarshalling them
func LoadConfig(path string, logger *log.Logger, errorLogger *log.Logger) (config Config, err error) {
	viper.SetConfigFile(path)

	err = viper.ReadInConfig()
	if err != nil {
		errorLogger.Println("Config file not found or cannot be read; using environment variables only")
		GetEnvironmentVariable(&config)
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

// SetTestConfigFile loads and/or sets environment variables from
// a .env or from environment variables `viper.AutomaticEnv`
func SetTestConfigFile() (config Config) {
	viper.SetConfigFile("../gocommerce.env")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config file not found or cannot be read; using environment variables only")
	} else {
		fmt.Println("Config file loaded successfully")
	}

	viper.AutomaticEnv()

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("Error unmarshalling config: %v", err)
		os.Exit(1) // Exit on failure
	}
	config.MongoDBName = "GoCommerceTest"
	config.MongoDBURI = GetMongoDBURI()

	return
}
