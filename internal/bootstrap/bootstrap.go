package bootstrap

import (
	"fmt"

	"github.com/geomodular/go-htmx-todo/internal/model"
)

func FillSimple(m model.Model, n int) {
	for i := range n {
		note := fmt.Sprintf("My task number #%d", i)
		m.Add(note)
	}
}
