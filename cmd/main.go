package main

//go:generate -command swag go run github.com/swaggo/swag/cmd/swag@latest
//go:generate swag init --parseDependency --parseInternal -o ../docs

import (
	"fmt"
	"github.com/leetatech/leeta_backend/adapt"
	"github.com/leetatech/leeta_backend/services/library/logger"
	"os"
)

// @title			LEETA BACKEND API
// @version		1.0
// @description	LEETA Application backend documentation
// @termsOfService	http://swagger.io/terms/
// @contact.name	LEETA Technologies
// @contact.email	admin@getlleta.com
// @license.name	Apache 3.0-or-later
// @host			localhost:3000
// @BasePath		/api
func main() {
	appLogger := logger.New()

	app, err := adapt.New(appLogger)
	if err != nil {
		appLogger.Error(fmt.Sprintf("Fatal error creating application: %v", err))
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		appLogger.Error(fmt.Sprintf("Fatal error running application: %v", err))
		os.Exit(1)
	}
}
