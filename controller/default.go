package controller

import (
	"time"

	"model"
	"appengine"
	"appengine/datastore"

)

const (
	//TIME_LAYOUT = "02.01Jan 2, 2006 15:04:05 AM"
	TIME_LAYOUT = "Jan 2, 2006 at 3:04pm (MST)"

)


type Default struct {

	key *datastore.Key `datastore:"-"`


	Created time.Time
}




type DefaultInterface interface {

	Key() *datastore.Key
	SetKey(*datastore.Key)
	Id() int64
}




func (s *Default) Key() *datastore.Key {
	return s.key
}


/*
func (src *Default) GetKey(c appengine.Context, id int64) (*datastore.Key, error) {

	return model.GetKey(c, src, id)
}
*/



func (s *Default) Id() int64 {

	
	if s.Key() == nil{
		return 0
	}
	
	return s.Key().IntID()
}


func (s *Default) SetKey(key *datastore.Key) {
	s.key = key
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

func NewDefault() (s Default) {
	return Default{
				Created: time.Now(),
			}
}

// timestamp formats date/time of the project.
func (s *Default) GetCreated() string {
	return s.Created.Format(TIME_LAYOUT)	
}