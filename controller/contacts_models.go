package controller

import "appengine/datastore"




type Contact struct {
	User_id	       int    `json:"user_id"`
	Desc 		   string `json:"desc"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`

	Group_id       *datastore.Key// `datastore:"Group_id,noindex"`
	Project_id     *datastore.Key// `datastore:"Project_id,noindex"`
	Custom_id      *datastore.Key// `datastore:"Custom_id,noindex"`
}




