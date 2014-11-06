package controller

import (
	"errors"

	"appengine/user"

	"github.com/crhym3/go-endpoints/endpoints"

	"config"

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



type BoardMsg struct {
	State string `json:"state" endpoints:"required"`
}



// Mindale API service
type ServiceApi struct {
}

/*
// BoardGetMove simulates a computer move in mindale.
// Exposed as API endpoint
func (api *ServiceApi) BoardGetMove(r *http.Request,
	req *BoardMsg, resp *BoardMsg) error {

	const boardLen = 9
	if len(req.State) != boardLen {
		return fmt.Errorf("Bad Request: Invalid board: %q", req.State)
	}
	runes := []rune(req.State)
	freeIndices := make([]int, 0)
	for pos, r := range runes {
		if r != 'O' && r != 'X' && r != '-' {
			return fmt.Errorf("Bad Request: Invalid rune: %q", r)
		}
		if r == '-' {
			freeIndices = append(freeIndices, pos)
		}
	}
	freeIdxLen := len(freeIndices)
	if freeIdxLen > 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomIdx := r.Intn(freeIdxLen)
		runes[freeIndices[randomIdx]] = 'O'
		resp.State = string(runes)
	} else {
		return fmt.Errorf("Bad Request: This board is full: %q", req.State)
	}
	return nil
}

// ProjectsList queries scontrollers for the current user.
// Exposed as API endpoint
func (api *ServiceApi) ProjectsList(r *http.Request,
	req *ProjectsListReq, resp *ProjectsListResp) error {

	c := endpoints.NewContext(r)
	u, err := getCurrentUser(c)
	if err != nil {
		return err
	}
	q := newUserProjectQuery(u)
	if req.Limit <= 0 {
		req.Limit = 10
	}
	scontrollers, err := fetchProjects(c, q, req.Limit)
	if err != nil {
		return err
	}
	resp.Items = make([]*ProjectRespMsg, len(scontrollers))
	for i, scontroller := range scontrollers {
		resp.Items[i] = scontroller.toMessage(nil)
	}
	return nil
}

*/


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


func createMethod(rpcService *endpoints.RpcService, service, path, method, name string){

	info := rpcService.MethodByName(name).Info()
	info.Path, info.HttpMethod, info.Name = path, method, name
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences
}

/*
type Entity interface {

	UnmarshalHTTP(*http.Request) error
}

func GetEntity(r *http.Request, v Entity) error {
	return v.UnmarshalHTTP(r)
}

func (api *ServiceApi) GroupsList(r *http.Request, req *GroupsListReq, resp *GroupsListResp) error {

}
*/



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




	info := rpcService.MethodByName("ProjectsList").Info()
	info.Path, info.HttpMethod, info.Name = "projects/list", "GET", "projects.list"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences


	info = rpcService.MethodByName("ProjectsCreate").Info()
	info.Path, info.HttpMethod, info.Name = "projects/create", "POST", "projects.create"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences


	info = rpcService.MethodByName("CustomsList").Info()
	info.Path, info.HttpMethod, info.Name = "params/custom/list", "GET", "params.custom.list"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

	info = rpcService.MethodByName("CustomsCreate").Info()
	info.Path, info.HttpMethod, info.Name = "params/custom/create", "POST", "params.custom.create"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences



	info = rpcService.MethodByName("GroupsList").Info()
	info.Path, info.HttpMethod, info.Name = "groups/list", "GET", "groups.list"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences


	info = rpcService.MethodByName("GroupsFetch").Info()
	info.Path, info.HttpMethod, info.Name = "groups/fetch", "POST", "groups.fetch"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences


/*
	info = rpcService.MethodByName("ProjectsInsert").Info()
	info.Path, info.HttpMethod, info.Name = "scontrollers", "POST", "projects.insert"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

	
	info := rpcService.MethodByName("PaymentsAdd").Info()
	info.Path, info.HttpMethod, info.Name = "payments", "GET", "payments.add"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences




	info = rpcService.MethodByName("BoardGetMove").Info()
	info.Path, info.HttpMethod, info.Name = "board", "POST", "board.getmove"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

	info = rpcService.MethodByName("ProjectsList").Info()
	info.Path, info.HttpMethod, info.Name = "scontrollers", "GET", "scontrollers.list"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

	info = rpcService.MethodByName("ProjectsInsert").Info()
	info.Path, info.HttpMethod, info.Name = "scontrollers", "POST", "scontrollers.insert"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

*/
	return rpcService, nil
}

