package controller

import (

"appengine"
"appengine/datastore"
"appengine/user"

"model"
)


const (
	PROJECT_KIND  = "Project"
)


type Project struct {

	Default

	User    string
	Name    string
}




func NewProject(name string, user string) (*Project){

	return &Project{
		Name: name,
		User: user,
		Default: NewDefault(),
	}
}


// Turns the Project struct/entity into a ProjectRespMsg which is then used
// as an API response.

func (s *Project) toMessage(msg *ProjectRespMsg) *ProjectRespMsg {

	if msg == nil {
		msg = &ProjectRespMsg{}
	}
	
	msg.Name = s.Name

	msg.Id 		= s.Id()
	msg.Created = s.GetCreated()

	return msg
}


func (src *Project) Get(c appengine.Context, id int64) (*datastore.Key, error) {

	key := datastore.NewKey(c, model.GetKind(src), "", id, nil)

	if err := datastore.Get(c, key, src); err != nil {
		return nil, err
	}

	src.SetKey(key)

	return key, nil
}



func (src *Project) Put(c appengine.Context) (*datastore.Key, error) {

	key, err := model.Put(c, src)
	if err != nil {
		return nil, err
	}

	src.SetKey(key)

	return key, nil
}






// newUserProjectQuery returns a Query which can be used to list all previous
// games of a user.
func queryProjectByUser(u *user.User) *datastore.Query {
	return datastore.NewQuery(PROJECT_KIND).Filter("User =", userId(u)).Order("Created")
}


// Datastore.
func fetchProjects(c appengine.Context, q *datastore.Query, limit int) ([]*Project, error) {

	items := make([]*Project, 0, limit)
	keys, err := q.Limit(limit).GetAll(c, &items)
	if err != nil {
		return nil, err
	}
	for i, item := range items {
		item.SetKey(keys[i])
	}
	return items, nil
}

// userId returns a string ID of the user u to be used as Player of Project.
func userId(u *user.User) string {
	return u.String()
}

