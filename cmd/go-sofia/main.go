package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rumyantseva/go-sofia/internal/diagnostics"
)

func main() {
	log.Print("Starting the application...")

	blPort := os.Getenv("PORT")
	if len(blPort) == 0 {
		log.Fatal("The application port should be set")
	}

	diagPort := os.Getenv("DIAG_PORT")
	if len(diagPort) == 0 {
		log.Fatal("The diagnostics port should be set")
	}

	router := mux.NewRouter()
	router.HandleFunc("/", hello)
	go func() {
		log.Print("The application server is preparing to handle connections...")
		err := http.ListenAndServe(":"+blPort, router)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Print("The diagnostics server is preparing to handle connections...")
	diagnostics := diagnostics.NewDiagnostics()
	err := http.ListenAndServe(":"+diagPort, diagnostics)
	if err != nil {
		log.Fatal(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Print("The hello handler was called")
	fmt.Fprint(w, http.StatusText(http.StatusOK))
}
