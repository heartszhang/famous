package main

import (
	"github.com/heartszhang/backend"
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe(backend.BackendConfig().Address(), nil)
	if err != nil {
		log.Fatal("backend: ", err)
	}
}
