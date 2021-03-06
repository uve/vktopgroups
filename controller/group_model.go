package controller

import (
	"appengine"
	"appengine/datastore"

)


const (

	GROUP_KIND  = "Group"
)




type GroupBase struct {

	Default

	VK_id	       int64  `json:"id"`

	Name 		   string `json:"name"`
	Screen_name	   string `json:"screen_name"`

	Is_closed	   int    `json:"is_closed"`
	Type		   string `json:"type"`

	Members_count  int64  `json:"members_count"`

	Site		   string `json:"site"`
	Photo_50	   string `json:"photo_50"`
	Photo_100	   string `json:"photo_100"`
	Photo_200	   string `json:"photo_200"`


	City_id 	  int64 `json:"city_id"`
	City_title    string `json:"city_title"`

	Country_id 	  int64 `json:"country_id"`
	Country_title string `json:"country_title"`
}


type GroupCity struct {

	Id 	  int64  `json:"id"`
	Title string `json:"title"`
}

type GroupCountry struct {

	Id 	  int64  `json:"id"`
	Title string `json:"title"`
}


type Group struct {

	GroupBase

	City           GroupCity    `datastore:"-" json:"city"`          
	Country        GroupCountry `datastore:"-" json:"country"`          

	Contacts 	   []Contact `datastore:"-" json:"contacts"`

	/*key 		   *datastore.Key*/

	Project_id     *datastore.Key// `datastore:"Project_id,noindex"`
	Custom_id      *datastore.Key// `datastore:"Custom_id,noindex"`
}



type GroupJson struct {

	GroupBase

	Id      	   int64  	  `json:"id"`

	Contacts 	   []Contact `json:"contacts"`
}



func (s *Group) toMessage(msg *GroupJson) *GroupJson {
	if msg == nil {
		msg = &GroupJson{}
	}
	msg.Id = s.key.IntID()
	msg.Name = s.Name
	msg.GroupBase = s.GroupBase

	return msg
}




func fetchGroups(c appengine.Context, project_id *datastore.Key, limit int) ([]*Group, error) {

	q:= datastore.NewQuery(GROUP_KIND).Filter("Project_id=", project_id).Order("-Members_count").Limit(limit)

	results := make([]*Group, 0, QUERY_MAX)
	keys, err := q.GetAll(c, &results)
	if err != nil {
		return nil, err
	}

	for i, item := range results {
		item.SetKey(keys[i])
	}

	return results, nil
}


func queryGroupsByCustom(c appengine.Context, custom_id *datastore.Key) (*datastore.Query){

	return datastore.NewQuery(GROUP_KIND).Filter("Custom_id=", custom_id).Order("-Members_count")
}
