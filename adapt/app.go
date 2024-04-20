package adapt

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/leetatech/leeta_backend/adapt/routes"
	"github.com/leetatech/leeta_backend/pkg/config"
	"github.com/leetatech/leeta_backend/pkg/database"
	stateApplication "github.com/leetatech/leeta_backend/services/state/application"
	stateInfrastructure "github.com/leetatech/leeta_backend/services/state/infrastructure"
	stateInterface "github.com/leetatech/leeta_backend/services/state/interfaces"

	authApplication "github.com/leetatech/leeta_backend/services/auth/application"
	authInfrastructure "github.com/leetatech/leeta_backend/services/auth/infrastructure"
	authInterface "github.com/leetatech/leeta_backend/services/auth/interfaces"
	"github.com/rs/zerolog/log"

	cartApplication "github.com/leetatech/leeta_backend/services/cart/application"
	cartInfrastructure "github.com/leetatech/leeta_backend/services/cart/infrastructure"
	cartInterface "github.com/leetatech/leeta_backend/services/cart/interfaces"

	gasrefillApplication "github.com/leetatech/leeta_backend/services/gasrefill/application"
	gasrefillInfrastructure "github.com/leetatech/leeta_backend/services/gasrefill/infrastructure"
	gasrefillInterface "github.com/leetatech/leeta_backend/services/gasrefill/interfaces"

	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/mailer"

	orderApplication "github.com/leetatech/leeta_backend/services/order/application"
	orderInfrastructure "github.com/leetatech/leeta_backend/services/order/infrastructure"
	orderInterface "github.com/leetatech/leeta_backend/services/order/interfaces"

	productApplication "github.com/leetatech/leeta_backend/services/product/application"
	productInfrastructure "github.com/leetatech/leeta_backend/services/product/infrastructure"
	productInterface "github.com/leetatech/leeta_backend/services/product/interfaces"

	userApplication "github.com/leetatech/leeta_backend/services/user/application"
	userInfrastructure "github.com/leetatech/leeta_backend/services/user/infrastructure"
	userInterface "github.com/leetatech/leeta_backend/services/user/interfaces"

	feesApplication "github.com/leetatech/leeta_backend/services/fees/application"
	feesInfrastructure "github.com/leetatech/leeta_backend/services/fees/infrastructure"
	feeInterface "github.com/leetatech/leeta_backend/services/fees/interfaces"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Application struct {
	Logger       *zap.Logger
	Config       *config.ServerConfig
	Db           *mongo.Client
	Ctx          context.Context
	Router       *chi.Mux
	EmailClient  mailer.MailerClient
	Repositories pkg.Repositories
}

// New instances a new application
// The application contains all the related components that allow the execution of the service
func New(logger *zap.Logger, configFile string) (*Application, error) {
	var app Application
	var err error
	app.Logger = logger
	app.Config, err = app.buildConfig(configFile)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//build application clients
	// verify application config
	if app.Config == nil {
		return nil, errors.New("application config is empty")
	}

	app.Db, err = database.MongoDBClient(ctx, app.Config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongo db client %w", err)
	}
	if err := app.Db.Ping(ctx, readpref.Primary()); err != nil {
		return nil, errors.New("error pinging database")
	}

	app.EmailClient = mailer.NewMailerClient(pkg.PostMarkAPIToken, app.Logger)

	tokenHandler, err := pkg.NewMiddlewares(app.Config.PublicKey, app.Config.PrivateKey, app.Logger)
	if err != nil {
		return nil, err
	}

	allInterfaces := app.buildApplicationConnection(*tokenHandler)

	router, _, err := routes.SetupRouter(tokenHandler, allInterfaces)
	if err != nil {
		return nil, err
	}
	app.Router = router

	app.Ctx = ctx

	return &app, nil
}

// Run executes the application
func (app *Application) Run() error {
	defer func() {
		if err := app.Db.Disconnect(app.Ctx); err != nil {
			log.Debug().Msgf("error disconnecting from database: %v", err)
		}
	}()

	app.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/swagger/", http.StatusFound)
	})

	app.Router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("Leeta Backend Server is running..."))
		if err != nil {
			log.Debug().Msgf("Error writing response: %v", err)
		}
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

func (app *Application) buildConfig(configFile string) (*config.ServerConfig, error) {
	return config.ReadConfig(*app.Logger, configFile)
}

func (app *Application) buildApplicationConnection(tokenHandler pkg.TokenHandler) *routes.AllHTTPHandlers {
	authPersistence := authInfrastructure.NewAuthPersistence(app.Db, app.Config.Database.DBName, app.Logger)
	orderPersistence := orderInfrastructure.NewOrderPersistence(app.Db, app.Config.Database.DBName, app.Logger)
	userPersistence := userInfrastructure.NewUserPersistence(app.Db, app.Config.Database.DBName, app.Logger)
	productPersistence := productInfrastructure.NewProductPersistence(app.Db, app.Config.Database.DBName, app.Logger)
	gasRefillPersistence := gasrefillInfrastructure.NewRefillPersistence(app.Db, app.Config.Database.DBName, app.Logger)
	cartPersistence := cartInfrastructure.NewCartPersistence(app.Db, app.Config.Database.DBName, app.Logger)
	feesPersistence := feesInfrastructure.NewFeesPersistence(app.Db, app.Config.Database.DBName, app.Logger)
	statePersistence := stateInfrastructure.NewStatePersistence(app.Db, app.Config.Database.DBName, app.Logger)

	allRepositories := pkg.Repositories{
		OrderRepository:     orderPersistence,
		AuthRepository:      authPersistence,
		UserRepository:      userPersistence,
		ProductRepository:   productPersistence,
		GasRefillRepository: gasRefillPersistence,
		CartRepository:      cartPersistence,
		FeesRepository:      feesPersistence,
		StatesRepository:    statePersistence,
	}

	app.Repositories = allRepositories

	request := pkg.DefaultApplicationRequest{
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
	cartsApplication := cartApplication.NewCartApplication(request)
	feeApplication := feesApplication.NewFeesApplication(request)
	statesApplication := stateApplication.NewStateApplication(request, app.Config.NgnStates)

	orderInterfaces := orderInterface.NewOrderHTTPHandler(orderApplications)
	authInterfaces := authInterface.NewAuthHttpHandler(authApplications)
	userInterfaces := userInterface.NewUserHttpHandler(userApplications)
	productInterfaces := productInterface.NewProductHTTPHandler(productApplications)
	gasRefillInterfaces := gasrefillInterface.NewGasRefillHTTPHandler(gasRefillApplications)
	cartInterfaces := cartInterface.NewCartHTTPHandler(cartsApplication)
	feesInterfaces := feeInterface.NewFeesHTTPHandler(feeApplication)
	statesInterfaces := stateInterface.NewStateHttpHandler(statesApplication)

	allInterfaces := routes.AllHTTPHandlers{
		Order:     orderInterfaces,
		Auth:      authInterfaces,
		User:      userInterfaces,
		Product:   productInterfaces,
		GasRefill: gasRefillInterfaces,
		Cart:      cartInterfaces,
		Fees:      feesInterfaces,
		State:     statesInterfaces,
	}
	return routes.AllInterfaces(&allInterfaces)
}
