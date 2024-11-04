package server

import (
	"context"
	"html/template"
	"net/http"
	"strconv"

	"github.com/geomodular/go-htmx-todo/internal/bootstrap"
	"github.com/geomodular/go-htmx-todo/internal/model"
	"github.com/geomodular/go-htmx-todo/internal/pagination"
	"github.com/gorilla/mux"
)

type homeHandler struct {
	template *template.Template
	db       model.Model
}

type homeData struct {
	Tasks []model.Task
	Pages []pagination.Page
}

func (h homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	items := h.db.List(0, 5, true)
	pages := pagination.ComputePages(5, 0, 5, h.db.Length())

	w.Header().Set("Content-Type", "text/html")
	err := h.template.Execute(w, homeData{items, pages})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type tasksHandler struct {
	template *template.Template
	db       model.Model
}

func (h tasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/html")
	if r.Method == "DELETE" {
		h.db.Remove(id)

		items := h.db.List(0, 5, true)
		err = h.template.ExecuteTemplate(w, "tasks", items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "GET" {
		item, ok := h.db.Get(id)
		if !ok {
			http.Error(w, "", http.StatusNotFound)
		}

		err = h.template.ExecuteTemplate(w, "task", item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		r.ParseForm()
		value := r.FormValue("action")
		if value == "check" {
			h.db.Complete(id)
		} else if value == "uncheck" {
			h.db.Uncomplete(id)
		}

		item, ok := h.db.Get(id)
		if !ok {
			http.Error(w, "", http.StatusNotFound)
		}

		err = h.template.ExecuteTemplate(w, "task", item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type taskCreateHandler struct {
	template *template.Template
	db       model.Model
}

func (h taskCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		value := r.FormValue("note")
		if value != "" {
			h.db.Add(value)
		}

		items := h.db.List(0, 5, true)
		err := h.template.ExecuteTemplate(w, "tasks", items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte("<input id=\"task-add-input\" name=\"note\" type=\"text\" class=\"w3-input w3-border\" hx-swap-oob=\"true\">"))
	}
}

func Run(ctx context.Context, host string) error {
	db := model.NewMemModel()
	bootstrap.FillSimple(db, 10)

	homeTmpl := template.Must(template.ParseFiles(
		"web/template/root.gohtml",
		"web/template/tasks.gohtml",
		"web/template/home.gohtml"))

	taskTmpl := template.Must(template.ParseFiles(
		"web/template/tasks.gohtml",
	))

	r := mux.NewRouter()
	r.Handle("/", homeHandler{homeTmpl, db})
	r.Handle("/tasks/{id:[0-9]+}", tasksHandler{taskTmpl, db}).Methods("DELETE", "GET", "POST")
	r.Handle("/tasks", taskCreateHandler{taskTmpl, db}).Methods("POST")

	// TODO: Catch ctrl+c signal
	// https://github.com/gorilla/mux
	return http.ListenAndServe(host, r)
}
