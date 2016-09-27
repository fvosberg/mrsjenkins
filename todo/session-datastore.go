package todo

type SessionDatastore struct {
	data []*Todo
}

func (s *SessionDatastore) Create(todo *Todo) {
	s.data = append(s.data, todo)
}

func (s *SessionDatastore) Index() []*Todo {
	return s.data
}
