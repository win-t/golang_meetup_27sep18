package main

import (
	"fmt"
	"net/http"

	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common"
	"github.com/payfazz/go-router/path"
)

func main() {
	storage, err := openStorage(
		getConf("POSTGRES_ENDPOINT"),
		getConf("POSTGRES_USER"),
		getConf("POSTGRES_PASSWORD"),
		getConf("POSTGRES_DB"),
	)
	errMustNil(err)
	defer storage.Close()

	handler := middleware.Compile(
		common.BasicPack(),
		path.C(path.H{
			"/users":  compileUserHandler(storage),
			"/coupon": compileCouponHandler(storage),
		}),
	)

	fmt.Println("Server is running")
	errMustNil(http.ListenAndServe(getConf("ADDR"), handler))
}

func errMustNil(err error) {
	if err != nil {
		panic(err)
	}
}
