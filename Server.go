package demo

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

var visitors int

func handleHi(w http.ResponseWriter, req *http.Request) {
	match, _ := regexp.MatchString(`^\w*$`, req.FormValue("color"))

	if !match {
		http.Error(w, "Optional color is invalid", http.StatusBadRequest)
		return
	}

	//visitors++
	now := visitors + 1
	visitors = now

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<h1 style='color: " + req.FormValue("color") + "'>Welcome!</h1> You are visitor number " + fmt.Sprint(visitors) + "!"))
}

func main() {
	log.Printf("Starting on port 8080")
	http.HandleFunc("/hi", handleHi)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))

}
