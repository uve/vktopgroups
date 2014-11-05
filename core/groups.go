package core

import (
	"net/http"
	"net/url"

	"fmt"

	"appengine/urlfetch"
	"github.com/crhym3/go-endpoints/endpoints"

	"encoding/json"


	"models"

)


const (

	GROUPS_LIMIT=3	
	MAX_LIMIT=5000
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




// GroupsList queries scores for the current user.
// Exposed as API endpoint
func (api *ServiceApi) GroupsFetch(r *http.Request, req *GroupsFetchReq, results *GroupsFetchResp) error {

	c := endpoints.NewContext(r)

	server := "https://api.vk.com/method"


	var custom models.Custom

	custom_key, _ := custom.Get(c, req.Custom_id)


	Project_id := custom.Project_id

	c.Infof("Custom_id: %v",  custom_key)
	c.Infof("Project_id: %v", Project_id)

	c.Infof("Seach id: %v",    req.Custom_id)
	c.Infof("Seach query: %v", custom.Name)

	method := "execute.scan_groups"


	limit := fmt.Sprintf("%d", req.Limit);


	v := url.Values{}
	v.Set("access_token", VK.Token)
	v.Add("limit", limit)

	v.Add("v", VK.ApiVersion)
	v.Add("q", custom.Name)

	api_url := fmt.Sprintf("%s/%s?%s", server, method, v.Encode())


	client := urlfetch.Client(c)
	resp, err := client.Get(api_url)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}


	err = models.DeleteAll(c, Group{})
	err = models.DeleteAll(c, Contact{})



	var m Message

	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		panic(fmt.Sprintf("Could not parse : %s", err))
	}


	all_groups, err := fetchGroups(c, Project_id, MAX_LIMIT)
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
		

		group.Custom_id = custom_key;
		group.Project_id = Project_id;		


		group.City_id    = group.City.Id
		group.City_title = group.City.Title

		group.Country_id    = group.Country.Id
		group.Country_title = group.Country.Title

		new_groups = append(new_groups, group)
	}

	c.Infof("finded groups: %v", len(m.Response))

	c.Infof("new groups: %v", len(new_groups))





	group_keys, err := models.PutMulti(c, new_groups);

	if err != nil {
		panic(fmt.Sprintf("Could put to database : %s", err))
	}

	new_contacts := []Contact{}


	for i := 0; i < len(group_keys); i++ {

		value := new_groups[i]

		for j := 0; j < len(value.Contacts); j++ {

			item := value.Contacts[j]

			item.Group_id   = group_keys[i]
			item.Project_id = value.Project_id
			item.Custom_id  = value.Custom_id

			new_contacts = append(new_contacts, item)					
		}
		
	}


	if _, err := models.PutMulti(c, new_contacts); err != nil {

		c.Infof("error: %v", err)
		return err
	}


	c.Infof("new contacts: %v", len(new_contacts))

	return nil
}






