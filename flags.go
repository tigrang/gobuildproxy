package main

import "flag"

type flags struct {
	notify         bool
	run            string
	proxyUrl       string
	proxyBind      string
	notifyRoute    string
	buildCmd       string
	connectTimeout int
	appPath        string
	watch          bool
	debug          bool
}

func parseFlags() flags {
	var f flags

	flag.BoolVar(&f.notify, "notify", false, "notify proxy to trigger build")
	flag.StringVar(&f.run, "run", "./run", "path to script that will run app")
	flag.StringVar(&f.proxyBind, "proxybind", "localhost:9000", "the addr for error proxy to listen on")
	flag.StringVar(&f.notifyRoute, "notifyroute", "/internal/build/notify", "path to trigger builds (must be the same when --notify is used)")
	flag.StringVar(&f.proxyUrl, "proxy", "localhost:3000", "url app is listening on to forward requests")
	flag.StringVar(&f.buildCmd, "build", "./build", "path to script that will build app")
	flag.IntVar(&f.connectTimeout, "timeout", 30, "the number of seconds to wait for proxy to be available")
	flag.StringVar(&f.appPath, "path", "", "Path to app")
	flag.BoolVar(&f.watch, "watch", false, "Enable watch mode")
	flag.BoolVar(&f.debug, "debug", false, "Enable debug mode")
	flag.Parse()

	return f
}
