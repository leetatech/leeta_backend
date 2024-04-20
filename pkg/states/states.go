package states

import (
	"context"
	"encoding/json"
	"github.com/leetatech/leeta_backend/pkg/restclient"
	"github.com/leetatech/leeta_backend/services/models"
	"net/http"
)

const (
	getStatePath = "/states/"
)

func GetState(ctx context.Context, name, url string) (state models.State, err error) {
	getStateUrl := url + getStatePath + name
	resp, err := restclient.DoHTTPRequest(ctx, http.MethodGet, nil, getStateUrl)
	if err != nil {
		return state, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&state)
	if err != nil {
		return state, err
	}

	return
}

func GetAllStates(ctx context.Context, url string) (stateList []models.State, err error) {
	listStateUrl := url + getStatePath
	resp, err := restclient.DoHTTPRequest(ctx, http.MethodGet, nil, listStateUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&stateList)
	if err != nil {
		return nil, err
	}

	return stateList, nil
}
