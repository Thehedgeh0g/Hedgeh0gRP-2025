package transport

import (
	"encoding/json"
	"net/http"
	"protokey/pkg/app/model"
)

type Handler interface {
	Set(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Keys(w http.ResponseWriter, r *http.Request)
}

func NewHandler(storage *model.Storage) Handler {
	return &handler{storage: storage}
}

type handler struct {
	storage *model.Storage
}

func (h *handler) Set(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ||
		!model.ValidKey.MatchString(req.Key) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.storage.Set(req.Key, req.Value); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if !model.ValidKey.MatchString(key) {
		http.Error(w, "Invalid key", http.StatusBadRequest)
		return
	}
	val, err := h.storage.Get(key)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(map[string]string{"value": val})
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
}

func (h *handler) Keys(w http.ResponseWriter, r *http.Request) {
	prefix := r.URL.Query().Get("prefix")
	if !model.ValidKey.MatchString(prefix) {
		http.Error(w, "Invalid prefix", http.StatusBadRequest)
		return
	}

	keys, err := h.storage.Keys(prefix)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(map[string][]string{"keys": keys})
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
}
