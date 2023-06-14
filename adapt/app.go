package adapt

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/leetatech/leeta_backend/adapt/routes"
	"github.com/leetatech/leeta_backend/services/library"
	orderApplication "github.com/leetatech/leeta_backend/services/order/application"
	"github.com/leetatech/leeta_backend/services/order/infrastructure"
	"github.com/leetatech/leeta_backend/services/order/interfaces"
	"go.uber.org/zap"
	"log"
	"net/http"
)

type application struct {
	logger       *zap.Logger
	config       *ServerConfig
	db           *sql.DB
	router       *chi.Mux
	repositories library.Repositories
}

// New instances a new application
// The application contains all the related components that allow the execution of the service
func New(logger *zap.Logger) (*application, error) {
	var app application
	var err error

	app.logger = logger
	app.config, err = app.buildConfig()

	if err != nil {
		return nil, err
	}
	//build application clients
	app.db = app.buildSqlClient()

	if err := app.db.PingContext(context.Background()); err != nil {
		app.logger.Info("msg", zap.String("msg", "failed to ping to database"))
		log.Fatal(err)
	}

	tokenHandler, err := library.NewMiddlewares()
	if err != nil {
		return nil, err
	}

	allInterfaces := app.buildApplicationConnection(*tokenHandler)

	router, tokenHandler, err := routes.SetupRouter(tokenHandler, allInterfaces)
	if err != nil {
		return nil, err
	}
	app.router = router

	return &app, nil
}

// Run executes the application
func (app *application) Run() error {
	defer app.db.Close()

	app.router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Welcome to the leeta Server.."))
	})

	svr := http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.HTTPPort),
		Handler: app.router,
	}
	err := svr.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (app *application) buildConfig() (*ServerConfig, error) {
	return Read(*app.logger)
}

func (app *application) buildSqlClient() *sql.DB {
	db := Database{Config: app.config, Log: app.logger}

	dbConn, err := db.ConnectDB()

	if err != nil {
		log.Fatal(err)
	}

	if err := db.RunMigration(dbConn); err != nil {
		log.Fatal(err)
	}

	return dbConn
}

func (app *application) buildApplicationConnection(tokenHandler library.TokenHandler) *routes.AllHTTPHandlers {
	orderPersistences := infrastructure.NewOrderPersistence(app.db, app.logger)

	allRepositories := library.Repositories{
		OrderRepository: orderPersistences,
	}
	app.repositories = allRepositories

	orderApplications := orderApplication.NewOrderApplication(tokenHandler, allRepositories)

	orderInterfaces := interfaces.NewOrderHTTPHandler(orderApplications)
	allInterfaces := routes.AllHTTPHandlers{
		Order: orderInterfaces,
	}
	return routes.AllInterfaces(&allInterfaces)
}
