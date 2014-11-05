package core

import (
	//"appengine"
	"appengine/datastore"
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


/*

func getContacts(c appengine.Context) ([]*Contact, error) {

	
	q:= datastore.NewQuery(CONTACT_KIND)

	results := make([]*Contact, 0, 5000)
	_, err := q.GetAll(c, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

*/





