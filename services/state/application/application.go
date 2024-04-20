package application

import (
	"context"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/config"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/states"
	"github.com/leetatech/leeta_backend/services/models"
	"go.uber.org/zap"
	"strings"
	"time"
)

type StateAppHandler struct {
	config        config.NgnStatesConfig
	idGenerator   pkg.IDGenerator
	tokenHandler  pkg.TokenHandler
	logger        *zap.Logger
	allRepository pkg.Repositories
}

func (s StateAppHandler) FetchStateDataFromAPI(ctx context.Context) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()
	stateList, err := states.GetAllStates(ctxWithTimeout, s.config.URL)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.InternalError, err)
	}
	if len(stateList) == 0 {
		return leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("no ngn states found from api %v", s.config.URL))
	}
	var allStates []any

	for _, eachState := range stateList {
		updatedState, err := states.GetState(ctxWithTimeout, eachState.Id, s.config.URL)
		if err != nil {
			return leetError.ErrorResponseBody(leetError.InternalError, err)
		}

		state := models.State{
			Id:       s.idGenerator.Generate(),
			Name:     strings.ToUpper(eachState.Name),
			Region:   eachState.Region,
			Capital:  eachState.Capital,
			Lgas:     updatedState.Lgas,
			Slogan:   eachState.Slogan,
			Towns:    updatedState.Towns,
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
	FetchStateDataFromAPI(ctx context.Context) error
	GetState(ctx context.Context, name string) (models.State, error)
	GetAllStates(ctx context.Context) ([]models.State, error)
}

func NewStateApplication(request pkg.DefaultApplicationRequest, config config.NgnStatesConfig) StateApplication {
	return &StateAppHandler{
		idGenerator:   pkg.NewIDGenerator(),
		logger:        request.Logger,
		tokenHandler:  request.TokenHandler,
		allRepository: request.AllRepository,
		config:        config,
	}
}
