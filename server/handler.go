package server

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/kofuk/spaghettini/server/types"
)

type Handler struct {
	logger    *slog.Logger
	evaluator *Evaluator
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Request received", slog.String("method", r.Method), slog.String("path", r.URL.Path))

	body, err := io.ReadAll(io.LimitReader(r.Body, 4096))
	if err != nil {
		h.logger.Error("failed to read request body", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if len(body) == 4096 {
		h.logger.Warn("request body is too large")
		http.Error(w, "request body is too large", http.StatusRequestEntityTooLarge)
		return
	}
	r.Body.Close()

	request := &types.Request{
		Method:  r.Method,
		Path:    r.URL.Path,
		Header:  r.Header,
		Body:    body,
		Trailer: r.Trailer,
	}
	response, err := h.evaluator.Evaluate(request)
	if err != nil {
		h.logger.Error("failed to execute template", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	for key, values := range response.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	for key, values := range response.Trailer {
		for _, value := range values {
			w.Header().Add(http.TrailerPrefix+key, value)
		}
	}

	w.WriteHeader(response.StatusCode)
	w.Write([]byte(response.Body))
}
