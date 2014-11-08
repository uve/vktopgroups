package controller

import (
	//"errors"
	//"fmt"
	//"math/rand"
	"net/http"

	//"time"

	"github.com/crhym3/go-endpoints/endpoints"
	
)



type RequestMsgCustom struct {
	Name       string  `json:"name"              endpoints:'required'`
	Project_id int64   `json:"project_id,string" endpoints:'required'`
}

type ResponseMsgCustom struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Created string `json:"created"`
	/*Outcome string `json:"outcome"`
	Played  string `json:"played"`*/
}

type ListRequestCustoms struct {
	Limit      int     `json:"limit"`
	Project_id int64   `json:"project_id,string" endpoints:'required'`
}


type ListResponseCustoms struct {
	Items []*ResponseMsgCustom `json:"items"`
}




func (api *ServiceApi) CustomsCreate(r *http.Request, req *RequestMsgCustom, resp *ResponseMsgCustom) error {

	c := endpoints.NewContext(r)

	var project Project

	project_id, err := project.Get(c, req.Project_id)
	if err != nil {
		return err
	}

	item := NewCustom(req.Name, project_id)

	_, err = item.Put(c)
	if err != nil {
		return err
	}

	item.toMessage(resp)

	return nil
}



func (s *Custom) toMessage(msg *ResponseMsgCustom) *ResponseMsgCustom {
	if msg == nil {
		msg = &ResponseMsgCustom{}
	}
	
	msg.Name = s.Name

	msg.Id 		= s.Id()
	msg.Created = s.GetCreated()

	return msg
}




func (api *ServiceApi) CustomsList(r *http.Request, req *ListRequestCustoms, resp *ListResponseCustoms) error {

	c := endpoints.NewContext(r)

	c.Infof("New List query")


	var project Project

	key, err := project.Get(c, req.Project_id)
	if err != nil {
		c.Errorf("Error: %v", err)
		return err
	}

	c.Infof("Project_id: %v", req.Project_id)
	c.Infof("Project key: %v", key)


	q := queryCustomByProject(key)

	results, err := fetchCustoms(c, q, req.Limit)
	if err != nil {
		return err
	}
	
	
	resp.Items = make([]*ResponseMsgCustom, len(results))
	for i, item := range results {
		resp.Items[i] = item.toMessage(nil)
		c.Infof("%v", i)
		c.Infof("%v",item)
		
	}
	

	return nil
}


