package core

import (
	"net/http"
	"net/url"

	"fmt"

	"appengine/urlfetch"
	"github.com/crhym3/go-endpoints/endpoints"

	"encoding/json"


)




/*
type UserResponse struct {
	Id		   int `json:"id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Photo      string `json:"photo"`
}
*/

type GroupsListReq struct {
	Limit 	   int    `json:"limit"`
	Project_id int64  `json:"project_id,string" endpoints:'required'`
}



type GroupsListResp struct {

	Items []*GroupJson `json:"items"`
}




type GroupsFetchReq struct {
	Limit 	   int    `json:"limit"`
	Custom_id  int64  `json:"custom_id,string" endpoints:'required'`

}

type GroupsFetchResp struct {

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
func (api *ServiceApi) GroupsFetch(r *http.Request, req *GroupsFetchReq, results *GroupsFetchResp) error {

	c := endpoints.NewContext(r)

	server := "https://api.vk.com/method"

	Custom_id, custom, _ := getCustom(c, req.Custom_id)
	Project_id := custom.Project_id

	c.Infof("Custom_id: %v",  Custom_id)
	c.Infof("Project_id: %v", Project_id)

	c.Infof("Seach id: %v",    req.Custom_id)
	c.Infof("Seach query: %v", custom.Name)

	method := "execute.scan_groups"


	v := url.Values{}
	v.Set("access_token", VK.Token)
	v.Add("limit", "3")

	v.Add("v", VK.ApiVersion)
	v.Add("q", custom.Name)

	api_url := fmt.Sprintf("%s/%s?%s", server, method, v.Encode())


	client := urlfetch.Client(c)
	resp, err := client.Get(api_url)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}


	var m Message

	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		panic(fmt.Sprintf("Could not parse : %s", err))
	}


	all_groups, err := fetchGroups(c, Project_id, req.Limit)
	if err != nil {
		return err
	}

	new_groups := []Group{}




	for _,group := range m.Response {

		is_exist := false

		for _,item := range all_groups {

			if group.VK_id == item.VK_id {
				is_exist = true
				break
			}
		}

		if is_exist {
			continue
		}



		group.Custom_id = Custom_id;
		group.Project_id = Project_id;

		new_groups = append(new_groups, group)
	}

	c.Infof("finded groups: %v", len(m.Response))

	c.Infof("new groups: %v", len(new_groups))





	if err = putMulti(c, &new_groups); err != nil {
		panic(fmt.Sprintf("Could put to database : %s", err))
	}

	return nil
}




func (api *ServiceApi) GroupsList(r *http.Request, req *GroupsListReq, resp *GroupsListResp) error {

	c := endpoints.NewContext(r)

	project_id, _ := getProjectKey(c, req.Project_id)


	results, err := fetchGroups(c, project_id, req.Limit)
	if err != nil {
		return err
	}

	resp.Items = make([]*GroupJson, len(results))
	for i, item := range results {
		resp.Items[i] = item.toMessage(nil)
	}

	return nil
}






