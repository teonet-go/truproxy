// Copyright 2023-2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The truproxy example web server package.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/teonet-go/truproxy/tru/server"
	"golang.org/x/crypto/acme/autocert"
)

const (
	appShort   = "truproxy-server"
	appName    = "Teonet truproxy server"
	appLong    = ``
	appVersion = "0.0.1"
)

var appStart = time.Now()

var domain string

// main is the entry point of the program.
//
// It parses application parameters, defines a handler function for HTTP requests,
// registers the handler function to handle all requests, creates a file server
// to serve static files, registers a teowebrtc server, starts an HTTPS server
// if a domain is set, or starts an HTTP server if a domain is not set.
func main() {

	// Print
	fmt.Printf("%s, ver. %s\n", appName, appVersion)

	// Parse application parameters
	var monitor, laddr string
	var gzip bool
	//
	flag.StringVar(&domain, "domain", "", "domain name to process HTTP/s server")
	flag.StringVar(&laddr, "laddr", "localhost:8085", "local address of http, used if domain doesn't set")
	flag.StringVar(&monitor, "monitor", "", "teonet monitor address")
	flag.BoolVar(&gzip, "gzip", false, "gzip http files")
	flag.Parse()

	// Define Hello handler function for the HTTP requests
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	}

	// Register the Hello handler function to handle all requests
	http.HandleFunc("/hello", handler)

	// Create a file server to serve static files from the "frontend" directory
	var frontendFS http.Handler
	if gzip {
		frontendFS = gziphandler.GzipHandler(http.FileServer(http.FS(getFrontendAssets())))
	} else {
		frontendFS = http.FileServer(http.FS(getFrontendAssets()))
	}
	http.Handle("/", frontendFS)

	// Start tru proxy server
	_, err := server.New(appName, appVersion, appStart, "", "server-2")
	if err != nil {
		fmt.Println("Create teonet proxy server error:", err)
		return
	}
	// http.HandleFunc("/ws", serve.HandleWebSocket)

	// Start HTTPS server if domain is set
	if len(domain) > 0 {

		// Redirect HTTP requests to HTTPS
		go func() {
			err := http.ListenAndServe(":80", http.HandlerFunc(redirectTLS))
			if err != nil {
				log.Fatalf("ListenAndServe error: %v", err)
			}
		}()

		// Start HTTPS server and create certificate for domain
		log.Println("Start https serve with domain:", domain)
		log.Fatal(http.Serve(autocert.NewListener(domain), nil))
		return
	}

	// Start HTTP server
	log.Println("Start http server at:", laddr)
	log.Fatalln(http.ListenAndServe(laddr, nil))
}

// redirectTLS redirects the HTTP request to HTTPS.
//
// It takes in the http.ResponseWriter and http.Request as parameters.
func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+domain+":443"+r.RequestURI,
		http.StatusMovedPermanently)
}
