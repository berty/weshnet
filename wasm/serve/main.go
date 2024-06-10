package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("serving on :4242")
	http.ListenAndServe(`:4242`, http.FileServer(http.Dir(`.`)))
}
