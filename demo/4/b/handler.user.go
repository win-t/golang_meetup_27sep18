package main

import (
	"net/http"

	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/kv"
	"github.com/payfazz/go-router/method"
	"github.com/payfazz/go-router/path"
	"github.com/payfazz/go-router/segment"
)

func compileUserHandler(s storage) http.HandlerFunc {
	return middleware.Compile(
		authMiddleware(s),
		path.C(path.H{
			"/:id": segment.E(method.C(method.H{
				"GET": func(w http.ResponseWriter, r *http.Request) {
					id, _ := segment.Get(r, "id")
					if id == "me" {
						uObj := kv.Get(r, "currentUser").(user)
						responseJson(w, struct {
							ID   int    `json:"id"`
							Name string `json:"name"`
						}{uObj.id, uObj.name})
					} else {
						http.Error(w, "Not Yet Implemented", http.StatusNotImplemented)
					}
				},
			})),
		}),
	)
}
