package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

func newLogger(f flags) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if f.debug {
		opts.Level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}

func main() {
	flags := parseFlags()

	logger := newLogger(flags)

	app := newApp(flags.appPath, flags.proxyURL, flags.run, flags.buildCmd)

	proxyURL := flags.proxyURL
	if flags.templ {
		proxyURL = "localhost:7331" // TODO: add flag for this
	}

	timeout := time.Duration(flags.connectTimeout) * time.Second
	proxy, err := newProxy(app, logger, timeout, proxyURL)
	if err != nil {
		log.Fatalln(err)
	}

	logger.Info("Listening on " + flags.proxyBind)
	logger.Info("Proxying requests to " + flags.proxyURL)

	if flags.templ {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go watchTempl(proxy, flags.appPath, flags.proxyURL, &wg)
		wg.Wait()
		logger.Info("Templ proxy is ready")
	}

	logger.Info("Starting file watcher")
	watcher, err := newWatcher(logger, flags.appPath)
	if err != nil {
		log.Fatalln(err)
	}

	go notifyOnChange(proxy, watcher, logger)

	if err := http.ListenAndServe(flags.proxyBind, proxy); err != nil {
		log.Fatalln(err)
	}
}
