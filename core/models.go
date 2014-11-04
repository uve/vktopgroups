package core

import (
	"appengine"
	"appengine/datastore"

	"reflect"
	"strings"
)



const (
	QUERY_MAX = 5000
)


func getStructName(any interface{}) (string){

 	full_type := reflect.TypeOf(any).String()
 	arr := strings.Split(full_type, ".")

 	return arr[len(arr)-1]
}




func deleteAll(c appengine.Context, any interface{}) (error) {

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




func PutMulti(c appengine.Context, src interface{}) ([]*datastore.Key, error){

/*
	items := reflect.ValueOf(src)	
	kind := getStructName(items)
	c.Infof("kind: %v", kind)


	

	//t2:= reflect.Indirect(reflect.ValueOf(src)).Type()

*/

	items := reflect.ValueOf(src)	
	

	cnt := items.Len()

	c.Infof("len: %v", items.Len())


	full_type :=  reflect.Indirect(items).Index(0).Type().String()

	arr := strings.Split(full_type, ".")

 	kind := arr[len(arr)-1]
	//kind := getStructName(t)
	c.Infof("kind: %v", kind)



	incomplete_keys := make([]*datastore.Key, cnt)

	for i := 0; i < cnt; i++ {
		incomplete_keys[i] = datastore.NewKey(c, kind, "", 0, nil)
	}

	keys, err := datastore.PutMulti(c, incomplete_keys, src)
	if err != nil {
		return nil, err
	}


	c.Infof("len(keys): %v", len(keys))

	return keys, nil


/*
	kind := getStructName(items)


	incomplete_keys := make([]*datastore.Key, len(items))

	for i := 0; i < len(*items); i++ {
		incomplete_keys[i] = datastore.NewKey(c, kind, "", 0, nil)
	}

	keys, err := datastore.PutMulti(c, incomplete_keys, *items)
	if err != nil {
		return nil, err
	}

	return keys, nil
	*/
	//return nil, nil
}



