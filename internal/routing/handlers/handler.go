package handlers

import (
	"encoding/json"
	ethereumRPC "ethereum-parser/internal/services/ethereum-rpc"
	"net/http"
)

type Handler interface {
	GetCurrentBlock(w http.ResponseWriter, r *http.Request)
	Subscribe(w http.ResponseWriter, r *http.Request)
	GetTransactions(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	ParserService ethereumRPC.Parser
}

func NewHandler(parser ethereumRPC.Parser) *handler {
	return &handler{
		ParserService: parser,
	}
}

func (h *handler) GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	lastParsedBlock := h.ParserService.GetCurrentBlock(r.Context())

	jsonResponse, marshalErr := json.Marshal(lastParsedBlock)
	if marshalErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func (h *handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ok := h.ParserService.Subscribe(r.Context(), address)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(map[string]string{"message": "address subscribed"})
	w.Write(response)
}

func (h *handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	transactions := h.ParserService.GetTransactions(r.Context(), address)
	resp, err := json.Marshal(transactions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
