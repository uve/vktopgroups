package controller

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



func (api *ServiceApi) ProjectsCreate(r *http.Request,	req *ProjectReqMsg, resp *ProjectRespMsg) error {

	c := endpoints.NewContext(r)

	u, err := GetCurrentUser(c)
	if err != nil {
		return err
	}

	project := NewProject(req.Name, userId(u))
	
	_, err = project.put(c)
	if err != nil {
		return err
	}

	project.toMessage(resp)

	return nil
}



func (api *ServiceApi) ProjectsList(r *http.Request, req *ProjectsListReq, resp *ProjectsListResp) error {

	c := endpoints.NewContext(r)
	u, err := GetCurrentUser(c)
	if err != nil {
		return err
	}
	q := newUserProjectQuery(u)
	if req.Limit <= 0 {
		req.Limit = 10
	}
	results, err := fetch(c, q, req.Limit)
	if err != nil {
		return err
	}
	
	resp.Items = make([]*ProjectRespMsg, len(results))
	for i, item := range results {
		resp.Items[i] = item.toMessage(nil)
	}

	return nil
}






