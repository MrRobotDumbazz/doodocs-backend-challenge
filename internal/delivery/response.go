package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type Response struct {
	Status string `json:"status"`
	Data   interface{}
}

type contextKey struct {
	name string
}

var StatusCtxKey = contextKey{"Status"}

func Status(r *http.Request, status int) {
	*r = *r.WithContext(context.WithValue(r.Context(), StatusCtxKey, status))
}

func (h *Handler) EncodeJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	resp := Response{
		Status: strconv.Itoa(status) + http.StatusText(status),
		Data:   data,
	}
	if err := enc.Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if status, ok := r.Context().Value(StatusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	if _, err := w.Write(buf.Bytes()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
