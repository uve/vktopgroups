package core

import (
	"errors"
	//"fmt"
	//"math/rand"
	"net/http"
	//"time"

	"appengine/user"

	"github.com/crhym3/go-endpoints/endpoints"
	
	//"appengine/datastore"

	"config"
)


var ClientId = config.Config.OAuthProviders.Google.ClientID
var RootUrl  = config.Config.RootUrl


var (
	scopes    = []string{endpoints.EmailScope}
	clientIds = []string{ClientId, endpoints.ApiExplorerClientId}
	// in case we'll want to use Mindale API from an Android app
	audiences = []string{ClientId}
)

type BoardMsg struct {
	State string `json:"state" endpoints:"required"`
}



type ProjectReqMsg struct {
	Name string `json:"name" endpoints:"required"`
}

type ProjectRespMsg struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	//Created string `json:"created"`
	/*Outcome string `json:"outcome"`
	Played  string `json:"played"`*/
}

type ProjectsListReq struct {
	Limit int `json:"limit"`
}

type ProjectsListResp struct {
	Items []*ProjectRespMsg `json:"items"`
}

// Mindale API service
type ServiceApi struct {
}

/*
// BoardGetMove simulates a computer move in mindale.
// Exposed as API endpoint
func (ttt *ServiceApi) BoardGetMove(r *http.Request,
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

// ProjectsList queries scores for the current user.
// Exposed as API endpoint
func (ttt *ServiceApi) ProjectsList(r *http.Request,
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
	scores, err := fetchProjects(c, q, req.Limit)
	if err != nil {
		return err
	}
	resp.Items = make([]*ProjectRespMsg, len(scores))
	for i, score := range scores {
		resp.Items[i] = score.toMessage(nil)
	}
	return nil
}

*/

// ProjectsInsert inserts a new score for the current user.
func (ttt *ServiceApi) ProjectsCreate(r *http.Request,
	req *ProjectReqMsg, resp *ProjectRespMsg) error {

	c := endpoints.NewContext(r)
	u, err := getCurrentUser(c)
	if err != nil {
		return err
	}
	project := newProject(req.Name, u)
	if err := project.put(c); err != nil {
		return err
	}
	project.toMessage(resp)

	return nil
}






// ProjectsList queries scores for the current user.
// Exposed as API endpoint
func (ttt *ServiceApi) ProjectsList(r *http.Request,
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
	results, err := fetchProjects(c, q, req.Limit)
	if err != nil {
		return err
	}
	
	resp.Items = make([]*ProjectRespMsg, len(results))
	for i, item := range results {
		resp.Items[i] = item.toMessage(nil)
	}
	return nil
}




// getCurrentUser retrieves a user associated with the request.
// If there's no user (e.g. no auth info present in the request) returns
// an "unauthorized" error.
func getCurrentUser(c endpoints.Context) (*user.User, error) {

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



// RegisterService exposes ServiceApi methods as API endpoints.
// 
// The registration/initialization during startup is not performed here but
// in app package. It is separated from this package so that the
// service and its methods defined here can be used in another app,
// e.g. http://github.com/crhym3/go-endpoints.appspot.com.
func RegisterService() (*endpoints.RpcService, error) {
	api := &ServiceApi{}
	rpcService, err := endpoints.RegisterService(api,
		"vktopgroups", "v1", "VK Top groups API", true)
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


/*
	info = rpcService.MethodByName("ProjectsInsert").Info()
	info.Path, info.HttpMethod, info.Name = "scores", "POST", "projects.insert"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

	
	info := rpcService.MethodByName("PaymentsAdd").Info()
	info.Path, info.HttpMethod, info.Name = "payments", "GET", "payments.add"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences




	info = rpcService.MethodByName("BoardGetMove").Info()
	info.Path, info.HttpMethod, info.Name = "board", "POST", "board.getmove"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

	info = rpcService.MethodByName("ProjectsList").Info()
	info.Path, info.HttpMethod, info.Name = "scores", "GET", "scores.list"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

	info = rpcService.MethodByName("ProjectsInsert").Info()
	info.Path, info.HttpMethod, info.Name = "scores", "POST", "scores.insert"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

*/
	return rpcService, nil
}


