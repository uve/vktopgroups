package controller

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"

)

const (
	PROJECT_KIND  = "Project"
)




// Project is an entity to store projects that have been inserted by users.
type Project struct {

	*Default



	key *datastore.Key `datastore:"-"`

	User    string
	Name    string

	/*Created time.Time*/
}




func projectCreate(name string, user string) (*Project){

	return &Project{
		Name: name,
		User: user,
		Default: NewDefault(),
	}
}


func (s *Project) put(c appengine.Context) (err error) {

	return s.Default.put(s, c)
}

// Turns the Project struct/entity into a ProjectRespMsg which is then used
// as an API response.

func (s *Project) toMessage(msg *ProjectRespMsg) *ProjectRespMsg {
	if msg == nil {
		msg = &ProjectRespMsg{}
	}
	//msg.Id = s.key.IntID()
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
func fetchProjects(c appengine.Context, q *datastore.Query, limit int) ([]*Project, error) {

	projects := make([]*Project, 0, limit)
	keys, err := q.Limit(limit).GetAll(c, &projects)
	if err != nil {
		return nil, err
	}
	for i, project := range projects {
		project.key = keys[i]
	}
	return projects, nil
}

// userId returns a string ID of the user u to be used as Player of Project.
func userId(u *user.User) string {
	return u.String()
}

