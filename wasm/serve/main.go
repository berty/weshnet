package main

import "net/http"

func main() {
	http.ListenAndServe(`:4242`, http.FileServer(http.Dir(`.`)))
}
