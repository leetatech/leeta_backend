package interfaces

import (
	"github.com/go-chi/chi/v5"
	_ "github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/jwtmiddleware"
	_ "github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/state/application"
	"net/http"
)

type StateHttpHandler struct {
	StateApplication application.State
}

func New(stateApplication application.State) *StateHttpHandler {
	return &StateHttpHandler{
		StateApplication: stateApplication,
	}
}

// RetrieveNGNStatesData is the endpoint responsible for calling external API to retrieve all states in NGN
// should not be listed as one of our accessible APIs. only for internal use
func (handler *StateHttpHandler) RetrieveNGNStatesData(w http.ResponseWriter, r *http.Request) {
	err := handler.StateApplication.UpdateStatesFromAPI(r.Context())
	if err != nil {
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	jwtmiddleware.WriteJSONResponse(w, nil, http.StatusCreated)
}

// GetState is the endpoint to get a state.
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
func (handler *StateHttpHandler) GetState(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	response, err := handler.StateApplication.StateByName(r.Context(), name)
	if err != nil {
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	jwtmiddleware.WriteJSONResponse(w, response, http.StatusOK)
}

// ListStates is the endpoint to get all states.
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
func (handler *StateHttpHandler) ListStates(w http.ResponseWriter, r *http.Request) {
	response, err := handler.StateApplication.States(r.Context())
	if err != nil {
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	jwtmiddleware.WriteJSONResponse(w, response, http.StatusOK)
}
