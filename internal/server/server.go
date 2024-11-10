package server

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/geomodular/go-htmx-todo/internal/bootstrap"
	"github.com/geomodular/go-htmx-todo/internal/model"
	"github.com/geomodular/go-htmx-todo/internal/pagination"
	"github.com/gorilla/mux"
)

const (
	maxPaginationSize = 5
	maxTasksPerPage   = 7
)

func reportInternalError(w http.ResponseWriter, err error, msg string) {
	code := http.StatusInternalServerError
	slog.Error(msg, "err", err)
	http.Error(w, http.StatusText(code), code)
}

func parseFormInteger(r *http.Request, key string, default_ int) int {
	ret := default_
	if arg := r.FormValue(key); arg != "" {
		var err error
		ret, err = strconv.Atoi(arg)
		if err != nil {
			msg := fmt.Sprintf("failed converting %s to integer", key)
			slog.Error(msg, key, arg)
			ret = default_
		}
	}
	return ret
}

type Task struct {
	Task model.Task
	Page pagination.Page // Info about the actual page
}

type homeHandler struct {
	template *template.Template
	db       model.Model
}

type todoData struct {
	Tasks []Task
	Pages []pagination.Page // Info about all the pages around the actual page
}

func (h homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	offset := parseFormInteger(r, "offset", 0)
	size := parseFormInteger(r, "size", maxTasksPerPage)

	pages := pagination.ComputePages(maxPaginationSize, offset, maxTasksPerPage, h.db.Length())
	pages = pagination.IncludeArrowPages(pages, maxPaginationSize, offset, maxTasksPerPage, h.db.Length())
	page := pagination.GetActualPage(maxPaginationSize, offset, maxTasksPerPage, h.db.Length())
	items := h.db.List(offset, size, true)

	var tasks []Task
	for _, item := range items {
		tasks = append(tasks, Task{item, page})
	}

	w.Header().Set("Content-Type", "text/html")
	err := h.template.Execute(w, todoData{tasks, pages})
	if err != nil {
		reportInternalError(w, err, "failed executing template")
	}
}

type taskHandler struct {
	template *template.Template
	db       model.Model
}

func (h taskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		reportInternalError(w, err, "failed converting to integer")
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if r.Method == "DELETE" {
		r.ParseForm()

		h.db.Remove(id)

		offset := parseFormInteger(r, "offset", 0)
		size := parseFormInteger(r, "size", maxTasksPerPage)

		pages := pagination.ComputePages(maxPaginationSize, offset, maxTasksPerPage, h.db.Length())
		pages = pagination.IncludeArrowPages(pages, maxPaginationSize, offset, maxTasksPerPage, h.db.Length())
		page := pagination.GetActualPage(maxPaginationSize, offset, maxTasksPerPage, h.db.Length())
		items := h.db.List(offset, size, true)

		var tasks []Task
		for _, item := range items {
			tasks = append(tasks, Task{item, page})
		}

		err = h.template.ExecuteTemplate(w, "todo", todoData{tasks, pages})
		if err != nil {
			reportInternalError(w, err, "failed executing template")
		}
	} else if r.Method == "POST" {
		r.ParseForm()

		value := r.FormValue("checkbox")
		if value == "on" {
			h.db.Complete(id)
		} else {
			h.db.Uncomplete(id)
		}
		offset := parseFormInteger(r, "offset", 0)
		page := pagination.GetActualPage(maxPaginationSize, offset, maxTasksPerPage, h.db.Length())

		item, ok := h.db.Get(id)
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		err = h.template.ExecuteTemplate(w, "task", Task{item, page})
		if err != nil {
			reportInternalError(w, err, "failed executing template")
		}
	}
}

type tasksHandler struct {
	template *template.Template
	db       model.Model
}

func (h tasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	value := r.FormValue("note")
	if value != "" {
		h.db.Add(value)
	}
	pages := pagination.ComputePages(maxPaginationSize, 0, maxTasksPerPage, h.db.Length())
	pages = pagination.IncludeArrowPages(pages, maxPaginationSize, 0, maxTasksPerPage, h.db.Length())
	page := pagination.GetActualPage(maxPaginationSize, 0, maxTasksPerPage, h.db.Length())
	items := h.db.List(0, maxTasksPerPage, true)

	var tasks []Task
	for _, item := range items {
		tasks = append(tasks, Task{item, page})
	}

	err := h.template.ExecuteTemplate(w, "todo", todoData{tasks, pages})
	if err != nil {
		reportInternalError(w, err, "failed executing template")
	}
}

func Run(ctx context.Context, host string) error {
	db := model.NewMemModel()
	bootstrap.FillSimple(db, 50)

	homeTmpl := template.Must(template.ParseFiles(
		"web/template/base.gohtml",
		"web/template/todo.gohtml",
		"web/template/home.gohtml"))

	todoTmpl := template.Must(template.ParseFiles(
		"web/template/todo.gohtml",
	))

	r := mux.NewRouter()
	r.Handle("/", homeHandler{homeTmpl, db}).Methods("GET")
	r.Handle("/tasks/{id:[0-9]+}", taskHandler{todoTmpl, db}).Methods("DELETE", "POST")
	r.Handle("/tasks", tasksHandler{todoTmpl, db}).Methods("POST")

	// TODO: Catch ctrl+c signal
	// https://github.com/gorilla/mux
	return http.ListenAndServe(host, r)
}
