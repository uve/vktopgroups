package controller

import (
	"net/http"
	"github.com/crhym3/go-endpoints/endpoints"


	"appengine"

	"appengine/taskqueue"
	"net/url"
)



func ContactsFetch(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	cursor := r.FormValue("cursor")
	c.Infof("Cursor: %v",  cursor)

	results, err := fetchContacts(c, cursor)

	if err != nil {
		if err != CURSOR_COMPLETE {
			c.Infof("Error: %v", err)
		}
		return
	}

	v := url.Values{}
	v.Add("cursor", results)

	task := taskqueue.NewPOSTTask("/fetch/contacts", v)
	_, err = taskqueue.Add(c, task, "")

}



func (api *ServiceApi) ContactsList(r *http.Request, req *GroupsListReq, resp *GroupsListResp) error {

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




