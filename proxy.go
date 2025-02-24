package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type handler struct {
	mu             sync.Mutex
	app            *app
	logger         *slog.Logger
	proxy          *httputil.ReverseProxy
	connectTimeout time.Duration
}

func newProxy(app *app, logger *slog.Logger, connectTimeout time.Duration, proxyURL string) (*handler, error) {
	remote, err := url.Parse("http://" + proxyURL)
	if err != nil {
		return nil, err
	}

	return &handler{
		app:            app,
		logger:         logger,
		proxy:          httputil.NewSingleHostReverseProxy(remote),
		connectTimeout: connectTimeout,
	}, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Serving request", slog.Any("url", r.URL), slog.String("method", r.Method))

	h.mu.Lock()
	err := h.app.rebuildIfDirty(h.connectTimeout)
	h.mu.Unlock()

	if err != nil {
		h.logger.Error(err.Error())
		h.respondWithError(w, err)
		return
	}

	h.proxy.ServeHTTP(w, r)
}

func (h *handler) captureError(tool string, line string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.app.buildOutput = line
	h.app.lines = parse(line, h.app.path)
	h.app.lastErr = fmt.Errorf("%s build error", tool)
	h.app.markAsDirty()
}

func (h *handler) clearErrorState() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.app.buildOutput = ""
	h.app.lastErr = nil
	h.app.lines = nil
	h.app.markAsDirty()
}

func (h *handler) respondWithError(w http.ResponseWriter, err error) {
	if err := tmpl.Execute(w, map[string]any{"lines": h.app.lines, "error": err.Error()}); err != nil {
		h.logger.Error("Failed to execute template", "err", err)
	}
}
