package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/leetatech/leeta_backend/docs"
	"github.com/leetatech/leeta_backend/pkg"
	authInterfaces "github.com/leetatech/leeta_backend/services/auth/interfaces"
	cartInterfaces "github.com/leetatech/leeta_backend/services/cart/interfaces"
	feesInterfaces "github.com/leetatech/leeta_backend/services/fees/interfaces"
	orderInterfaces "github.com/leetatech/leeta_backend/services/order/interfaces"
	productInterfaces "github.com/leetatech/leeta_backend/services/product/interfaces"
	stateInterfaces "github.com/leetatech/leeta_backend/services/state/interfaces"
	userInterfaces "github.com/leetatech/leeta_backend/services/user/interfaces"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type AllHTTPHandlers struct {
	Order   *orderInterfaces.OrderHttpHandler
	Auth    *authInterfaces.AuthHttpHandler
	User    *userInterfaces.UserHttpHandler
	Product *productInterfaces.ProductHttpHandler
	Cart    *cartInterfaces.CartHttpHandler
	Fees    *feesInterfaces.FeesHttpHandler
	State   *stateInterfaces.StateHttpHandler
}

func AllInterfaces(interfaces *AllHTTPHandlers) *AllHTTPHandlers {
	return &AllHTTPHandlers{
		Order:   interfaces.Order,
		Auth:    interfaces.Auth,
		User:    interfaces.User,
		Product: interfaces.Product,
		Cart:    interfaces.Cart,
		Fees:    interfaces.Fees,
		State:   interfaces.State,
	}
}

func SetupRouter(tokenHandler *pkg.TokenHandler, interfaces *AllHTTPHandlers) (*chi.Mux, *pkg.TokenHandler, error) {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Use(middleware.Logger)

	orderRouter := buildOrderEndpoints(*interfaces.Order, tokenHandler)
	authRouter := buildAuthEndpoints(*interfaces.Auth)
	userRouter := buildUserEndpoints(*interfaces.User, tokenHandler)
	productRouter := buildProductEndpoints(*interfaces.Product, tokenHandler)
	cartRouter := buildCartEndpoints(*interfaces.Cart, tokenHandler)
	feesRouter := buildFeesEndpoints(*interfaces.Fees, tokenHandler)
	stateRouter := buildStatesEndpoints(*interfaces.State, tokenHandler)

	router.Route("/api", func(r chi.Router) {
		r.Handle("/swagger/*", httpSwagger.WrapHandler)
		r.Mount("/session", authRouter)
		r.Mount("/order", orderRouter)
		r.Mount("/user", userRouter)
		r.Mount("/product", productRouter)
		r.Mount("/cart", cartRouter)
		r.Mount("/fees", feesRouter)
		r.Mount("/state", stateRouter)
	})

	return router, tokenHandler, nil
}

func buildAuthEndpoints(session authInterfaces.AuthHttpHandler) http.Handler {
	router := chi.NewRouter()

	// Signing
	router.Post("/signup", session.SignUpHandler)
	router.Post("/signin", session.SignInHandler)
	router.Post("/admin/signup", session.AdminSignUpHandler)

	// guest session management
	router.Post("/guest", session.ReceiveGuestTokenHandler)
	router.Get("/guest/{device_id}", session.GetGuestRecordHandler)
	router.Put("/guest", session.UpdateGuestRecordHandler)

	// otp
	router.Route("/otp", func(r chi.Router) {
		r.Post("/request", session.RequestOTPHandler)
		r.Post("/validate", session.ValidateOTPHandler)
	})

	// password
	router.Route("/password", func(r chi.Router) {
		r.Post("/forgot", session.ForgotPasswordHandler)
		r.Post("/create", session.CreateNewPasswordHandler)
	})

	// earlyAccess
	router.Route("/early_access", func(r chi.Router) {
		r.Post("/", session.EarlyAccessHandler)
	})

	return router
}

func buildOrderEndpoints(order orderInterfaces.OrderHttpHandler, tokenHandler *pkg.TokenHandler) http.Handler {
	router := chi.NewRouter()
	router.Use(tokenHandler.ValidateMiddleware)
	router.Put("/status", order.UpdateOrderStatusHandler)
	router.Get("/id/{order_id}", order.GetOrderByIDHandler)
	router.Get("/", order.GetCustomerOrdersByStatusHandler)
	router.Put("/", order.ListOrdersHandler)
	router.Get("/options", order.ListOrdersOptions)
	router.Get("/status/history/{order_id}", order.ListOrderStatusHistoryHandler)
	return router
}

func buildUserEndpoints(user userInterfaces.UserHttpHandler, tokenHandler *pkg.TokenHandler) http.Handler {
	router := chi.NewRouter()

	router.Mount("/vendor", buildVendorEndpoints(user, tokenHandler))

	router.Group(func(r chi.Router) {
		r.Use(tokenHandler.ValidateMiddleware)
		r.Put("/", user.UpdateUserData)
	})

	return router
}

func buildVendorEndpoints(handler userInterfaces.UserHttpHandler, tokenHandler *pkg.TokenHandler) http.Handler {
	router := chi.NewRouter()
	router.Use(tokenHandler.ValidateRestrictedAccessMiddleware)

	// authentication group here
	router.Group(func(r chi.Router) {
		r.Use(tokenHandler.ValidateMiddleware)
		r.Post("/verification", handler.VendorVerificationHandler)
		r.Post("/admin/vendor", handler.AddVendorByAdminHandler)
	})

	// non-authentication group here

	return router
}

func buildProductEndpoints(product productInterfaces.ProductHttpHandler, tokenHandler *pkg.TokenHandler) http.Handler {
	router := chi.NewRouter()
	router.Use(tokenHandler.ValidateMiddleware)

	// Restricted route group
	router.Route("/", func(r chi.Router) {
		r.Use(tokenHandler.ValidateRestrictedAccessMiddleware)
		r.Post("/", product.CreateGasProductHandler)
	})

	// Unrestricted routes
	router.Get("/id/{product_id}", product.GetProductByIDHandler)
	router.Get("/", product.GetAllVendorProductsHandler)
	router.Put("/", product.ListProductsHandler)
	router.Get("/options", product.ListProductOptions)
	router.Post("/create", product.CreateProductHandler)

	return router
}

func buildCartEndpoints(handler cartInterfaces.CartHttpHandler, tokenHandler *pkg.TokenHandler) http.Handler {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(tokenHandler.ValidateMiddleware)
		// post endpoints
		r.Post("/add", handler.AddToCart)
		r.Put("/", handler.ListCart)
		r.Post("/checkout", handler.Checkout)

		// get endpoints
		r.Get("/options", handler.ListCartOptions)

		// update endpoints
		r.Put("/item/quantity", handler.UpdateCartItemQuantity)

		// delete endpoints
		r.Delete("/{cart_id}", handler.DeleteCart)
		r.Delete("/item/{cart_item_id}", handler.DeleteCartItem)
	})

	return router
}

func buildFeesEndpoints(handler feesInterfaces.FeesHttpHandler, tokenHandler *pkg.TokenHandler) http.Handler {
	router := chi.NewRouter()
	router.Use(tokenHandler.ValidateMiddleware)

	// Restricted route group
	router.Route("/", func(r chi.Router) {
		r.Use(tokenHandler.ValidateRestrictedAccessMiddleware)
		router.Post("/", handler.CreateFeeHandler)
	})
	router.Put("/", handler.FetchFeesHandler)
	router.Get("/options", handler.ListFeesOptions)
	return router
}

func buildStatesEndpoints(handler stateInterfaces.StateHttpHandler, tokenHandler *pkg.TokenHandler) http.Handler {
	router := chi.NewRouter()

	router.Use(tokenHandler.ValidateMiddleware)
	router.Post("/", handler.RetrieveNGNStatesData)
	router.Get("/{name}", handler.GetState)
	router.Get("/", handler.ListStates)

	return router
}
