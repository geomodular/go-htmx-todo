package model

import (
	"fmt"
	"testing"
)

func makeMemModel() *MemModel {
	return &MemModel{0, map[int]*Task{}, []int{}}
}

func addTasks(mm *MemModel, n int) {
	for i := range n {
		mm.Add(fmt.Sprintf("my task %d", i+1))
	}
}

func mapTasksToIDs(tasks []Task) []int {
	ids := []int{}
	for _, task := range tasks {
		ids = append(ids, task.ID)
	}
	return ids
}

func areListsSame(l1 []int, l2 []int) bool {
	if len(l1) != len(l2) {
		return false
	}
	for i := range l1 {
		if l1[i] != l2[i] {
			return false
		}
	}
	return true
}

func TestMemAdd(t *testing.T) {
	m := makeMemModel()
	addTasks(m, 3)

	if len(m.data) != 3 {
		t.Errorf("wrong number of entries, want 3, got %d", len(m.data))
	}

	if len(m.keys) != 3 {
		t.Errorf("wrong number of keys, want 3, got %d", len(m.data))
	}

	if m.keys[0] != 0 {
		t.Errorf("wrong id, want 0, got %d", m.keys[0])
	}

	if m.keys[1] != 1 {
		t.Errorf("wrong id, want 1, got %d", m.keys[1])
	}

	if m.keys[2] != 2 {
		t.Errorf("wrong id, want 2, got %d", m.keys[2])
	}
}

func TestMemGet(t *testing.T) {
	m := makeMemModel()
	addTasks(m, 5)

	task, ok := m.Get(3)
	if !ok {
		t.Errorf("incorrect flag, item should exist")
	}

	if task.ID != 3 {
		t.Errorf("wrong id, want 3, got %d", task.ID)
	}

	task, ok = m.Get(15)
	if ok {
		t.Errorf("incorrect flag, item should not exist")
	}

}

func TestMemRemove(t *testing.T) {
	m := makeMemModel()
	addTasks(m, 3)

	m.Remove(1)

	if len(m.data) != 3 {
		t.Errorf("wrong number of entries, want 3, got %d", len(m.data))
	}

	if m.data[1].Deleted != true {
		t.Error("wrong flag, want deleted to be true")
	}

	if len(m.keys) != 2 {
		t.Errorf("wrong number of keys, want 2, got %d", len(m.keys))
	}

	if m.keys[0] != 0 || m.keys[1] != 2 {
		t.Errorf("incorrect keys, want [0 2], got %v", m.keys)
	}
}

func TestMemComplete(t *testing.T) {
	m := makeMemModel()
	addTasks(m, 2)

	if m.data[1].Done != false {
		t.Fatal("wrong flag, want false")
	}

	m.Complete(1)

	if m.data[1].Done != true {
		t.Error("wrong flag, want true")
	}
}

func TestMemUncomplete(t *testing.T) {
	m := makeMemModel()
	addTasks(m, 2)

	m.data[1].Done = true

	m.Uncomplete(1)

	if m.data[1].Done != false {
		t.Error("wrong flag, want false")
	}
}

func TestMemEdit(t *testing.T) {
	m := makeMemModel()
	addTasks(m, 2)

	if m.data[1].Note != "my task 2" {
		t.Fatalf("wrong note, want `my task 2`, got %s", m.data[1].Note)
	}

	m.Edit(1, "my task 3")

	if m.data[1].Note != "my task 3" {
		t.Errorf("wrong note, want `my task 3`, got %s", m.data[1].Note)
	}
}

func TestMemList(t *testing.T) {
	m := makeMemModel()
	addTasks(m, 10)

	tasks := m.List(0, 10, false)
	if len(tasks) != 10 {
		t.Errorf("wrong number of tasks, want 10, got %d", len(tasks))
	}

	tasks = m.List(0, 100, false)
	if len(tasks) != 10 {
		t.Errorf("wrong number of tasks, want 10, got %d", len(tasks))
	}

	tasks = m.List(5, 4, false)
	if len(tasks) != 4 {
		t.Errorf("wrong number of tasks, want 4, got %d", len(tasks))
	}
	if !areListsSame(mapTasksToIDs(tasks), []int{5, 6, 7, 8}) {
		t.Errorf("incorrect keys, want [5 6 7 8], got %v", mapTasksToIDs(tasks))
	}

	tasks = m.List(5, 10, false)
	if len(tasks) != 5 {
		t.Errorf("wrong number of tasks, want 5, got %d", len(tasks))
	}
	if !areListsSame(mapTasksToIDs(tasks), []int{5, 6, 7, 8, 9}) {
		t.Errorf("incorrect keys, want [5 6 7 8 9], got %v", mapTasksToIDs(tasks))
	}

	tasks = m.List(11, 10, false)
	if len(tasks) != 0 {
		t.Errorf("wrong number of tasks, want 0, got %d", len(tasks))
	}

	tasks = m.List(5, -1, false)
	if len(tasks) != 0 {
		t.Errorf("wrong number of tasks, want 0, got %d", len(tasks))
	}

	tasks = m.List(5, 0, false)
	if len(tasks) != 0 {
		t.Errorf("wrong number of tasks, want 0, got %d", len(tasks))
	}
}

func TestMemListReversed(t *testing.T) {
	m := makeMemModel()
	addTasks(m, 10)

	tasks := m.List(0, 5, true)
	if !areListsSame(mapTasksToIDs(tasks), []int{9, 8, 7, 6, 5}) {
		t.Errorf("incorrect keys, want [9 8 7 6 5], got %v", mapTasksToIDs(tasks))
	}

	tasks = m.List(1, 3, true)
	if !areListsSame(mapTasksToIDs(tasks), []int{8, 7, 6}) {
		t.Errorf("incorrect keys, want [8 7 6], got %v", mapTasksToIDs(tasks))
	}
}

func TestMemRemoveAndList(t *testing.T) {
	m := makeMemModel()
	addTasks(m, 10)

	m.Remove(3)
	m.Remove(7)

	tasks := m.List(0, 100, false)
	if len(tasks) != 8 {
		t.Errorf("wrong number of tasks, want 8, got %d", len(tasks))
	}
	if !areListsSame(mapTasksToIDs(tasks), []int{0, 1, 2, 4, 5, 6, 8, 9}) {
		t.Errorf("incorrect keys, want [0 1 2 4 5 6 8 9], got %v", mapTasksToIDs(tasks))
	}

	m.Remove(5)

	tasks = m.List(0, 100, false)
	if len(tasks) != 7 {
		t.Errorf("wrong number of tasks, want 7, got %d", len(tasks))
	}
	if !areListsSame(mapTasksToIDs(tasks), []int{0, 1, 2, 4, 6, 8, 9}) {
		t.Errorf("incorrect keys, want [0 1 2 4 6 8 9], got %v", mapTasksToIDs(tasks))
	}
}

func TestMemLength(t *testing.T) {
	m := makeMemModel()
	addTasks(m, 10)

	m.Remove(3)

	if m.Length() != 9 {
		t.Errorf("incorrect number of itemst, want 9, got %d", m.Length())
	}
}
