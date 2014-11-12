package controller

import (
	"net/http"
	"net/url"

	"fmt"

	"appengine/urlfetch"
	"appengine/taskqueue"

	"github.com/crhym3/go-endpoints/endpoints"

	"encoding/json"

	"model"
	"strconv"
)


const (

	GROUPS_LIMIT = 10

)



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



type Message struct {

	Response []Group `json:"response"`
}





func (api *ServiceApi) GroupsList(r *http.Request, req *GroupsListReq, resp *GroupsListResp) error {

	c := endpoints.NewContext(r)

	var project Project

	project_id, err := project.Get(c, req.Project_id)
	if err != nil {
		return err
	}

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



func GroupsFetch(w http.ResponseWriter, r *http.Request) {

	c := endpoints.NewContext(r)

	c.Infof("Fetch groups")

	form_custom_id_str := r.FormValue("custom_id")
	form_custom_id, err := strconv.ParseInt(form_custom_id_str, 10, 64)

	limit := strconv.FormatInt(GROUPS_LIMIT, 10)

	if err != nil{
		panic(err)
	}

	var custom Custom

	custom_id, _ := custom.Get(c, form_custom_id)


	Project_id := custom.Project_id

	c.Infof("Custom_id: %v",  custom_id)
	c.Infof("Project_id: %v", Project_id)

	c.Infof("Seach id: %v",    form_custom_id)
	c.Infof("Seach query: %v", custom.Name)

	c.Infof("Seach query: %v", custom.Name)
	c.Infof("Find Group Limit: %v",  GROUPS_LIMIT)


	method := "execute.scan_groups"

	v := url.Values{}
	v.Set("access_token", VK.Token)
	v.Add("limit", limit)

	v.Add("v", VK.ApiVersion)
	v.Add("q", custom.Name)

	api_url := fmt.Sprintf("%s/%s?%s", VK.Server, method, v.Encode())


	client := urlfetch.Client(c)
	resp, err := client.Get(api_url)
	if err != nil {
		panic(err)
	}

	/*
	err = model.DeleteByModel(c, Group{})
	err = model.DeleteByModel(c, Contact{})
	*/


	var m Message

	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		panic(fmt.Sprintf("Could not parse : %s", err))
	}



	all_groups, err := fetchGroups(c, Project_id, QUERY_MAX)
	if err != nil {
		panic(err)
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

		group.Default = NewDefault()

		group.Custom_id = custom_id;
		group.Project_id = Project_id;		


		group.City_id    = group.City.Id
		group.City_title = group.City.Title

		group.Country_id    = group.Country.Id
		group.Country_title = group.Country.Title

		new_groups = append(new_groups, group)
	}

	c.Infof("finded groups: %v", len(m.Response))

	c.Infof("new groups: %v", len(new_groups))



	group_keys, err := model.PutMulti(c, new_groups);
	if err != nil {
		panic(err)
	}

	new_contacts := []Contact{}


	for i := 0; i < len(group_keys); i++ {

		value := new_groups[i]

		for j := 0; j < len(value.Contacts); j++ {

			item := value.Contacts[j]


			item.Default = NewDefault()
			item.Is_exist = false

			item.Group_id   = group_keys[i]
			item.Project_id = value.Project_id
			item.Custom_id  = value.Custom_id

			new_contacts = append(new_contacts, item)					
		}
		
	}


	if _, err := model.PutMulti(c, new_contacts); err != nil {
		panic(err)
	}


	c.Infof("new contacts: %v", len(new_contacts))

	task := taskqueue.NewPOSTTask("/fetch/contacts", url.Values{})
	_, err = taskqueue.Add(c, task, "")

}






