package main

import (
	"fmt"
	"net/http"

	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/kv"
	"github.com/payfazz/go-router/method"
	"github.com/payfazz/go-router/path"
	"github.com/payfazz/go-router/segment"
)

func compileCouponHandler(s storage) http.HandlerFunc {
	return middleware.Compile(
		authMiddleware(s),
		path.C(path.H{
			"/:id/check": segment.E(method.C(method.H{
				"GET": func(w http.ResponseWriter, r *http.Request) {
					couponID, _ := segment.Get(r, "id")
					uObj := kv.Get(r, "currentUser").(user)
					res := executeCoupon(s, uObj.id, couponID, false)
					fmt.Fprintf(w, "%d", res)
				},
			})),
			"/:id/use": segment.E(method.C(method.H{
				"POST": func(w http.ResponseWriter, r *http.Request) {
					couponID, _ := segment.Get(r, "id")
					uObj := kv.Get(r, "currentUser").(user)
					res := executeCoupon(s, uObj.id, couponID, true)
					fmt.Fprintf(w, "%d", res)
				},
			})),
		}),
	)
}
