package core

import (
	"appengine"
	"appengine/datastore"
)




type Contacts struct {
	User_id	       int    `json:"user_id"`
	Desc 		   string `json:"desc"`
	Phone          string `json:"phone"`
}


type Group struct {
	Id		       int    `json:"id"`
	Name 		   string `json:"name"`
	Screen_name	   string `json:"screen_name"`

	Is_closed	   int    `json:"is_closed"`
	Type		   string `json:"type"`

	Members_count  int    `json:"members_count"`

	Contacts 	   []Contacts `json:"contacts"`


	Site		   string `json:"site"`
	Photo_50	   string `json:"photo_50"`
	Photo_100	   string `json:"photo_100"`
	Photo_200	   string `json:"photo_200"`

}


func getCustom(c appengine.Context, id int64) (*Custom, error){

	k := datastore.NewKey(c, "Custom", "", id, nil)
	e := new(Custom)
	if err := datastore.Get(c, k, e); err != nil {
		//http.Error(w, err.Error(), 500)
		return nil, err
	}


	return e, nil
}


func putMulti(c appengine.Context, s *[]Group) (int, error){

	// A batch put.
	//_, err := datastore.PutMulti(c, []*datastore.Key{k1, k2, k3}, []interface{}{e1, e2, e3})


	keys := make([]*datastore.Key, len(*s))
	//rs := make([]*Record, len(emails))

	for i := 0; i < len(*s); i++ {
		keys[i] = datastore.NewKey(c, "Group", "", 0, nil)
	}

	ks, err := datastore.PutMulti(c, keys, *s)
	if err != nil {
		return 0, nil
	}

	return len(ks), nil


}


