package confighandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/wrrn/restapi/configuration"
	"github.com/wrrn/restapi/configuration/configsort"
	"github.com/wrrn/restapi/utils/request"
	"github.com/wrrn/restapi/utils/response"
)

type Handler struct {
	configuration.ConfigurationController
}

func (ch Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	variables := request.GetURLVariables(path)
	switch {
	case request.Is(r, "GET") && path == "/":
		ch.handleGetAll(w, r)
	case request.Is(r, "GET") && len(variables) == 1:
		ch.handleGet(w, r, variables[0])
	case request.Is(r, "POST") && path == "/":
		ch.handleAdd(w, r)
	case request.Is(r, "DELETE") && len(variables) == 1:
		ch.handleDelete(w, r, variables[0])
	case request.Is(r, "PATCH") && len(variables) == 1:
		ch.handleModify(w, r, variables[0])
	default:
		http.Error(w, "", http.StatusNotImplemented)

	}

}

// handleGetAll sends a list of all the configurations with a 200 code
func (ch Handler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	configs, err := ch.GetAll()
	if err != nil {
		response.ServerError(w)
		return
	}

	configs, err = handleParameters(w, r, configs)
	if err != nil {
		http.Error(w, "Bad Query String", http.StatusBadRequest)
	}

	response.WriteJson(w, http.StatusOK, configuration.Configurations{configs})
}

// handleGet sends a list of configurations containing only one configuration
// whose name matches the name specified in the url with a 200 code. If no such
// configuration can be found sends a 404 code.
func (ch Handler) handleGet(w http.ResponseWriter, r *http.Request, configName string) {
	configs, err := ch.Get(configName)
	if err == configuration.DoesNotExistErr {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	if err != nil {
		response.ServerError(w)
		return
	}

	response.WriteJson(w, http.StatusOK, configuration.Configurations{configs})
}

// handleAdd parses the json in the request body and creates a configuration with the fields
// indicated in the json. If successful it sends a 200 code. If two configurations
//have the same name then it sends a 409 code with the configuration in the body of
// the response.
func (ch Handler) handleAdd(w http.ResponseWriter, r *http.Request) {
	config := configuration.Configuration{}
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, "Bad Format", http.StatusBadRequest)
		return
	}

	configs, err := ch.Add(config)
	if configErr, ok := err.(configuration.Error); ok && configErr.Err == configuration.DuplicateConfigErr {
		response.WriteJson(w, http.StatusConflict, configuration.Configurations{[]configuration.Configuration{configErr.Configuration}})
		return
	}

	if err != nil {
		response.ServerError(w)
	}

	response.WriteJson(w, http.StatusOK, configuration.Configurations{configs})

}

// handleDelete deletes the configuration whose name matches the name specified
// in the url. If no such configuration exists do nothing.
// Always sends a 204 code.
func (ch Handler) handleDelete(w http.ResponseWriter, r *http.Request, configName string) {
	if err := ch.Delete(configName); err != nil {
		response.ServerError(w)
		return
	}
	response.Write(w, http.StatusNoContent, nil)
}

// handleModify modifies the configuration whose name matches the name specified
// in the url. If no such configuration exists sends a 404 code. If the modification
// would cause two configurations to have the same name then sends a 409 code with
// the configuration in the body of the response. If successful sends a 200 code/
func (ch Handler) handleModify(w http.ResponseWriter, r *http.Request, configName string) {
	config := configuration.Configuration{}
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, "Bad Format", http.StatusBadRequest)
		return
	}

	config, err = ch.Modify(configName, config)

	if err == configuration.DoesNotExistErr {
		http.Error(w, "", http.StatusNotFound)
		return
	} else if confErr, ok := err.(configuration.Error); ok && confErr.Err == configuration.DuplicateConfigErr {
		response.WriteJson(w, http.StatusConflict, configuration.Configurations{[]configuration.Configuration{confErr.Configuration}})
		return
	} else if err != nil {
		response.ServerError(w)
		return
	}

	response.WriteJson(w, http.StatusOK, configuration.Configurations{[]configuration.Configuration{config}})
}

// handleParameters loops through the parameters of request and performs actions
// based on their values.
func handleParameters(w http.ResponseWriter, r *http.Request, configs []configuration.Configuration) ([]configuration.Configuration, error) {
	configs = handleSortParameters(r, configs)
	return handlePaginateParameters(r, configs)

}

func handleSortParameters(r *http.Request, configs []configuration.Configuration) (sortedConfigs []configuration.Configuration) {
	sortBy := r.FormValue("sort")
	sorter := &configsort.Sorter{}
	var orderer configsort.Orderer
	switch sortBy {
	case "name":
		orderer = sorter.ByName
	case "hostname":
		orderer = sorter.ByHostName
	case "port":
		orderer = sorter.ByPort
	case "username":
		orderer = sorter.ByUsername
	}

	if orderer != nil {
		configs = sorter.Sort(orderer, configs)
	}

	return configs
}

func handlePaginateParameters(r *http.Request, configs []configuration.Configuration) (page []configuration.Configuration, err error) {
	pageNum, pnErr := strconv.Atoi(r.FormValue("page"))
	perPage, ppErr := strconv.Atoi(r.FormValue("per_page"))

	switch {
	case (pnErr != nil && ppErr == nil) || (pnErr == nil && ppErr != nil):
		return configs, fmt.Errorf("Bad page parameters")
	case (pnErr == nil && ppErr == nil) && perPage < 100:
		return configuration.GetPage(configs, pageNum, perPage), nil
	case perPage > 100:
		return configs, fmt.Errorf("Pagination format error")
	default:
		return configs, nil
	}
}
