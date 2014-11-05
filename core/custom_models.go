package core

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
	key *datastore.Key
	
	Name    string
	Created time.Time	
	Project_id *datastore.Key// `datastore:"Project_id,noindex"`
}


func (s *Custom) toMessage(msg *ResponseMsgCustom) *ResponseMsgCustom {
	if msg == nil {
		msg = &ResponseMsgCustom{}
	}
	msg.Id = s.key.IntID()
	msg.Name = s.Name

	return msg
}



func (s *Custom) timestamp() string {
	return s.Created.Format(TIME_LAYOUT)	
}


func (s *Custom) put(c appengine.Context) (err error) {

	key := s.key
	if key == nil {
		key = datastore.NewKey(c, CUSTOM_KIND, "", 0,  nil)
	}
	key, err = datastore.Put(c, key, s)
	if err != nil {
		return err
	}
	s.key = key
	
	return nil
}




// fetchProjects runs Query q and returns Project entities fetched from the
// Datastore.
func fetchCustoms(c appengine.Context, project_id *datastore.Key, limit int) ([]*Custom, error) {

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
		item.key = keys[i]
	}

	return results, nil
}




func (s *Custom) Get(c appengine.Context, id int64) (*datastore.Key, error){

	key := datastore.NewKey(c, "Custom", "", id, nil)

	if err := datastore.Get(c, key, s); err != nil {
		return nil,  err
	}

	return key, nil
}
