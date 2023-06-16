package adapt

import (
	"fmt"
	"github.com/caarlos0/env"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"time"
)

type ServerConfig struct {
	AppEnv     string `env:"APP_ENV" envDefault:"dev" envWhitelisted:"true"`
	HTTPPort   int    `env:"PORT" envDefault:"3000" envWhitelisted:"true"`
	Database   DatabaseConfig
	PrivateKey string `env:"PRIVATE_KEY"`
	PublicKey  string `env:"PUBLIC_KEY"`
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

func Read(logger zap.Logger) (*ServerConfig, error) {
	var serverConfig ServerConfig

	for _, target := range []interface{}{
		&serverConfig,
		&serverConfig.Database,
		//&serverConfig.Security,
	} {
		if err := env.Parse(target); err != nil {
			return &serverConfig, err
		}
	}

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
		SetHosts([]string{config.Database.Host + config.Database.Port}).
		SetAuth(options.Credential{
			AuthMechanism: "SCRAM-SHA-256",
			Username:      config.Database.UserName,
			Password:      config.Database.Password,
		})
}
