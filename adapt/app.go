package adapt

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/leetatech/leeta_backend/adapt/routes"
	authApplication "github.com/leetatech/leeta_backend/services/auth/application"
	authInfrastructure "github.com/leetatech/leeta_backend/services/auth/infrastructure"
	authInterface "github.com/leetatech/leeta_backend/services/auth/interfaces"

	cartApplication "github.com/leetatech/leeta_backend/services/cart/application"
	cartInfrastructure "github.com/leetatech/leeta_backend/services/cart/infrastructure"
	cartInterface "github.com/leetatech/leeta_backend/services/cart/interfaces"

	gasrefillApplication "github.com/leetatech/leeta_backend/services/gasrefill/application"
	gasrefillInfrastructure "github.com/leetatech/leeta_backend/services/gasrefill/infrastructure"
	gasrefillInterface "github.com/leetatech/leeta_backend/services/gasrefill/interfaces"

	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/library/mailer"

	orderApplication "github.com/leetatech/leeta_backend/services/order/application"
	orderInfrastructure "github.com/leetatech/leeta_backend/services/order/infrastructure"
	orderInterface "github.com/leetatech/leeta_backend/services/order/interfaces"

	productApplication "github.com/leetatech/leeta_backend/services/product/application"
	productInfrastructure "github.com/leetatech/leeta_backend/services/product/infrastructure"
	productInterface "github.com/leetatech/leeta_backend/services/product/interfaces"

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//build application clients
	app.Db = app.buildMongoClient(ctx)
	if err := app.Db.Ping(ctx, readpref.Primary()); err != nil {
		app.Logger.Info("msg", zap.String("msg", "failed to ping to database"))
		log.Fatal(err)
	}

	//defer app.db.Disconnect(ctx)

	app.EmailClient = mailer.NewMailerClient(library.PostMarkAPIToken, app.Logger)

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

	app.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/swagger/", http.StatusFound)
	})

	app.Router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Welcome to the leeta Server.."))
	})

	svr := http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.HTTPPort),
		Handler: app.Router,
	}
	fmt.Println("Application running on port ", app.Config.HTTPPort)
	fmt.Printf("Access swagger docs on host://%v/api/swagger/, app.Config.HTTPPort)", app.Config.HTTPPort)
	err := svr.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) buildConfig() (*ServerConfig, error) {
	return ReadConfig(*app.Logger)
}

func (app *Application) buildApplicationConnection(tokenHandler library.TokenHandler) *routes.AllHTTPHandlers {
	authPersistence := authInfrastructure.NewAuthPersistence(app.Db, app.Config.Database.DBName, app.Logger)
	orderPersistence := orderInfrastructure.NewOrderPersistence(app.Db, app.Config.Database.DBName, app.Logger)
	userPersistence := userInfrastructure.NewUserPersistence(app.Db, app.Config.Database.DBName, app.Logger)
	productPersistence := productInfrastructure.NewProductPersistence(app.Db, app.Config.Database.DBName, app.Logger)
	gasRefillPersistence := gasrefillInfrastructure.NewRefillPersistence(app.Db, app.Config.Database.DBName, app.Logger)
	cartPersistence := cartInfrastructure.NewCartPersistence(app.Db, app.Config.Database.DBName, app.Logger)

	allRepositories := library.Repositories{
		OrderRepository:     orderPersistence,
		AuthRepository:      authPersistence,
		UserRepository:      userPersistence,
		ProductRepository:   productPersistence,
		GasRefillRepository: gasRefillPersistence,
		CartRepository:      cartPersistence,
	}

	app.Repositories = allRepositories
	request := library.DefaultApplicationRequest{
		TokenHandler:  tokenHandler,
		Logger:        app.Logger,
		AllRepository: allRepositories,
		EmailClient:   app.EmailClient,
		Domain:        app.Config.Leeta.Domain,
	}

	orderApplications := orderApplication.NewOrderApplication(request)
	authApplications := authApplication.NewAuthApplication(request)
	userApplications := userApplication.NewUserApplication(request)
	productApplications := productApplication.NewProductApplication(request)
	gasRefillApplications := gasrefillApplication.NewGasRefillApplication(request)
	cartApplication := cartApplication.NewCartApplication(request)

	orderInterfaces := orderInterface.NewOrderHTTPHandler(orderApplications)
	authInterfaces := authInterface.NewAuthHttpHandler(authApplications)
	userInterfaces := userInterface.NewUserHttpHandler(userApplications)
	productInterfaces := productInterface.NewProductHTTPHandler(productApplications)
	gasRefillInterfaces := gasrefillInterface.NewGasRefillHTTPHandler(gasRefillApplications)
	cartInterfaces := cartInterface.NewCartHTTPHandler(cartApplication)

	allInterfaces := routes.AllHTTPHandlers{
		Order:     orderInterfaces,
		Auth:      authInterfaces,
		User:      userInterfaces,
		Product:   productInterfaces,
		GasRefill: gasRefillInterfaces,
		Cart:      cartInterfaces,
	}
	return routes.AllInterfaces(&allInterfaces)
}
