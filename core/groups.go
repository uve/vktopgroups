package core

import (
	"net/http"
	"net/url"


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
	Limit 	   int    `json:"limit"`
	Custom_id  int64  `json:"custom_id,string" endpoints:'required'`

}



type UserResponse struct {
	Id		   int `json:"id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Photo      string `json:"photo"`
}


type GroupsListResp struct {

	Items []Group `json:"items"`
}






/*
type GroupsResponse struct {
	Group
}
*/
type Message struct {

	Response []Group `json:"response"`

}


// GroupsList queries scores for the current user.
// Exposed as API endpoint
func (api *ServiceApi) GroupsList(r *http.Request, req *GroupsListReq, results *GroupsListResp) error {


	c := endpoints.NewContext(r)
	client := urlfetch.Client(c)


	server := "https://api.vk.com/method"


	custom, _ := getCustom(c, req.Custom_id)

	c.Infof("Seach id: %v", req.Custom_id)
	c.Infof("Seach query: %v", custom.Name)

	/*
	method := "users.get"
	parameters := "user_ids=1184396&fields=photo"
	*/

	method := "execute.scan_groups"


	v := url.Values{}
	v.Set("access_token", Token)
	v.Add("limit", "3")

	v.Add("v", ApiVersion)
	v.Add("q", custom.Name)


	api_url := fmt.Sprintf("%s/%s?%s", server, method, v.Encode())


	resp, err := client.Get(api_url)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}


	var m Message

	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		panic(fmt.Sprintf("Could not parse : %s", err))
	}




	/*
	for i := 0; i < len(m.Response); i++ {

		user := m.Response[i];

		c.Infof("New response!!!")

		c.Infof("Id: %v", user.Id)
		c.Infof("First_name: %v", user.First_name)
		c.Infof("Last_name: %v", user.Last_name)
		c.Infof("Photo: %v", user.Photo)

		c.Infof("Id: %v", user.Id)
		c.Infof("Name: %v", user.Name)
		c.Infof("Members_count: %v", user.Members_count)

	}
	*/

	c.Infof("Values from vk: %v", len(m.Response))

	cnt,_ := putMulti(c, &m.Response)


	c.Infof("Values created: %v", cnt)


	results.Items = m.Response

	return nil
}






