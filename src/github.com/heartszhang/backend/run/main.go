package main

import (
	"github.com/heartszhang/backend"
	"log"
	"net/http"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(16)
	err := http.ListenAndServe(backend.BackendConfig().Address(), nil)
	if err != nil {
		log.Fatal("backend: ", err)
	}
}
