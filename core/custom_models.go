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
	Project_id *datastore.Key `datastore:"Project_id,noindex"`
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
		key = datastore.NewIncompleteKey(c, CUSTOM_KIND, nil)
	}
	key, err = datastore.Put(c, key, s)
	if err == nil {
		s.key = key
	}
	return
}

// newProject returns a new Project ready to be stored in the Datastore.
func newCustom(name string, project_id *datastore.Key) *Custom {
	//return &Project{Outcome: outcome, Played: time.Now(), Player: userId(u)}

	return &Custom{
				Name: name,				
				Created: time.Now(),
				Project_id: project_id,
			}
}






// newUserProjectQuery returns a Query which can be used to list all previous
// games of a user.
func newCustomQuery(project_id *datastore.Key) *datastore.Query {
	return datastore.NewQuery(CUSTOM_KIND).Filter("Project_id =", project_id).Order("Created")
}

// fetchProjects runs Query q and returns Project entities fetched from the
// Datastore.
func fetchCustoms(c appengine.Context, q *datastore.Query, limit int) ([]*Custom, error) {

	results := make([]*Custom, 0, limit)
	keys, err := q.Limit(limit).GetAll(c, &results)
	if err != nil {
		return nil, err
	}
	for i, item := range results {
		item.key = keys[i]
	}
	return results, nil
}

