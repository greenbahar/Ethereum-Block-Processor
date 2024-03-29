package routing

import (
	"ethereum-parser/internal/routing/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func SetUpRouters(handler handlers.Handler) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/get-current-block", handler.GetCurrentBlock).Methods("GET")
	router.HandleFunc("/subscribe", handler.Subscribe).Methods("POST")
	router.HandleFunc("/get-transactions", handler.GetTransactions).Methods("GET")

	return router
}
