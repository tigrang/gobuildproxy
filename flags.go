package main

import "flag"

type flags struct {
	appPath        string
	run            string
	buildCmd       string
	proxyURL       string
	proxyBind      string
	connectTimeout int
	debug          bool
	templ          bool
}

func parseFlags() flags {
	var f flags

	flag.StringVar(&f.appPath, "path", ".", "Path to app")
	flag.StringVar(&f.run, "run", "./run", "path to script that will run app")
	flag.StringVar(&f.buildCmd, "build", "./build", "path to script that will build app")
	flag.StringVar(&f.proxyBind, "proxybind", "localhost:9000", "the addr for error proxy to listen on")
	flag.StringVar(&f.proxyURL, "proxy", "localhost:3000", "url app is listening on to forward requests")
	flag.IntVar(&f.connectTimeout, "timeout", 30, "the number of seconds to wait for proxy to be available")
	flag.BoolVar(&f.debug, "debug", false, "Enable debug mode")
	flag.BoolVar(&f.templ, "templ", false, "run templ generate and capture error logs")
	flag.Parse()

	return f
}
