package core

import (
	//"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

const (
	TIME_LAYOUT = "Jan 2, 2006 15:04:05 AM"
	PROJECT_KIND  = "Project"
)

// Project is an entity to store projects that have been inserted by users.
type Project struct {
	key *datastore.Key

/*
	Outcome string    `datastore:"outcome"`
	Played  time.Time `datastore:"played"`
	Player  string    `datastore:"player"`
*/
	User    string
	Name    string
}





// Turns the Project struct/entity into a ProjectRespMsg which is then used
// as an API response.

func (s *Project) toMessage(msg *ProjectRespMsg) *ProjectRespMsg {
	if msg == nil {
		msg = &ProjectRespMsg{}
	}
	msg.Id = s.key.IntID()
	msg.Name = s.Name
	//msg.Outcome = s.Outcome
	//msg.Played = s.timestamp()
	return msg
}


// timestamp formats date/time of the project.
func (s *Project) timestamp() string {
	//return s.Played.Format(TIME_LAYOUT)
	return ""
}

// put stores the project in the Datastore.
func (s *Project) put(c appengine.Context) (err error) {
	key := s.key
	if key == nil {
		key = datastore.NewIncompleteKey(c, PROJECT_KIND, nil)
	}
	key, err = datastore.Put(c, key, s)
	if err == nil {
		s.key = key
	}
	return
}

// newProject returns a new Project ready to be stored in the Datastore.
func newProject(name string, u *user.User) *Project {
	//return &Project{Outcome: outcome, Played: time.Now(), Player: userId(u)}
	return &Project{
				Name: name,
				User: userId(u) }
}



// newUserProjectQuery returns a Query which can be used to list all previous
// games of a user.
func newUserProjectQuery(u *user.User) *datastore.Query {
	return datastore.NewQuery(PROJECT_KIND).Filter("User =", userId(u))
}

// fetchProjects runs Query q and returns Project entities fetched from the
// Datastore.
func fetchProjects(c appengine.Context, q *datastore.Query, limit int) (
	[]*Project, error) {

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
