package core

import (
	//"errors"
	//"fmt"
	//"math/rand"
	"net/http"
	//"time"

	"github.com/crhym3/go-endpoints/endpoints"

)



type RequestMsgCustom struct {
	Name       string  `json:"Name"              endpoints:'required'`
	Project_id int64   `json:"Project_id,string" endpoints:'required'`
}

type ResponseMsgCustom struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	//Created string `json:"created"`
	/*Outcome string `json:"outcome"`
	Played  string `json:"played"`*/
}

type ListRequestCustoms struct {
	Limit      int     `json:"limit"`
	Project_id int64   `json:"Project_id,string" endpoints:'required'`
}


type ListResponseCustoms struct {
	Items []*ResponseMsgCustom "json:'items'"
}



func (ttt *ServiceApi) CustomsCreate(r *http.Request, req *RequestMsgCustom, resp *ResponseMsgCustom) error {

	c := endpoints.NewContext(r)

	c.Infof("CustomsCreate Params:")
	c.Infof("Name: %v", req.Name)
	c.Infof("Project_id:%v", req.Project_id)


	project_id, err := getProjectKey(c, req.Project_id)
	if err != nil {
		return err
	}

	c.Infof("Project key:  %v", project_id)


	item := newCustom(req.Name, project_id)
	if err := item.put(c); err != nil {
		return err
	}
	item.toMessage(resp)

	return nil
}





func (ttt *ServiceApi) CustomsList(r *http.Request, req *ListRequestCustoms, resp *ListResponseCustoms) error {

	c := endpoints.NewContext(r)

	project_id, _ := getProjectKey(c, req.Project_id)
	/*
	if err != nil {
		return errz
	}
	*/

	results, err := fetchCustoms(c, project_id, req.Limit)
	if err != nil {
		return err
	}
	
	resp.Items = make([]*ResponseMsgCustom, len(results))
	for i, item := range results {
		resp.Items[i] = item.toMessage(nil)
	}

	return nil
}

