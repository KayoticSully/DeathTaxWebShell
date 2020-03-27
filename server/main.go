package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	router := fasthttprouter.New()
	router.GET("/", Index)

	log.Fatal(fasthttp.ListenAndServe(":5000", router.Handler))
}

func Index(ctx *fasthttp.RequestCtx) {
	cmd := exec.Command("echo", "DeathTax")
	output := &bytes.Buffer{}
	cmd.Stdout = output
	cmd.Stderr = output

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	cmd.Wait()
	fmt.Fprint(ctx, fmt.Sprintf("%s", output.String()))
}
