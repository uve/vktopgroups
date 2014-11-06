package controller

import (
	"time"

	"appengine"
	"appengine/datastore"
	
)

const (
	//TIME_LAYOUT = "02.01Jan 2, 2006 15:04:05 AM"
	CUSTOM_KIND  = "Custom"
)


type Custom struct {
	Key *datastore.Key 		`datastore:"-"`     
	
	Name    string
	Created time.Time	
	Project_id *datastore.Key// `datastore:"Project_id,noindex"`
}



func (s *Custom) timestamp() string {
	return s.Created.Format(TIME_LAYOUT)	
}


func (s *Custom) put(c appengine.Context) (err error) {

	key := s.Key
	if key == nil {
		key = datastore.NewKey(c, CUSTOM_KIND, "", 0,  nil)
	}
	key, err = datastore.Put(c, key, s)
	if err != nil {
		return err
	}
	s.Key = key
	
	return nil
}









// fetchProjects runs Query q and returns Project entities fetched from the
// Datastore.
func FetchCustoms(c appengine.Context, project_id *datastore.Key, limit int) ([]*Custom, error) {

	if limit<= 0 {
		limit = 10
	}

	q:= datastore.NewQuery(CUSTOM_KIND).Order("-Created").Limit(limit).Filter("Project_id=", project_id)

	results := make([]*Custom, 0, limit)
	keys, err := q.GetAll(c, &results)
	if err != nil {
		return nil, err
	}

	for i, item := range results {

		//c.Infof("Key:  %v : is_equal: %v", item.Project_id, project_id.Equal(item.Project_id))
		item.Key = keys[i]
	}

	return results, nil
}




func (s *Custom) Get(c appengine.Context, id int64) (*datastore.Key, error){

	c.Infof("################### id:  %v", id)

	key := datastore.NewKey(c, "Custom", "", id, nil)

	if err := datastore.Get(c, key, s); err != nil {
		return nil,  err
	}

	return key, nil
}
