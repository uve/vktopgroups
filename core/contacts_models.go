package core

import (
	"appengine"
	"appengine/datastore"
)


const (

	CONTACT_KIND  = "Contact"
)


type Contact struct {
	User_id	       int    `json:"user_id"`
	Desc 		   string `json:"desc"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`

	Group_id       *datastore.Key// `datastore:"Group_id,noindex"`
	Project_id     *datastore.Key// `datastore:"Project_id,noindex"`
	Custom_id      *datastore.Key// `datastore:"Custom_id,noindex"`
}



func putContacts(c appengine.Context, items *[]Contact) ([]*datastore.Key, error){

	incomplete_keys := make([]*datastore.Key, len(*items))

	for i := 0; i < len(*items); i++ {
		incomplete_keys[i] = datastore.NewKey(c, CONTACT_KIND, "", 0, nil)

	}

	keys, err := datastore.PutMulti(c, incomplete_keys, *items)

	if err != nil {
		return nil, err
	}

	return keys, nil
}




func getContacts(c appengine.Context) ([]*Contact, error) {

	
	q:= datastore.NewQuery(CONTACT_KIND)

	results := make([]*Contact, 0, QUERY_MAX)
	_, err := q.GetAll(c, &results)
	if err != nil {
		return nil, err
	}



	return results, nil
}





func deleteContacts(c appengine.Context) (error) {

	
	q:= datastore.NewQuery(CONTACT_KIND).KeysOnly()

	results := make([]*Contact, 0, QUERY_MAX)
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





