package controller

import (
	//"errors"
	//"fmt"
	//"math/rand"
	"net/http"

	//"time"

	"github.com/crhym3/go-endpoints/endpoints"

	"model"
	

	"time"
	//"errors"
)



type RequestMsgCustom struct {
	Name       string  `json:"name"              endpoints:'required'`
	Project_id int64   `json:"project_id,string" endpoints:'required'`
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
	Project_id int64   `json:"project_id,string" endpoints:'required'`
}


type ListResponseCustoms struct {
	Items []*ResponseMsgCustom `json:"items"`
}




func (api *ServiceApi) CustomsCreate(r *http.Request, req *RequestMsgCustom, resp *ResponseMsgCustom) error {

	c := endpoints.NewContext(r)

	var project Project

	project_id, err := project.get(c, req.Project_id)
	if err != nil {
		return err
	}

	item := &Custom{

		Name: req.Name,
		Created: time.Now(),
		Project_id: project_id,
	}

	key, err := model.Put(c, item)
	if err != nil {
		return err
	}

	item.Key = key

	item.toMessage(resp)

	return nil
}



func (s *Custom) toMessage(msg *ResponseMsgCustom) *ResponseMsgCustom {
	if msg == nil {
		msg = &ResponseMsgCustom{}
	}
	msg.Id = s.Key.IntID()
	msg.Name = s.Name

	return msg
}




func (api *ServiceApi) CustomsList(r *http.Request, req *ListRequestCustoms, resp *ListResponseCustoms) error {

	c := endpoints.NewContext(r)

	c.Infof("New List query")


	var project Project

	project_id, err := project.get(c, req.Project_id)
	if err != nil {
		return err
	}


	results, err := FetchCustoms(c, project_id, req.Limit)
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
