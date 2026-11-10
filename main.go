package main

import (
	"go_homeworks/internal"
	"log"
	"net/http"
)

func main() {
	internal.RegisterRouters()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
