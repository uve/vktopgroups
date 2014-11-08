package controller

import (
	"time"

	"model"
	"appengine"
	"appengine/datastore"

)

const (
	TIME_LAYOUT = "02.01Jan 2, 2006 15:04:05 AM"

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


func (s *Default) Id() int64 {

	/*
	if s.Key() == nil{
		return 0
	}
	*/
	return s.Key().IntID()
}


func (s *Default) setKey(key *datastore.Key) {
	s.key = key
}



func Put(c appengine.Context, src *DefaultInterface) (*datastore.Key, error) {


	key := src.Key()


	if key == nil {
		key = datastore.NewKey(c, model.GetKind(src), "", 0,  nil)
	}


	c.Infof("PUT: %s", key)

	key, err = datastore.Put(c, key, src)

	if err != nil {
		return err, nil
	}

	src.SetKey(key)

	return key, nil
}


/*

func (src *Project) put(c appengine.Context) (*datastore.Key, error) {

	key, err := model.Put(c, src)
	if err != nil{
		return nil, err
	}

	src.setKey(key)

	return key, nil
}


func (s *Default) put(src interface {}, c appengine.Context) (err error) {
	key := s.key

	//c.Infof("PUT: %s", model.GetKind(src))

	if key == nil {
		key = datastore.NewKey(c, model.GetKind(src), "", 0,  nil)
	}

	c.Infof("PUT: %s", key)

	key, err = datastore.Put(c, key, src)

	if err != nil {
		return err
	}

	s.key = key

	return nil
}
*/



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
	return Default{Created: time.Now()}
}
