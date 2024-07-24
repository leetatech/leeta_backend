package application

import (
	"context"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/config"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/pkg/idgenerator"
	"github.com/leetatech/leeta_backend/pkg/states"
	"github.com/leetatech/leeta_backend/services/models"
	"strings"
	"time"
)

type StateAppHandler struct {
	config            config.NgnStatesConfig
	idGenerator       idgenerator.Generator
	repositoryManager pkg.RepositoryManager
}

type State interface {
	UpdateStatesFromAPI(ctx context.Context) error
	StateByName(ctx context.Context, name string) (models.State, error)
	States(ctx context.Context) ([]models.State, error)
}

func New(request pkg.ApplicationContext, config config.NgnStatesConfig) State {
	return &StateAppHandler{
		idGenerator:       idgenerator.New(),
		repositoryManager: request.RepositoryManager,
		config:            config,
	}
}

func (s *StateAppHandler) UpdateStatesFromAPI(ctx context.Context) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()
	stateList, err := states.GetAllStates(ctxWithTimeout, s.config.URL)
	if err != nil {
		return errs.Body(errs.InternalError, err)
	}
	if len(stateList) == 0 {
		return errs.Body(errs.InternalError, fmt.Errorf("no ngn states found from api %v", s.config.URL))
	}
	var allStates []any

	for _, eachState := range stateList {
		updatedState, err := states.GetState(ctxWithTimeout, eachState.Id, s.config.URL)
		if err != nil {
			return errs.Body(errs.InternalError, err)
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

	err = s.repositoryManager.StatesRepository.SaveStates(ctx, allStates)
	if err != nil {
		return errs.Body(errs.DatabaseError, err)
	}

	return nil
}

func (s *StateAppHandler) StateByName(ctx context.Context, name string) (models.State, error) {
	state, err := s.repositoryManager.StatesRepository.GetState(ctx, strings.ToUpper(name))
	if err != nil {
		return models.State{}, errs.Body(errs.DatabaseError, fmt.Errorf("could not get state %s: %w", name, err))
	}

	return state, nil
}

func (s *StateAppHandler) States(ctx context.Context) ([]models.State, error) {
	allStates, err := s.repositoryManager.StatesRepository.GetAllStates(ctx)
	if err != nil {
		return nil, errs.Body(errs.DatabaseError, err)
	}

	return allStates, nil
}
