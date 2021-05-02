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
	"syscall"

	"github.com/halimath/kvlog"
)

const (
	version = "0.1.0"
)

var (
	docRoot        = flag.String("doc-root", ".", "Document root")
	httpAddress    = flag.String("http-address", ":8080", "Network address to bind to listen for incoming requests")
	disableLogging = flag.Bool("no-log", false, "Disable logging")
)

func main() {
	fmt.Printf("mini-httpd v%s https://github.com/halimath/mini-httpd\nPublished under the Apache License version 2.\nServing at %s\n\n", version, *httpAddress)

	var handler http.Handler
	handler = http.FileServer(http.Dir(*docRoot))

	if !*disableLogging {
		handler = kvlog.Middleware(kvlog.L, handler)
	}

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

	err := server.ListenAndServe()

	if err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error running mini-httpd: %s", err)
	}
}
