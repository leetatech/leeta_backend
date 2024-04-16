package application

import (
	"context"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/address"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/services/models"
	"go.uber.org/zap"
	"strings"
	"time"
)

type AddressAppHandler struct {
	idGenerator   pkg.IDGenerator
	tokenHandler  pkg.TokenHandler
	logger        *zap.Logger
	allRepository pkg.Repositories
	addressConfig address.Address
}

func (a AddressAppHandler) Save(ctx context.Context, name string) error {
	fetchedState, err := a.addressConfig.GetState(name)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.InternalError, err)
	}

	state := models.State{
		Id:       a.idGenerator.Generate(),
		Name:     strings.ToUpper(fetchedState.Name),
		Region:   fetchedState.Region,
		Capital:  fetchedState.Capital,
		Lgas:     fetchedState.Lgas,
		Slogan:   fetchedState.Slogan,
		Towns:    fetchedState.Towns,
		StatusTs: time.Now().Unix(),
		Ts:       time.Now().Unix(),
	}

	err = a.allRepository.AddressRepository.Upsert(ctx, state)
	if err != nil {
		return err
	}

	return nil
}

func (a AddressAppHandler) Update(ctx context.Context, state models.State) error {

	storedState, err := a.allRepository.AddressRepository.GetState(ctx, state.Name)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)

	}

	storedState.Lgas = state.Lgas

	err = a.allRepository.AddressRepository.Update(ctx, storedState)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}

func (a AddressAppHandler) GetState(ctx context.Context, name string) (models.State, error) {
	state, err := a.allRepository.AddressRepository.GetState(ctx, strings.ToUpper(name))
	if err != nil {
		a.logger.Error("could not get state", zap.String("name", name), zap.Error(err))
		return models.State{}, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return state, nil
}

func (a AddressAppHandler) GetAllStates(ctx context.Context) ([]models.State, error) {
	states, err := a.allRepository.AddressRepository.GetAllStates(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return states, nil
}

type AddressApplication interface {
	Save(ctx context.Context, name string) error
	Update(ctx context.Context, state models.State) error
	GetState(ctx context.Context, name string) (models.State, error)
	GetAllStates(ctx context.Context) ([]models.State, error)
}

func NewAddressApplication(request pkg.DefaultApplicationRequest) AddressApplication {
	return &AddressAppHandler{
		idGenerator:   pkg.NewIDGenerator(),
		logger:        request.Logger,
		tokenHandler:  request.TokenHandler,
		allRepository: request.AllRepository,
		addressConfig: request.Address,
	}
}
