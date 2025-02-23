package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", servePage)

	fmt.Println("Starting front end service on port 80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Panic(err)
	}
}

func servePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/test.page.gohtml")
}
