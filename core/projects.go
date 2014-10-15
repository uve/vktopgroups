package core

import (
	"net/http"
	"github.com/crhym3/go-endpoints/endpoints"
)



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






