package diagnostics

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func NewDiagnostics() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/healthz", helthz)
	router.HandleFunc("/readyz", ready)
	return router
}

func helthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, http.StatusText(http.StatusOK))
}

func ready(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, http.StatusText(http.StatusOK))
}
