//
// This file is part of mini-httpd.
//
// Copyright (c) 2021 Alexander Metzner.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package main contains the mini-httpd application.
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/halimath/kvlog"
)

var (
	version        = "0.1.0"
	buildTimestamp = "now"
	commit         = "local"

	docRoot        = flag.String("doc-root", ".", "Document root")
	httpAddress    = flag.String("http-address", "localhost:8080", "Network address to bind to listen for incoming requests")
	disableLogging = flag.Bool("no-log", false, "Disable logging")
	cache          = flag.Bool("enable-caching", false, "Enable caching of unchanged files")
)

func main() {
	flag.Parse()

	fmt.Printf("mini-httpd v%s (%s; built %s)\nhttps://github.com/halimath/mini-httpd\nPublished under the Apache License version 2.\n\n", version, commit, buildTimestamp)

	docRootAbs, err := filepath.Abs(*docRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to determine absolute path for %s: %s\n", os.Args[0], *docRoot, err)
		os.Exit(1)
	}

	var handler http.Handler
	handler = http.FileServer(http.Dir(docRootAbs))

	fmt.Printf("%20s: %s\n", "Docroot", docRootAbs)
	fmt.Printf("%20s: %s\n", "Address", *httpAddress)

	if !*cache {
		handler = noCache(handler)
		fmt.Printf("%20s: disabled\n", "Cache")
	}

	if !*disableLogging {
		handler = kvlog.Middleware(kvlog.L, handler)
	}

	fmt.Println()

	server := &http.Server{
		Addr:    *httpAddress,
		Handler: handler,
	}

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-termChan
		fmt.Printf("Got SIGTERM. Shutting down...\n")
		server.Shutdown(context.Background())
	}()

	err = server.ListenAndServe()

	if err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error running mini-httpd: %s", err)
	}
}

var noCacheHeadersToAdd = map[string]string{
	"Expires":         time.Unix(0, 0).Format(time.RFC1123),
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var noCacheHeadersToRemove = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func noCache(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, h := range noCacheHeadersToRemove {
			r.Header.Del(h)
		}

		for h, v := range noCacheHeadersToAdd {
			w.Header().Set(h, v)
		}

		handler.ServeHTTP(w, r)
	})
}
