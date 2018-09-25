package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/kv"
)

func authMiddleware(s storage) middleware.Func {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !func() bool {
				authParts := strings.Split(r.Header.Get("Authorization"), " ")
				if len(authParts) != 2 {
					return false
				}
				if strings.ToLower(authParts[0]) != "basic" {
					return false
				}
				authBytes, err := base64.StdEncoding.DecodeString(authParts[1])
				if err != nil {
					return false
				}
				authStr := string(authBytes)
				sep := strings.Index(authStr, ":")
				if sep == -1 {
					return false
				}
				user := authStr[:sep]
				pass := authStr[sep+1:]
				uObj, err := s.getUser(user, pass)
				if err != nil {
					fmt.Println("DEBUG", "FAIL s.getUser", err)
					return false
				}
				kv.Set(r, "currentUser", uObj)
				return true
			}() {
				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
			} else {
				next(w, r)
			}
		}
	}
}
