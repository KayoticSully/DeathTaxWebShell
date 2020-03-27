package main

import (
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	router := fasthttprouter.New()
	router.GET("/", Index)

	log.Fatal(fasthttp.ListenAndServe(":5000", router.Handler))
}

func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Welcome!\n")
}