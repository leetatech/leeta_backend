package adapt

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"net/url"
	"os"
	"strings"
	"time"
)

type ServerConfig struct {
	AppEnv     string `env:"APP_ENV" envDefault:"dev" envWhitelisted:"true"`
	HTTPPort   int    `env:"PORT" envDefault:"3000" envWhitelisted:"true"`
	Database   DatabaseConfig
	PrivateKey string `env:"PRIVATE_KEY"`
	PublicKey  string `env:"PUBLIC_KEY"`
	Postmark   PostmarkConfig
	Leeta      LeetaConfig
}

type DatabaseConfig struct {
	Host     string `env:"MONGO_HOST" envDefault:"localhost:"`
	Port     string `env:"MONGO_PORT" envDefault:"27017"`
	Timeout  int    `env:"MONGO_CONNECTION_TIMEOUT_SECONDS" envDefault:"10"`
	DbName   string `env:"MONGO_DB_NAME" envDefault:"leeta"`
	UserName string `env:"MONGO_USERNAME" envDefault:"leeta"`
	Password string `env:"MONGO_PASSWORD" envDefault:"leet"`
	DbUrl    string `env:"DATABASE_URL" envDefault:"" envWhitelisted:"true"`
}

type PostmarkConfig struct {
	URL string `env:"POSTMARK_URL"`
	Key string `env:"POSTMARK_KEY"`
}

type LeetaConfig struct {
	Domain string `env:"DOMAIN"`
}

func Read(logger zap.Logger) (*ServerConfig, error) {
	var serverConfig ServerConfig

	if err := godotenv.Load("../local.env"); err != nil {
		return nil, err
	}

	for _, target := range []interface{}{
		&serverConfig,
		&serverConfig.Database,
		&serverConfig.Postmark,
		&serverConfig.Leeta,
		//&serverConfig.Security,
	} {
		if err := env.Parse(target); err != nil {
			return &serverConfig, err
		}
	}
	overrideWithCommandLine(serverConfig)

	out := serverConfig.formartUri()
	logger.Info(out)
	return &serverConfig, nil
}

func (config *ServerConfig) formartUri() string {
	format := "database: {host: %s port:%s timeout:%d, username-hidden password-hidden}"
	host := config.Database.Host
	port := config.Database.Port
	timeout := config.Database.Timeout

	if config.Database.DbUrl != "" {
		if connString, err := url.Parse(config.Database.DbUrl); err == nil {

			result := strings.Split(connString.Host, ":")
			host = result[0]
			port = result[1]
		}
	}

	return fmt.Sprintf(format, host, port, timeout)
}

func (config *ServerConfig) GetClientOptions() *options.ClientOptions {
	return options.Client().
		SetConnectTimeout(time.Duration(config.Database.Timeout) * time.Second).
		SetHosts([]string{config.Database.Host + config.Database.Port})
}

func overrideWithCommandLine(serverConfig ServerConfig) {
	if privateKey := os.Getenv("PRIVATE_KEY"); privateKey != "" {
		serverConfig.PrivateKey = privateKey
	}

	if publicKey := os.Getenv("PUBLIC_KEY"); publicKey != "" {
		serverConfig.PublicKey = publicKey
	}

	if postmarkURL := os.Getenv("POSTMARK_URL"); postmarkURL != "" {
		serverConfig.Postmark.URL = postmarkURL
	}

	if postmarkKey := os.Getenv("POSTMARK_KEY"); postmarkKey != "" {
		serverConfig.Postmark.Key = postmarkKey
	}
}
