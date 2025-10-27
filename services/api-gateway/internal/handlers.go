package internal

import (
	"context"
	"encoding/json"
	"net/http"
)

type Handlers struct {
	client *Client
}

func NewHandlers(client *Client) *Handlers {
	return &Handlers{client: client}
}

type putRequest struct {
	Value string `json:"value"`
}

func (h *Handlers) Put(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	var req putRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.client.Put(context.Background(), key, []byte(req.Value)); err != nil {
		http.Error(w, "failed to put", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	value, err := h.client.Get(context.Background(), key)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	resp := putRequest{Value: string(value)}
	json.NewEncoder(w).Encode(resp)
}

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	if err := h.client.Delete(context.Background(), key); err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
