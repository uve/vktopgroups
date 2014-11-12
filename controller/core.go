package controller

import (
	"errors"

	"appengine/user"

	"github.com/crhym3/go-endpoints/endpoints"

	"config"

	"os"
	"fmt"
	"encoding/json"
	"strings"
)


var (
	scopes    = []string{endpoints.EmailScope}
	clientIds = []string{ClientId, endpoints.ApiExplorerClientId}
	// in case we'll want to use Mindale API from an Android app
	audiences = []string{ClientId}


	VK = config.Config.VK

	ClientId = config.Config.OAuthProviders.Google.ClientID
	RootUrl  = config.Config.RootUrl
)

const config_api  = "config/config_api.json"


// Mindale API service
type ServiceApi struct {
}




// getCurrentUser retrieves a user associated with the request.
// If there's no user (e.g. no auth info present in the request) returns
// an "unauthorized" error.
func GetCurrentUser(c endpoints.Context) (*user.User, error) {

	user, err := endpoints.CurrentUser(c, scopes, audiences, clientIds)
	if err != nil {
		c.Errorf("User not found: %s", err)
		return nil, err
	}

	if user == nil {
		return nil, errors.New("Unauthorized: Please, sign in.")
	}

	c.Debugf("Current user: %s", user)
	return user, nil
}



type Service struct{
	Method     string
	Name       string
	Path       string
	HttpMethod string
}

// RegisterService exposes ServiceApi methods as API endpoints.
// 
// The registration/initialization during startup is not performed here but
// in app package. It is separated from this package so that the
// service and its methods defined here can be used in another app,
// e.g. http://github.com/crhym3/go-endpoints.appspot.com.
func RegisterService() (*endpoints.RpcService, error) {
	api := &ServiceApi{}
	rpcService, err := endpoints.RegisterService(api, "vktopgroups", "v1", "VK Top groups API", true)
	if err != nil {
		return nil, err
	}

	configFile, err := os.Open(config_api)
	if err != nil {
		panic(fmt.Sprintf("Could not open %s: %s", config_api, err))
	}
	defer configFile.Close()


	config := make([]Service,0)

	if err = json.NewDecoder(configFile).Decode(&config); err != nil {
		panic(fmt.Sprintf("Could not parse %s: %s", config_api, err))
	}

	for _, item := range config{

		info := rpcService.MethodByName(item.Method).Info()
		info.Path = item.Path
		info.Name = strings.Replace(item.Path, "/", ".", -1)
		info.HttpMethod = item.HttpMethod
		info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences
	}

	return rpcService, nil
}

