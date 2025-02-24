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

func newWatcher(logger *slog.Logger, appPath string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if appPath == "" {
		appPath, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	path, err := filepath.Abs(appPath)
	if err != nil {
		return nil, err
	}

	if err := watcher.Add(path); err != nil {
		log.Fatal(err)
	}

	if err := filepath.WalkDir(path, recursiveWatcher(watcher, logger, path)); err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return watcher, nil
}

func recursiveWatcher(watcher *fsnotify.Watcher, logger *slog.Logger, base string) func(string, fs.DirEntry, error) error {
	baseHiddenPath := base + string(filepath.Separator) + "."
	return func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return nil
		}

		if strings.Contains(path, "node_modules") || strings.HasPrefix(path, baseHiddenPath) {
			return nil
		}

		logger.Debug("watching path", slog.String("path", path))
		if err := watcher.Add(path); err != nil {
			logger.Warn("error adding watcher", slog.String("path", path), slog.String("error", err.Error()))
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

			if event.Op == fsnotify.Chmod {
				continue
			}

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

			if strings.HasSuffix(event.Name, "_templ.go") {
				continue
			}

			if strings.HasSuffix(event.Name, ".go") {
				logger.Debug("fsnotify event", slog.Any("event", event))
				proxy.clearErrorState()
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.Error("watcher error", slog.Any("error", err))
		}
	}
}
