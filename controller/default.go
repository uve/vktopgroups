package controller

import (
	"time"

	"appengine"
	"appengine/datastore"

	"model"
)

const (
	TIME_LAYOUT = "02.01Jan 2, 2006 15:04:05 AM"

)

type Default struct {

	key *datastore.Key `datastore:"-"`


	Created time.Time
}


type DefaultService interface {

	Key(*datastore.Key, error)
	SetKey(*datastore.Key) error
}



func (s *Default) put(src interface {}, c appengine.Context) (err error) {
	key := s.key

	//c.Infof("PUT: %s", model.GetKind(src))

	if key == nil {
		key = datastore.NewKey(c, model.GetKind(src), "", 0,  nil)
	}

	c.Infof("PUT: %s", key)

	key, err = datastore.Put(c, key, s)

	if err != nil {
		return err
	}

	s.key = key

	return nil
}


func (s *Default) get(c appengine.Context, id int64) (*datastore.Key, error){

	key := datastore.NewKey(c, model.GetKind(s), "", id, nil)

	if err := datastore.Get(c, key, s); err != nil {
		return nil, err
	}

	return key, nil
}



// timestamp formats date/time of the project.
func (s *Default) timestamp() string {
	return s.Created.Format(TIME_LAYOUT)	
}

func NewDefault() (s *Default) {
	return &Default{Created: time.Now()}
}
