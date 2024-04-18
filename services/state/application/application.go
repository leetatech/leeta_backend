package application

import (
	"context"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/states"
	"github.com/leetatech/leeta_backend/services/models"
	"go.uber.org/zap"
	"strings"
	"time"
)

type StateAppHandler struct {
	idGenerator   pkg.IDGenerator
	tokenHandler  pkg.TokenHandler
	logger        *zap.Logger
	allRepository pkg.Repositories
	stateConfig   states.StateMethods
}

func (s StateAppHandler) Save(ctx context.Context) error {
	fetchedStates, err := s.stateConfig.GetAllStates()
	if err != nil {
		return leetError.ErrorResponseBody(leetError.InternalError, err)
	}
	var allStates []interface{}

	for _, perState := range *fetchedStates {
		fetchedState, err := s.stateConfig.GetState(perState.Id)
		if err != nil {
			return leetError.ErrorResponseBody(leetError.InternalError, err)
		}

		state := models.State{
			Id:       s.idGenerator.Generate(),
			Name:     strings.ToUpper(perState.Name),
			Region:   perState.Region,
			Capital:  perState.Capital,
			Lgas:     fetchedState.Lgas,
			Slogan:   perState.Slogan,
			Towns:    fetchedState.Towns,
			StatusTs: time.Now().Unix(),
			Ts:       time.Now().Unix(),
		}

		allStates = append(allStates, state)
	}

	err = s.allRepository.StatesRepository.SaveStates(ctx, allStates)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}

func (s StateAppHandler) GetState(ctx context.Context, name string) (models.State, error) {
	state, err := s.allRepository.StatesRepository.GetState(ctx, strings.ToUpper(name))
	if err != nil {
		s.logger.Error("could not get state", zap.String("name", name), zap.Error(err))
		return models.State{}, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return state, nil
}

func (s StateAppHandler) GetAllStates(ctx context.Context) ([]models.State, error) {
	allStates, err := s.allRepository.StatesRepository.GetAllStates(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return allStates, nil
}

type StateApplication interface {
	Save(ctx context.Context) error
	GetState(ctx context.Context, name string) (models.State, error)
	GetAllStates(ctx context.Context) ([]models.State, error)
}

func NewStateApplication(request pkg.DefaultApplicationRequest) StateApplication {
	return &StateAppHandler{
		idGenerator:   pkg.NewIDGenerator(),
		logger:        request.Logger,
		tokenHandler:  request.TokenHandler,
		allRepository: request.AllRepository,
		stateConfig:   request.States,
	}
}
