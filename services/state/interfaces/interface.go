package interfaces

import (
	"github.com/go-chi/chi/v5"
	"github.com/leetatech/leeta_backend/pkg"
	_ "github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/state/application"
	"net/http"
)

type StateHttpHandler struct {
	StateApplication application.StateApplication
}

func NewStateHttpHandler(stateApplication application.StateApplication) *StateHttpHandler {
	return &StateHttpHandler{
		StateApplication: stateApplication,
	}
}

// SaveStateHandler is the endpoint to save state
// @Summary Save state in the DB
// @Description The endpoint ensures that all 36 states and their LGAs are saved in the DB
// @Tags State
// @Accept json
// @Produce json
// @Security BearerToken
// @Success 201
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /state/ [post]
func (handler *StateHttpHandler) SaveStateHandler(w http.ResponseWriter, r *http.Request) {
	err := handler.StateApplication.Save(r.Context())
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
		return
	}
	pkg.EncodeResult(w, nil, http.StatusCreated)
}

// GetStateHandler is the endpoint to get a state.
// @Summary Get a state.
// @Description The endpoint to get a state and all its LGAs.
// @Tags state
// @Accept json
// @Produce json
// @Param			name	path		string	true	"name"
// @Security BearerToken
// @Success 200 {object} models.State
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /state/{name} [get]
func (handler *StateHttpHandler) GetStateHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	response, err := handler.StateApplication.GetState(r.Context(), name)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
		return
	}

	pkg.EncodeResult(w, response, http.StatusOK)
}

// GetAllStatesHandler is the endpoint to get all states.
// @Summary Get all states.
// @Description The endpoint to get all states and all their LGAs.
// @Tags state
// @Accept json
// @Produce json
// @Security BearerToken
// @Success 200 {object} []models.State
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /state [get]
func (handler *StateHttpHandler) GetAllStatesHandler(w http.ResponseWriter, r *http.Request) {
	response, err := handler.StateApplication.GetAllStates(r.Context())
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
		return
	}

	pkg.EncodeResult(w, response, http.StatusOK)
}
