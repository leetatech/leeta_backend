package main

//go:generate -command swag go run github.com/swaggo/swag/cmd/swag@v1.8.9
//go:generate swag init --parseDependency --parseInternal -o ../docs

import (
	"flag"
	"github.com/leetatech/leeta_backend/adapt"
	"github.com/rs/zerolog/log"
	"os"
)

// @title			LEETA BACKEND API
// @version		1.0
// @description	LEETA Application backend documentation
// @termsOfService	http://swagger.io/terms/
// @contact.name	LEETA Technologies
// @contact.email	admin@getlleta.com
// @license.name	Apache 3.0-or-later
// @host			https://leetabackend-e6d948d15ae2.herokuapp.com
// @BasePath		/api
// @securityDefinitions.apikey BearerToken
// @in header
// @name authorization
func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "local.env", "configuration file")
	flag.StringVar(&configFile, "c", "local.env", "configuration file (shorthand)")
	flag.Parse()

	app, err := adapt.New(configFile)
	if err != nil {
		log.Error().Msgf("Fatal error creating application: %v", err)
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		log.Error().Msgf("Fatal error running application: %v", err)
		os.Exit(1)
	}
}
