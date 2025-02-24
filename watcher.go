package main

import (
	"github.com/fsnotify/fsnotify"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func newWatcher(logger *slog.Logger) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := watcher.Add("."); err != nil {
		log.Fatal(err)
	}

	if err := filepath.WalkDir(".", recursiveWatcher(watcher, logger)); err != nil {
		return nil, err
	}

	return watcher, nil
}

func recursiveWatcher(watcher *fsnotify.Watcher, logger *slog.Logger) func(string, fs.DirEntry, error) error {
	return func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			if !strings.Contains(path, "node_modules") && !strings.HasPrefix(path, ".") {
				logger.Debug("watching path", slog.String("path", path))
				watcher.Add(path)
			}
		}
		return nil
	}
}

func notifyOnChange(proxy *handler, watcher *fsnotify.Watcher, logger *slog.Logger) {
	defer watcher.Close()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			slog.Debug("watcher event", slog.Any("event", event))
			if event.Has(fsnotify.Create) {
				info, err := os.Stat(event.Name)
				if err != nil {
					logger.Error("failed to stat", slog.Any("error", err))
					continue
				}

				if info.IsDir() {
					logger.Debug("directory created, adding to watcher")
					if err := watcher.Add(event.Name); err != nil {
						logger.Error("failed to add dir to watcher", slog.Any("error", err))
					}
				}
			}

			if strings.HasSuffix(event.Name, ".go") {
				logger.Debug("detected change, notifying")
				if err := proxy.notify(); err != nil {
					logger.Error("failed to notify proxy", slog.Any("error", err))
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.Error("watcher error", slog.Any("error", err))
		}
	}
}
