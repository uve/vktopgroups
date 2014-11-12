package controller

import (
	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"

	"net/url"
	"fmt"
	"encoding/json"
	"strings"
)



const (
	CONTACT_KIND   = "Contact"
	CONTACT_CURSOR = CONTACT_KIND + "_cursor"
	CONTACTS_LIMIT = 100

)


type ContactBase struct{

	VK_id	       int64  `datastore:"-" json:"id"`

	First_name 	   string `json:"first_name"`
	Last_name	   string `json:"last_name"`
	Photo_50	   string `json:"photo_50"`
	Photo_100	   string `json:"photo_100"`
}



type Contact struct {

	Default
	ContactBase

	User_id	       int64  `json:"user_id"`
	Desc 		   string `json:"desc"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`

	Group_id       *datastore.Key// `datastore:"Group_id,noindex"`
	Project_id     *datastore.Key// `datastore:"Project_id,noindex"`
	Custom_id      *datastore.Key// `datastore:"Custom_id,noindex"`

	Is_exist       bool
}





type ContactMessage struct {

	Response []ContactBase `json:"response"`
}



func fetchContacts(c appengine.Context, cursor_start string) (string, error) {

	limit := CONTACTS_LIMIT

	q:= datastore.NewQuery(CONTACT_KIND).Filter("Is_exist=", false).Limit(limit)

	cursor, err := datastore.DecodeCursor(cursor_start)
	if err == nil {
		q = q.Start(cursor)
	}

	var ids []int64

	var all_keys []*datastore.Key
	var all_items []Contact

	// Iterate over the results.
	t := q.Run(c)
	for {
		var p Contact
		key, err := t.Next(&p)
		if err == datastore.Done {
			break
		}
		if err != nil {
			c.Errorf("fetching next Person: %v", cursor)
			break
		}
		// Do something with the Person p

		ids = append(ids, p.User_id)

		all_keys  = append(all_keys, key)
		all_items = append(all_items, p)
	}

	contacts, err := saveContacts(c, ids)
	if err != nil {
		return "", err
	}

	for i,contact := range *contacts{
		all_items[i].ContactBase = contact
		all_items[i].Is_exist = true
	}


	_, err = datastore.PutMulti(c, all_keys, all_items)
	if err != nil {
		return "", err
	}


	return GetCursor(t, cursor)
}




func saveContacts(c appengine.Context, list []int64) (*[]ContactBase, error) {

	method := "execute.users_get"

	var q []string

	for _, item := range list{
		q = append(q, fmt.Sprintf("%d", item))
	}

	query := strings.Join(q, ",")

	c.Infof("query: %v", query)

	v := url.Values{}
	v.Set("access_token", VK.Token)
	v.Add("lang", "ru")
	v.Add("https", "1")

	v.Add("v", VK.ApiVersion)
	v.Add("method", method)

	v.Add("q", query)

	api_url := fmt.Sprintf("%s/%s", VK.Server, method)

	client := urlfetch.Client(c)
	resp, err := client.PostForm(api_url, v)

	if err != nil {
		return nil, err
	}

	var m ContactMessage

	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		panic(fmt.Sprintf("Could not parse : %s", err))
	}

	return &m.Response, nil
}



func queryContactsByCustom(c appengine.Context, custom_id *datastore.Key) (*datastore.Query){

	return datastore.NewQuery(CONTACT_KIND).Filter("Custom_id=", custom_id).Order("Created")
}
