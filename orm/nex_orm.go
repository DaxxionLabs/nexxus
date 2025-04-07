package nex_orm

import (
	"fmt"
	"net/http"
	"nexxus"
	"reflect"
	"time"
)

type test struct {
	Model
}

func (t test) hello() {
	t.Create(t)
}

type Model struct {
	id          int64
	createdDate time.Time
	updatedDate time.Time

	handler nexxus.HandlerFunc
}

func (m *Model) Create(i interface{}) {
	name := reflect.TypeOf(i).Name()
	route := fmt.Sprintf("/%s", name)
	m.handler(route, func(w http.ResponseWriter, r *http.Request) {
	})
}
