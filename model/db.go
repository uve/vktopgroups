package model

import (
	"appengine"
	"appengine/datastore"

	"reflect"
	"strings"

	"errors"
)



const (

	TIME_LAYOUT = "02.01Jan 2, 2006 15:04:05 AM"
	
	QUERY_MAX = 5000
)


func getStructName(any interface{}) (string){

 	full_type := reflect.TypeOf(any).String()
 	arr := strings.Split(full_type, ".")

 	return arr[len(arr)-1]
}



func DeleteAll(c appengine.Context, any interface{}) (error) {

	kind := getStructName(any)
		
	c.Infof("Delete all: %v", kind)

	q:= datastore.NewQuery(kind).KeysOnly()


	results := make([]reflect.Value, 0, QUERY_MAX)
	keys, err := q.GetAll(c, &results)
	if err != nil {
		return err
	}


	err = datastore.DeleteMulti(c, keys)
	if err != nil {
		return err
	}
	
	return nil
}



func GetKind(src interface{}) (string){

	items := reflect.ValueOf(src)

	type_str := ""

	if items.Kind() == reflect.Slice {
		type_str = reflect.Indirect(items).Index(0).Type().String()

	} else {
		type_str = reflect.Indirect(items).Type().String()
	}


	arr := strings.Split(type_str, ".")
	kind := arr[len(arr)-1]

	return kind
}

/*
type Service interface {
  func Key()
  func SetKey(v)
}

func Key(p Project) {*datastore.Key, error){

}
*/

func Put(c appengine.Context, src interface{}) (*datastore.Key, error) {

	//key := src.key

	if reflect.ValueOf(src).Kind() != reflect.Ptr{
		return nil, errors.New("Value's type is not Ptr")
	}

	kind := GetKind(src)

	incomplete_keys := datastore.NewKey(c, kind, "", 0, nil)

	keys, err := datastore.Put(c, incomplete_keys, src)
	if err != nil {
		return nil, err
	}

	return keys, nil

}



func PutMulti(c appengine.Context, src interface{}) ([]*datastore.Key, error){


	items := reflect.ValueOf(src)

	if items.Kind() != reflect.Slice {
		return nil, errors.New("Value's type is not Slice")
	}

	cnt := items.Len()

	if cnt < 1{
		return nil, errors.New("PutMulti: slice length < 1")
	}

	kind := GetKind(src)
	//c.Infof("kind: %v", kind)


	incomplete_keys := make([]*datastore.Key, cnt)

	for i := 0; i < cnt; i++ {
		incomplete_keys[i] = datastore.NewKey(c, kind, "", 0, nil)
	}

	keys, err := datastore.PutMulti(c, incomplete_keys, src)
	if err != nil {
		return nil, err
	}


	return keys, nil

}


func GetKey(c appengine.Context, src interface {}, id int64) (*datastore.Key, error){

	kind := GetKind(src)

	key := datastore.NewKey(c, kind, "", id, nil)

	if err := datastore.Get(c, key, src); err != nil {
		return nil, err
	}

	return key, nil
}

