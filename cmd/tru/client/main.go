// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Teonet fortune web-server microservice. This is simple Teonet web-server
// micriservice application which get fortune message from Teonet Fortune
// microservice and show it in the site web page.
package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/teonet-go/teonet"
)

const (
	appShort = "teoproxy-client"
	appName  = "Teonet truproxy client web server"
	appLong  = `This is simple <a href="https://github.com/teonet-go">Teonet</a> web-server. Here are its main functions:
		- it get message from web form and send it to server using teowebrtc;
		- it send this message from wasm code using truproxy client.`
	appVersion = "0.0.1"

	appPort = "8080"
)

var appStartTime = time.Now()
var domain /* fortune, */, monitor string

// Params is teonet command line parameters
var Params struct {
	appShort    string
	port        int
	httpAddr    string
	stat        bool
	hotkey      bool
	showPrivate bool
	loglevel    string
	logfilter   string
}

func main() {

	// Application logo
	teonet.Logo(appName, appVersion)

	// Get HTTP port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = appPort
	}

	// Parse application command line parameters
	flag.StringVar(&Params.appShort, "name", appShort, "application short name")
	flag.IntVar(&Params.port, "p", 0, "local port")
	flag.StringVar(&Params.httpAddr, "addr", ":"+port, "http server local address")
	flag.BoolVar(&Params.stat, "stat", false, "show statistic")
	flag.BoolVar(&Params.hotkey, "hotkey", false, "start hotkey menu")
	flag.BoolVar(&Params.showPrivate, "show-private", false, "show private key")
	flag.StringVar(&Params.loglevel, "loglevel", "debug", "log level")
	flag.StringVar(&Params.logfilter, "logfilter", "", "log filter")
	//
	flag.StringVar(&domain, "domain", "", "domain name to process HTTP/s server")
	// flag.StringVar(&fortune, "fortune", "", "fortune microservice address")
	flag.StringVar(&monitor, "monitor", "", "monitor address")
	//
	flag.Parse()

	// Get fortune address from environment variable
	// if len(fortune) == 0 {
	// 	fortune = os.Getenv("TEO_FORTUNE")
	// }

	// Check requered parameters
	// teonet.CheckRequeredParams("fortune")

	// Initialize and run Teonet
	teo, err := newTeonet()
	if err != nil {
		log.Panic(err)
		return
	}

	// Initialize and run web-server
	err = newServe(domain, appLong, Params.httpAddr, teo)
	if err != nil {
		log.Panic(err)
		return
	}
}
