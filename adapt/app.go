package adapt

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/leetatech/leeta_backend/adapt/routes"
	"github.com/leetatech/leeta_backend/pkg/config"
	"github.com/leetatech/leeta_backend/pkg/database"
	"github.com/leetatech/leeta_backend/pkg/jwtmiddleware"
	"github.com/leetatech/leeta_backend/pkg/mailer/aws"
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

	"github.com/leetatech/leeta_backend/pkg"
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

	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Application struct {
	Config            *config.ServerConfig
	Db                *mongo.Client
	Ctx               context.Context
	Router            *chi.Mux
	EmailClient       aws.MailClient
	RepositoryManager pkg.RepositoryManager
}

// New instances a new application
// The application contains all the related components that allow the execution of the service
func New(configFile string) (*Application, error) {
	var app Application
	var err error
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

	app.Db, err = database.Client(ctx, app.Config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongo db client %w", err)
	}
	if err := app.Db.Ping(ctx, readpref.Primary()); err != nil {
		return nil, errors.New("error pinging database")
	}

	app.EmailClient = aws.MailClient{
		Config: &app.Config.AWSConfig,
	}
	err = app.EmailClient.Connect()
	if err != nil {
		return nil, err
	}

	jwtManager, err := jwtmiddleware.New(app.Config.PublicKey, app.Config.PrivateKey)
	if err != nil {
		return nil, err
	}

	allInterfaces := app.buildApplicationConnection(*jwtManager, *app.Config)

	router, _, err := routes.SetupRouter(jwtManager, allInterfaces)
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
	return config.ReadConfig(configFile)
}

func (app *Application) buildApplicationConnection(jwtManager jwtmiddleware.Manager, config config.ServerConfig) *routes.AllHTTPHandlers {
	authPersistence := authInfrastructure.New(app.Db, app.Config.Database.DBName)
	orderPersistence := orderInfrastructure.New(app.Db, app.Config.Database.DBName)
	userPersistence := userInfrastructure.New(app.Db, app.Config.Database.DBName)
	productPersistence := productInfrastructure.New(app.Db, app.Config.Database.DBName)
	cartPersistence := cartInfrastructure.New(app.Db, app.Config.Database.DBName)
	feesPersistence := feesInfrastructure.New(app.Db, app.Config.Database.DBName)
	statePersistence := stateInfrastructure.New(app.Db, app.Config.Database.DBName)

	repositoryManager := pkg.RepositoryManager{
		OrderRepository:   orderPersistence,
		AuthRepository:    authPersistence,
		UserRepository:    userPersistence,
		ProductRepository: productPersistence,
		CartRepository:    cartPersistence,
		FeesRepository:    feesPersistence,
		StatesRepository:  statePersistence,
	}

	app.RepositoryManager = repositoryManager

	request := pkg.ApplicationContext{
		JwtManager:        jwtManager,
		RepositoryManager: repositoryManager,
		MailClient:        app.EmailClient,
		Domain:            app.Config.Notification.Domain,
		Config:            config,
	}

	orderApplications := orderApplication.New(request)
	authApplications := authApplication.New(request)
	userApplications := userApplication.New(request)
	productApplications := productApplication.New(request)
	cartsApplication := cartApplication.New(request)
	feeApplication := feesApplication.New(request)
	statesApplication := stateApplication.New(request, app.Config.NgnStates)

	orderInterfaces := orderInterface.New(orderApplications)
	authInterfaces := authInterface.New(authApplications)
	userInterfaces := userInterface.New(userApplications)
	productInterfaces := productInterface.New(productApplications)
	cartInterfaces := cartInterface.New(cartsApplication)
	feesInterfaces := feeInterface.New(feeApplication)
	statesInterfaces := stateInterface.New(statesApplication)

	allInterfaces := routes.AllHTTPHandlers{
		Order:   orderInterfaces,
		Auth:    authInterfaces,
		User:    userInterfaces,
		Product: productInterfaces,
		Cart:    cartInterfaces,
		Fees:    feesInterfaces,
		State:   statesInterfaces,
	}
	return routes.AllInterfaces(&allInterfaces)
}
