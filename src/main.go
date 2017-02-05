package munhasir

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// route
	http.HandleFunc("/", indexHandler)
	http.Handle("/statics/",
		http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics"))),
	)
	http.ListenAndServe(":" + port, nil)
}