package model

type Model interface {
	List(offset int, size int, reverse bool) []Task
	// ListReverse?
	Get(id int) (Task, bool)
	Add(note string)
	Remove(id int)
	Complete(id int)
	Uncomplete(id int)
	Edit(id int, note string)
	Length() int
}

type Task struct {
	ID      int
	Done    bool
	Deleted bool
	Note    string
}
