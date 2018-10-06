package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rumyantseva/go-sofia/internal/diagnostics"
)

type serverConf struct {
	port   string
	router http.Handler
	name   string
}

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

	diagnostics := diagnostics.NewDiagnostics()

	possibleErrors := make(chan error, 2)

	configurations := []serverConf{
		{
			port:   blPort,
			router: router,
			name:   "application server",
		},
		{

			port:   diagPort,
			router: diagnostics,
			name:   "diagnostics server",
		},
	}

	servers := make([]*http.Server, 2)

	for i, c := range configurations {
		go func(conf serverConf, i int) {
			log.Printf("The %s is preparing to handle connections...", conf.name)
			servers[i] = &http.Server{
				Addr:    ":" + conf.port,
				Handler: conf.router,
			}
			err := servers[i].ListenAndServe()
			if err != nil {
				possibleErrors <- err
			}
		}(c, i)
	}

	select {
	case err := <-possibleErrors:
		for _, s := range servers {
			timeout := 5 * time.Second
			log.Printf("\nShutdown with timeout: %s\n", timeout)
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			customError := s.Shutdown(ctx)
			if customError != nil {
				fmt.Println(customError)
			}
			log.Printf("Server gracefully stopped")
		}
		log.Fatal(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Print("The hello handler was called")
	fmt.Fprint(w, http.StatusText(http.StatusOK))
}
