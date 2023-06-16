package adapt

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/leetatech/leeta_backend/adapt/routes"
	"github.com/leetatech/leeta_backend/services/library"
	orderApplication "github.com/leetatech/leeta_backend/services/order/application"
	"github.com/leetatech/leeta_backend/services/order/infrastructure"
	"github.com/leetatech/leeta_backend/services/order/interfaces"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

type application struct {
	logger       *zap.Logger
	config       *ServerConfig
	db           *mongo.Client
	ctx          context.Context
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//build application clients
	app.db = app.buildMongoClient(ctx)

	if err := app.db.Ping(ctx, readpref.Primary()); err != nil {
		app.logger.Info("msg", zap.String("msg", "failed to ping to database"))
		log.Fatal(err)
	}

	//defer app.db.Disconnect(ctx)

	tokenHandler, err := library.NewMiddlewares(app.config.PublicKey, app.config.PrivateKey)
	if err != nil {
		return nil, err
	}

	allInterfaces := app.buildApplicationConnection(*tokenHandler)

	router, tokenHandler, err := routes.SetupRouter(tokenHandler, allInterfaces)
	if err != nil {
		return nil, err
	}
	app.router = router

	app.ctx = ctx
	return &app, nil
}

// Run executes the application
func (app *application) Run() error {
	defer app.db.Disconnect(app.ctx)

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

func (app *application) buildApplicationConnection(tokenHandler library.TokenHandler) *routes.AllHTTPHandlers {
	orderPersistence := infrastructure.NewOrderPersistence(app.db, app.config.Database.DbName, app.logger)

	allRepositories := library.Repositories{
		OrderRepository: orderPersistence,
	}
	app.repositories = allRepositories

	orderApplications := orderApplication.NewOrderApplication(tokenHandler, allRepositories)

	orderInterfaces := interfaces.NewOrderHTTPHandler(orderApplications)
	allInterfaces := routes.AllHTTPHandlers{
		Order: orderInterfaces,
	}
	return routes.AllInterfaces(&allInterfaces)
}
