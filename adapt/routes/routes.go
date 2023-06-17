package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/leetatech/leeta_backend/docs"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/order/interfaces"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type AllHTTPHandlers struct {
	Order *interfaces.HTTPHandler
}

func AllInterfaces(interfaces *AllHTTPHandlers) *AllHTTPHandlers {
	return &AllHTTPHandlers{Order: interfaces.Order}
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

	router.Route("/leeta", func(r chi.Router) {
		r.Handle("/swagger/*", httpSwagger.WrapHandler)
		r.Mount("/order", orderRouter)
	})

	return router, tokenHandler, nil
}

func buildOrderEndpoints(order interfaces.HTTPHandler, tokenHandler *library.TokenHandler) http.Handler {
	router := chi.NewRouter()
	router.Use(tokenHandler.ValidateMiddleware)
	router.Post("/make_order", order.CreateOrder)
	return router
}
