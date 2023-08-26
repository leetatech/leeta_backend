package adapt

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/leetatech/leeta_backend/adapt/routes"
	authApplication "github.com/leetatech/leeta_backend/services/auth/application"
	authInfrastructure "github.com/leetatech/leeta_backend/services/auth/infrastructure"
	authInterface "github.com/leetatech/leeta_backend/services/auth/interfaces"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/library/mailer"
	orderApplication "github.com/leetatech/leeta_backend/services/order/application"
	orderInfrastructure "github.com/leetatech/leeta_backend/services/order/infrastructure"
	orderInterface "github.com/leetatech/leeta_backend/services/order/interfaces"
	productInfrastructure "github.com/leetatech/leeta_backend/services/product/infrastructure"
	userApplication "github.com/leetatech/leeta_backend/services/user/application"
	userInfrastructure "github.com/leetatech/leeta_backend/services/user/infrastructure"
	userInterface "github.com/leetatech/leeta_backend/services/user/interfaces"
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
	EmailClient  mailer.MailerClient
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

	app.EmailClient = mailer.NewMailerClient(app.Config.Postmark.Key, app.Logger)

	tokenHandler, err := library.NewMiddlewares(app.Config.PublicKey, app.Config.PrivateKey, app.Logger)
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
	log.Println("Application running on port ", app.Config.HTTPPort)
	log.Println("Access swagger docs on {PORT}/api/swagger/", app.Config.HTTPPort) //should be updated if route is ever changed
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
	authPersistence := authInfrastructure.NewAuthPersistence(app.Db, app.Config.Database.DbName, app.Logger)
	orderPersistence := orderInfrastructure.NewOrderPersistence(app.Db, app.Config.Database.DbName, app.Logger)
	userPersistence := userInfrastructure.NewUserPersistence(app.Db, app.Config.Database.DbName, app.Logger)
	productPersistence := productInfrastructure.NewProductPersistence(app.Db, app.Config.Database.DbName, app.Logger)

	allRepositories := library.Repositories{
		OrderRepository:   orderPersistence,
		AuthRepository:    authPersistence,
		UserRepository:    userPersistence,
		ProductRepository: productPersistence,
	}

	app.Repositories = allRepositories
	request := library.DefaultApplicationRequest{
		TokenHandler:  tokenHandler,
		Logger:        app.Logger,
		AllRepository: allRepositories,
		EmailClient:   app.EmailClient,
		Domain:        app.Config.Leeta.Domain,
	}

	orderApplications := orderApplication.NewOrderApplication(tokenHandler, allRepositories)
	authApplications := authApplication.NewAuthApplication(request)
	userApplications := userApplication.NewUserApplication(request)

	orderInterfaces := orderInterface.NewOrderHTTPHandler(orderApplications)
	authInterfaces := authInterface.NewAuthHttpHandler(authApplications)
	userInterfaces := userInterface.NewUserHttpHandler(userApplications)

	allInterfaces := routes.AllHTTPHandlers{
		Order: orderInterfaces,
		Auth:  authInterfaces,
		User:  userInterfaces,
	}
	return routes.AllInterfaces(&allInterfaces)
}
