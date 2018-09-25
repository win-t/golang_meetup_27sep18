package main

import (
	"os"
)

var defConf = map[string]string{
	"ADDR": ":8080",

	"POSTGRES_ENDPOINT": "localhost:5432",
	"POSTGRES_USER":     "testuser",
	"POSTGRES_PASSWORD": "testpassword",
	"POSTGRES_DB":       "testdb",
}

func getConf(key string) string {
	ret := os.Getenv(key)
	if ret == "" {
		return defConf[key]
	}
	return ret
}
