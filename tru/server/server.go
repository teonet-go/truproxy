// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Server package contains the server implementation for the Truproxy.
package server

import (
	"log"
	"time"

	"github.com/teonet-go/teogw"
	webrtc "github.com/teonet-go/teowebrtc_server"
	"github.com/teonet-go/tru"
)

const defaultSignalServer = "signal.teonet.dev"

// init initializes the Go program.
//
// It sets the log flags to include the standard date and time format, as well
// as microseconds.
func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

// TruServer is the main truproxy server type that contains the core components.
// It has mutexes for synchronization, the teowebrtc server,
// the Tru client, and API clients. As an exported type, it is part of the
// public API.
type TruServer struct {
	*webrtc.WebRTC
	*tru.Tru
}

// New creates a new TruServer instance.
func New(appName, appVersion string, appStart time.Time, sigAddr string,
	ownSign bool, name string,
	port int, params ...interface{}) (t *TruServer, err error) {

	t = &TruServer{}

	// Use default teonet signal server if not specified
	if sigAddr == "" {
		sigAddr = defaultSignalServer
	}

	// Create WebRTC server
	t.WebRTC, err = webrtc.New(
		sigAddr, // signal server address
		ownSign, // connect to remote signal server if false or create own if true
		name,    // Name of this server
		new(teogw.TeogwData).MarshalJson,
		new(teogw.TeogwData).UnmarshalJson,
		func(peer string, dc webrtc.DataChannel) {},
		func(peer string, dc webrtc.DataChannel) {},
	)
	if err != nil {
		return nil, err
	}
	t.webrtcCommands(appName, appVersion, appStart)

	// Create Tru server
	t.Tru, err = tru.New(port, params...)
	if err != nil {
		return nil, err
	}

	return
}
