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

type Application struct {
	Logger       *zap.Logger
	Config       *ServerConfig
	Db           *mongo.Client
	Ctx          context.Context
	Router       *chi.Mux
	Repositories library.Repositories
}

// New instances a new application
// The application contains all the related components that allow the execution of the service
func New(logger *zap.Logger) (*Application, error) {
	var app Application
	var err error
	app.Logger = logger
	app.Config, err = app.buildConfig()

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//build application clients
	app.Db = app.buildMongoClient(ctx)

	if err := app.Db.Ping(ctx, readpref.Primary()); err != nil {
		app.Logger.Info("msg", zap.String("msg", "failed to ping to database"))
		log.Fatal(err)
	}

	//defer app.db.Disconnect(ctx)

	tokenHandler, err := library.NewMiddlewares(app.Config.PublicKey, app.Config.PrivateKey)
	if err != nil {
		return nil, err
	}

	allInterfaces := app.buildApplicationConnection(*tokenHandler)

	router, tokenHandler, err := routes.SetupRouter(tokenHandler, allInterfaces)
	if err != nil {
		return nil, err
	}
	app.Router = router

	app.Ctx = ctx
	return &app, nil
}

// Run executes the application
func (app *Application) Run() error {
	defer app.Db.Disconnect(app.Ctx)

	app.Router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Welcome to the leeta Server.."))
	})

	svr := http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.HTTPPort),
		Handler: app.Router,
	}
	err := svr.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) buildConfig() (*ServerConfig, error) {
	return Read(*app.Logger)
}

func (app *Application) buildApplicationConnection(tokenHandler library.TokenHandler) *routes.AllHTTPHandlers {
	orderPersistence := infrastructure.NewOrderPersistence(app.Db, app.Config.Database.DbName, app.Logger)

	allRepositories := library.Repositories{
		OrderRepository: orderPersistence,
	}
	app.Repositories = allRepositories

	orderApplications := orderApplication.NewOrderApplication(tokenHandler, allRepositories)

	orderInterfaces := interfaces.NewOrderHTTPHandler(orderApplications)
	allInterfaces := routes.AllHTTPHandlers{
		Order: orderInterfaces,
	}
	return routes.AllInterfaces(&allInterfaces)
}
