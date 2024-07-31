package config

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	databaseTimeout = 10 * time.Second
	devMongoURI     = "mongodb+srv://admin:qT5IsndbYrzmq9eW@cluster0.rt4wdpi.mongodb.net/?retryWrites=true&w=majority"
)

type ServerConfig struct {
	AppEnv       string `env:"APP_ENV" envDefault:"staging" envWhitelisted:"true"`
	HTTPPort     int    `env:"PORT" envDefault:"3000" envWhitelisted:"true"`
	Database     DatabaseConfig
	PrivateKey   string `env:"PRIVATE_KEY"`
	PublicKey    string `env:"PUBLIC_KEY"`
	Postmark     PostmarkConfig
	Notification NotificationConfig
	NgnStates    NgnStatesConfig // configure resource API to retrieve NGN states
	AWSConfig    AWSConfig
}

type DatabaseConfig struct {
	Host     string `env:"MONGO_HOST" envDefault:"localhost:"`
	Port     string `env:"MONGO_PORT" envDefault:"27017"`
	Timeout  int    `env:"MONGO_CONNECTION_TIMEOUT_SECONDS" envDefault:"10"`
	DBName   string `env:"MONGO_DB_NAME" envDefault:"leeta"`
	UserName string `env:"MONGO_USERNAME" envDefault:"leeta"`
	Password string `env:"MONGO_PASSWORD" envDefault:"leet"`
	DbURL    string `env:"DATABASE_URL" envDefault:"" envWhitelisted:"true"`
}

type PostmarkConfig struct {
	URL string `env:"POSTMARK_URL"`
	Key string `env:"POSTMARK_KEY"`
}

type NotificationConfig struct {
	Domain            string `env:"DOMAIN"`
	VerificationEmail string `env:"LEETA_VERIFICATION_EMAIL" envDefault:"admin@getleeta.com"`
	DoNotReplyEmail   string `env:"LEETA_DONOTREPLY_EMAIL"`
}

type NgnStatesConfig struct {
	URL string `env:"URL" envDefault:"https://api.facts.ng/v1"`
}

type AWSConfig struct {
	Region   string `env:"AWS_REGION"`
	Endpoint string `env:"AWS_ENDPOINT"`
	Secret   string `env:"AWS_SECRET"`
}

func LoadEnv(configFile string) error {
	err := godotenv.Load(configFile)
	if err != nil {
		return fmt.Errorf("failed to load environment variables: %v", err)
	}
	return nil
}

func ReadConfig(configFile string) (*ServerConfig, error) {
	err := LoadEnv(configFile)
	if err != nil {
		return nil, err
	}

	var serverConfig ServerConfig
	targets := []interface{}{
		&serverConfig,
		&serverConfig.Database,
		&serverConfig.Postmark,
		&serverConfig.Notification,
		&serverConfig.NgnStates,
		&serverConfig.AWSConfig,
	}

	for _, target := range targets {
		if err := env.Parse(target); err != nil {
			return nil, fmt.Errorf("failed to parse environment variables: %v", err)
		}
	}

	overrideWithEnvVars(&serverConfig)
	out := serverConfig.formatURI()
	log.Debug().Msgf("config: %v", out)

	return &serverConfig, nil
}

func (config *ServerConfig) formatURI() string {
	format := "database: {host: %s port:%s timeout:%d, username-hidden password-hidden}"
	host := config.Database.Host
	port := config.Database.Port
	timeout := config.Database.Timeout

	if config.Database.DbURL != "" {
		if connString, err := url.Parse(config.Database.DbURL); err == nil {
			result := strings.Split(connString.Host, ":")
			host = result[0]
			port = result[1]
		}
	}

	return fmt.Sprintf(format, host, port, timeout)
}

func overrideWithEnvVars(config *ServerConfig) {
	config.PrivateKey = os.Getenv("PRIVATE_KEY")
	config.PublicKey = os.Getenv("PUBLIC_KEY")
	config.Postmark.URL = os.Getenv("POSTMARK_URL")
	config.Postmark.Key = os.Getenv("POSTMARK_KEY")
}

func (config *ServerConfig) GetClientOptions() *options.ClientOptions {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	if config.AppEnv == "dev" {
		return options.Client().
			SetConnectTimeout(databaseTimeout).
			ApplyURI(devMongoURI).
			SetServerAPIOptions(serverAPI)
	}

	connectionString := fmt.Sprintf("mongodb://%s%s", config.Database.Host, config.Database.Port)
	return options.Client().
		SetConnectTimeout(databaseTimeout).
		ApplyURI(connectionString)
}
