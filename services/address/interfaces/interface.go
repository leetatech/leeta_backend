package interfaces

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/services/address/application"
	"github.com/leetatech/leeta_backend/services/models"
	"net/http"
)

type AddressHttpHandler struct {
	AddressApplication application.AddressApplication
}

func NewAddressHttpHandler(addressApplication application.AddressApplication) *AddressHttpHandler {
	return &AddressHttpHandler{
		AddressApplication: addressApplication,
	}
}

// SaveAddressHandler is the endpoint to save address
// @Summary Save address in the DB
// @Description The endpoint ensures that all 36 states and their LGAs are saved in the DB
// @Tags Address
// @Accept json
// @Produce json
// @Param			name	path		string	true	"name"
// @Security BearerToken
// @Success 201
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /address/ [post]
func (handler *AddressHttpHandler) SaveAddressHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	err := handler.AddressApplication.Save(r.Context(), name)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
		return
	}
	pkg.EncodeResult(w, nil, http.StatusCreated)
}

// UpdateAddress is the endpoint to update address
// @Summary Update address in the DB
// @Description The endpoint is used to update an address in the DB
// @Tags Address
// @Accept json
// @Produce json
// @Param models.State body models.State true "update state request body"
// @Security BearerToken
// @Success 202
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /address/ [put]
func (handler *AddressHttpHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	var request models.State
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusBadRequest, err)
		return
	}

	err = handler.AddressApplication.Update(r.Context(), request)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
		return
	}
	pkg.EncodeResult(w, nil, http.StatusOK)
}

// GetStateHandler is the endpoint to get a state.
// @Summary Get a state.
// @Description The endpoint to get a state and all its LGAs.
// @Tags Address
// @Accept json
// @Produce json
// @Param			name	path		string	true	"name"
// @Security BearerToken
// @Success 200 {object} models.State
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /address/{name} [get]
func (handler *AddressHttpHandler) GetStateHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	response, err := handler.AddressApplication.GetState(r.Context(), name)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
		return
	}

	pkg.EncodeResult(w, response, http.StatusOK)
}

// GetAllStatesHandler is the endpoint to get all states.
// @Summary Get all states.
// @Description The endpoint to get all states and all their LGAs.
// @Tags Address
// @Accept json
// @Produce json
// @Security BearerToken
// @Success 200 {object} []models.State
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /address [get]
func (handler *AddressHttpHandler) GetAllStatesHandler(w http.ResponseWriter, r *http.Request) {
	response, err := handler.AddressApplication.GetAllStates(r.Context())
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
		return
	}

	pkg.EncodeResult(w, response, http.StatusOK)
}
