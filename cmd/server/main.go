package main

import (
	"fmt"
	"net/http"
)

// default route for "/"
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ai Task Processor running on port 6969")
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Health is OK")
}

func main() {
	println("AI-TASK-PROCESSOR")

	//Routes
	http.HandleFunc("/", handler)
	http.HandleFunc("/health", health)

	//Server
	http.ListenAndServe(":6969", nil)
}
