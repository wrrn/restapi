package configuration

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/warrenharper/restapi/utils/request"
	"github.com/warrenharper/restapi/utils/response"
)

type ConfigurationHandler struct {
	ConfigurationController
}

func (ch ConfigurationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// if _, err := authentication.CheckSession(r); err != nil {
	// 	auth.Forbidden(w)
	// 	return
	// }

	path := r.URL.Path
	variables := request.GetURLVariables(path)
	log.Println(variables)
	log.Println("PATH", path)
	switch {
	case request.Is(r, "GET") && path == "/":
		ch.handleGetAll(w, r)
	case request.Is(r, "GET") && len(variables) == 1:
		ch.handleGet(w, r, variables)
	case request.Is(r, "POST") && path == "/":
		ch.handleAdd(w, r)
	case request.Is(r, "DELETE") && len(variables) == 1:
	case request.Is(r, "PATCH") && len(variables) == 1:
	default:
		fmt.Fprintf(w, "PATH: %s", path)

	}

}

func (ch ConfigurationHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	configs, err := ch.GetAll()
	if err != nil {
		response.ServerError(w)
		return
	}
	response.WriteJson(w, http.StatusOK, Configurations{configs})
}

func (ch ConfigurationHandler) handleGet(w http.ResponseWriter, r *http.Request, vars []string) {
	configName := vars[0]
	configs, err := ch.Get(configName)
	if err == DoesNotExistErr {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	if err != nil {
		response.ServerError(w)
		return
	}

	response.WriteJson(w, http.StatusOK, Configurations{configs})
}

func (ch ConfigurationHandler) handleAdd(w http.ResponseWriter, r *http.Request) {
	config := Configuration{}
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, "Bad Format", http.StatusBadRequest)
		return
	}

	configs, err := ch.Add(config)
	if configErr, ok := err.(ConfigurationError); ok && configErr.Err == DuplicateConfigErr {
		response.WriteJson(w, http.StatusConflict, Configurations{[]Configuration{configErr.Configuration}})
		return
	}

	if err != nil {
		response.ServerError(w)
	}

	response.WriteJson(w, http.StatusOK, Configurations{configs})

}
