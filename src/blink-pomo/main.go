package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("blink-pomo: Pretty lights and productivity")

	http.Handle("/", http.FileServer(http.Dir("assets")))

	http.ListenAndServe(":9913", nil)
}
