package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/leetatech/leeta_backend/docs"
	authInterfaces "github.com/leetatech/leeta_backend/services/auth/interfaces"
	"github.com/leetatech/leeta_backend/services/library"
	orderInterfaces "github.com/leetatech/leeta_backend/services/order/interfaces"
	productInterfaces "github.com/leetatech/leeta_backend/services/product/interfaces"
	userInterfaces "github.com/leetatech/leeta_backend/services/user/interfaces"

	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type AllHTTPHandlers struct {
	Order   *orderInterfaces.OrderHttpHandler
	Auth    *authInterfaces.AuthHttpHandler
	User    *userInterfaces.UserHttpHandler
	Product *productInterfaces.ProductHttpHandler
}

func AllInterfaces(interfaces *AllHTTPHandlers) *AllHTTPHandlers {
	return &AllHTTPHandlers{Order: interfaces.Order, Auth: interfaces.Auth, User: interfaces.User, Product: interfaces.Product}
}

func SetupRouter(tokenHandler *library.TokenHandler, interfaces *AllHTTPHandlers) (*chi.Mux, *library.TokenHandler, error) {
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
	router.Route("/api", func(r chi.Router) {
		r.Handle("/swagger/*", httpSwagger.WrapHandler)
		r.Mount("/session", authRouter)
		r.Mount("/order", orderRouter)
		r.Mount("/user", userRouter)
		r.Mount("/product", productRouter)
	})

	return router, tokenHandler, nil
}

func buildAuthEndpoints(session authInterfaces.AuthHttpHandler) http.Handler {
	router := chi.NewRouter()

	// Signing
	router.Post("/signup", session.SignUpHandler)
	router.Post("/signin", session.SignInHandler)
	router.Post("/admin/signup", session.AdminSignUpHandler)

	// otp
	router.Post("/otp/request", session.RequestOTPHandler)
	router.Post("/otp/validate", session.ValidateOTPHandler)

	// password
	router.Post("/forgot_password", session.ForgotPasswordHandler)
	router.Post("/reset_password", session.ResetPasswordHandler)

	// earlyAccess
	router.Post("/early_access", session.EarlyAccessHandler)

	return router
}

func buildOrderEndpoints(order orderInterfaces.OrderHttpHandler, tokenHandler *library.TokenHandler) http.Handler {
	router := chi.NewRouter()
	router.Use(tokenHandler.ValidateMiddleware)
	router.Post("/make_order", order.CreateOrderHandler)
	router.Put("/status", order.UpdateOrderStatusHandler)
	router.Get("/id/{order_id}", order.GetOrderByIDHandler)
	router.Get("/", order.GetCustomerOrdersByStatusHandler)
	return router
}

func buildUserEndpoints(user userInterfaces.UserHttpHandler, tokenHandler *library.TokenHandler) http.Handler {
	router := chi.NewRouter()

	router.Mount("/vendor", buildVendorEndpoints(user, tokenHandler))

	return router
}

func buildVendorEndpoints(user userInterfaces.UserHttpHandler, tokenHandler *library.TokenHandler) http.Handler {
	router := chi.NewRouter()

	// authentication group here
	router.Group(func(r chi.Router) {
		r.Use(tokenHandler.ValidateMiddleware)
		r.Post("/verification", user.VendorVerificationHandler)
		r.Post("/admin/vendor", user.AddVendorByAdminHandler)
	})

	// non-authentication group here

	return router
}

func buildProductEndpoints(product productInterfaces.ProductHttpHandler, tokenHandler *library.TokenHandler) http.Handler {
	router := chi.NewRouter()
	router.Use(tokenHandler.ValidateMiddleware)
	router.Post("/create", product.CreateProductHandler)
	router.Get("/id/{product_id}", product.GetProductByIDHandler)
	router.Get("/", product.GetAllVendorProductsHandler)
	return router
}
