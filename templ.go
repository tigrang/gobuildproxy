package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
)

func watchTempl(proxy *handler, appDir string, proxyURL string, wg *sync.WaitGroup) {
	cmd := exec.Command(
		"templ",
		"generate",
		"--open-browser=false",
		"--watch",
		"--proxy",
		"http://"+proxyURL,
		"cmd",
		"echo 'gobuildproxy: rebuild'",
	)
	cmd.Dir = appDir

	r, w := io.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w

	go func() {
		defer w.Close()
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatalf("templ error: %v", err)
		}
	}()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "gobuildproxy: rebuild" || strings.Contains(line, "Error cleared") {
			proxy.clearErrorState()
		} else if strings.HasPrefix(line, "(✗)") {
			proxy.captureError("templ", line)
		} else if strings.HasPrefix(line, "(✓) Proxying") {
			wg.Done() // wait for templ to be ready
		}
	}
}
