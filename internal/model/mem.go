package model

type MemModel struct {
	nextID int
	data   map[int]*Task
	keys   []int // List of not deleted, sorted ids.
}

func NewMemModel() Model {
	return &MemModel{0, map[int]*Task{}, []int{}}
}

func (mm *MemModel) List(offset int, size int, reverse bool) []Task {
	tasks := []Task{}
	i := offset
	if i >= 0 && i <= len(mm.keys) {
		j := offset + size
		if j > len(mm.keys) {
			j = len(mm.keys)
		}
		if j < i {
			j = i
		}
		var keys []int
		if !reverse {
			keys = mm.keys[i:j]
		} else {
			i = (len(mm.keys) - i)
			j = (len(mm.keys) - j)
			for _, k := range mm.keys[j:i] {
				keys = append([]int{k}, keys...)
			}
		}
		for _, k := range keys {
			tasks = append(tasks, *mm.data[k])
		}
	}
	return tasks
}

func (mm *MemModel) Get(id int) (Task, bool) {
	if task, ok := mm.data[id]; ok {
		return *task, true
	}
	return Task{}, false
}

func (mm *MemModel) Add(note string) {
	mm.data[mm.nextID] = &Task{
		ID:   mm.nextID,
		Done: false,
		Note: note,
	}
	mm.keys = append(mm.keys, mm.nextID)
	mm.nextID++
}

func (mm *MemModel) Remove(id int) {
	if task, ok := mm.data[id]; ok {
		task.Deleted = true
	}
	for i, key := range mm.keys {
		if key == id {
			mm.keys = append(mm.keys[:i], mm.keys[i+1:]...)
			return
		}
	}
}

func (mm *MemModel) Complete(id int) {
	if task, ok := mm.data[id]; ok {
		task.Done = true
	}
}

func (mm *MemModel) Uncomplete(id int) {
	if task, ok := mm.data[id]; ok {
		task.Done = false
	}
}

func (mm *MemModel) Edit(id int, note string) {
	if task, ok := mm.data[id]; ok {
		task.Note = note
	}
}

func (mm *MemModel) Length() int {
	return len(mm.keys)
}
