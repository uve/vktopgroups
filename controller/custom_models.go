package controller

import (

	"appengine"
	"appengine/datastore"

	"model"
	
)

const (
	CUSTOM_KIND  = "Custom"
)


type Custom struct {

	Default 

	Name    string
	Project_id *datastore.Key// `datastore:"Project_id,noindex"`
}





func NewCustom(name string, project *datastore.Key) (*Custom){

	return &Custom{
		Name: name,
		Project_id: project,
		Default: NewDefault(),
	}
}


func (src *Custom) Put(c appengine.Context) (*datastore.Key, error) {

	key, err := model.Put(c, src)
	if err != nil {
		return nil, err
	}

	src.SetKey(key)

	return key, nil
}






func queryCustomByProject(project *datastore.Key) *datastore.Query {
	return datastore.NewQuery(CUSTOM_KIND).Filter("Project_id =", project).Order("Created")
}


// Datastore.
func fetchCustoms(c appengine.Context, q *datastore.Query, limit int) ([]*Custom, error) {

	if limit <= 0 {
		limit = 10
	}


	items := make([]*Custom, 0, limit)
	keys, err := q.Limit(limit).GetAll(c, &items)
	if err != nil {
		return nil, err
	}
	for i, item := range items {
		item.SetKey(keys[i])
	}
	return items, nil
}



func (s *Custom) Get(c appengine.Context, id int64) (*datastore.Key, error){

	c.Infof("################### id:  %v", id)

	key := datastore.NewKey(c, "Custom", "", id, nil)

	if err := datastore.Get(c, key, s); err != nil {
		return nil,  err
	}

	return key, nil
}
