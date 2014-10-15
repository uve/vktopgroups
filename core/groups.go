package core

import (
	"net/http"


	"fmt"

	"appengine/urlfetch"
	"github.com/crhym3/go-endpoints/endpoints"

	"config"


	"encoding/json"
)


var ApiVersion = config.Config.VK.ApiVersion
var Token      = config.Config.VK.Token
var AppId      = config.Config.VK.AppId



type GroupRespMsg struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	//Created string `json:"created"`
	/*Outcome string `json:"outcome"`
	Played  string `json:"played"`*/
}

type GroupsListReq struct {
	Limit int `json:"limit"`
}

type GroupsListResp struct {
	//Items []*GroupRespMsg `json:"items"`
	Items []Response `json:"items"`
}


type Response struct {
	Id		   int `json:"id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Photo      string `json:"photo"`
}

type Message struct {

	Response []Response `json:"response"`

}


// GroupsList queries scores for the current user.
// Exposed as API endpoint
func (ttt *ServiceApi) GroupsList(r *http.Request, req *GroupsListReq, results *GroupsListResp) error {


	c := endpoints.NewContext(r)
	client := urlfetch.Client(c)


	api    := ApiVersion
	server := "https://api.vk.com/method"
	method := "users.get"
	parameters := "user_ids=1184396&fields=photo"
	access_token := Token
	//AppId


	url := fmt.Sprintf("%s/%s?access_token=%s&v=%s&%s", server, method, access_token, api, parameters)

	c.Infof("URL: %v", url)


	resp, err := client.Get(url)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}


	var m Message

	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		panic(fmt.Sprintf("Could not parse : %s", err))
	}



	for i := 0; i < len(m.Response); i++ {

		user := m.Response[i];

		c.Infof("New user!!!")
		c.Infof("Id: %v", user.Id)

		c.Infof("First_name: %v", user.First_name)
		c.Infof("Last_name: %v", user.Last_name)
		c.Infof("Photo: %v", user.Photo)

	}

	results.Items = m.Response



	//results.Body = body
	//fmt.Fprintf(w, "HTTP GET returned status %v", resp.Status)


	/*
	c := endpoints. NewContext(r)
	u, err := getCurrentUser(c)
	if err != nil {
		return err
	}
	q := newUserGroupQuery(u)
	if req.Limit <= 0 {
		req.Limit = 10
	}
	results, err := fetchGroups(c, q, req.Limit)
	if err != nil {
		return err
	}

	resp.Items = make([]*GroupRespMsg, len(results))
	for i, item := range results {
		resp.Items[i] = item.toMessage(nil)
	}
	*/
	return nil
}






