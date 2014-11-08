package controller

import "appengine/datastore"

const (
	PROJECT_KIND  = "Project"
)




// Project is an entity to store projects that have been inserted by users.
type Project struct {

	Default



	key *datastore.Key `datastore:"-"`

	User    string
	Name    string

	/*Created time.Time*/
}




func NewProject(name string, user string) (*DefaultInterface){

	return &DefaultInterface{
		Name: name,
		User: user,
		Default: NewDefault(),
	}
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
*/

// Turns the Project struct/entity into a ProjectRespMsg which is then used
// as an API response.

func (s *Project) toMessage(msg *ProjectRespMsg) *ProjectRespMsg {

	if msg == nil {
		msg = &ProjectRespMsg{}
	}
	msg.Id = s.Id()
	msg.Name = s.Name
	//msg.Created = s.timestamp()

	return msg
}





// newUserProjectQuery returns a Query which can be used to list all previous
// games of a user.
func newUserProjectQuery(u *user.User) *datastore.Query {
	return datastore.NewQuery(PROJECT_KIND).Filter("User =", userId(u)).Order("-Created")
}

// fetchProjects runs Query q and returns Project entities fetched from the
// Datastore.
func fetch(c appengine.Context, q *datastore.Query, limit int) ([]*Project, error) {

	projects := make([]*Project, 0, limit)
	keys, err := q.Limit(limit).GetAll(c, &projects)
	if err != nil {
		return nil, err
	}
	for i, project := range projects {
		project.setKey(keys[i])
	}
	return projects, nil
}

// userId returns a string ID of the user u to be used as Player of Project.
func userId(u *user.User) string {
	return u.String()
}

