package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func newLogger(debug bool) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if debug {
		opts.Level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}

func main() {
	flags := parseFlags()

	logger := newLogger(flags.debug)

	app := newApp(flags.appPath, flags.proxyUrl, flags.run, flags.buildCmd)

	proxy, err := newProxy(
		flags.proxyBind,
		flags.notifyRoute,
		time.Duration(flags.connectTimeout)*time.Second,
		app,
	)

	if err != nil {
		log.Fatalln(err)
	}

	if flags.notify && flags.watch {
		log.Fatalln("--notify and --watch cannot be used at the same time")
	}

	if flags.notify {
		if err := proxy.notify(); err != nil {
			log.Fatalln(err)
		}
		return
	}

	logger.Info("Listening on " + flags.proxyBind)
	logger.Info("Proxy to " + flags.proxyUrl)

	if flags.watch {
		watcher, err := newWatcher(logger)
		if err != nil {
			log.Fatalln(err)
		}

		go notifyOnChange(proxy, watcher, logger)
	}

	if err := http.ListenAndServe(flags.proxyBind, proxy); err != nil {
		log.Fatalln(err)
	}
}
