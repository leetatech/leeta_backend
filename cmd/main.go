package main

import (
	"fmt"
	"github.com/leetatech/leeta_backend/adapt"
	"github.com/leetatech/leeta_backend/services/library/logger"
	"os"
)

//go:generate -command swag go run github.com/swaggo/swag/cmd/swag@latest

// @title			LEETA BACKEND API
// @version		1.0
// @description	This is the entire doc
// @termsOfService	http://swagger.io/terms/
// @contact.name	LEETA Engineering
// @contact.email	leeta.org
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:3000
// @BasePath		/leeta
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
