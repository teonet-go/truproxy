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
//
// Parameters:
// appName: the name of the application, used in the WebRTC server.
// appVersion: the version of the application, used in the WebRTC server.
// appStart: the start time of the application, used in the WebRTC server.
// sigAddr: the address of the signal server. If empty, the default
//          teonet signal server is used.
// ownSign: if true, the WebRTC server will create its own signal server,
//          otherwise it will connect to the remote signal server at sigAddr.
// name: the name of this server, used in the WebRTC server.
// connect: a callback function called when a new peer is connected, or nil.
// disconn: a callback function called when a peer is disconnected, or nil.
// port: the port number of the Tru server.
// params: additional parameters passed to the Tru server.
//
// Returns a new TruServer instance, or an error if the WebRTC server or the
// Tru server cannot be created.
func New(appName, appVersion string, appStart time.Time, sigAddr string,
	ownSign bool, name string,
	connect func(peer string, dc webrtc.DataChannel),
	disconn func(peer string, dc webrtc.DataChannel),
	port int, params ...interface{}) (t *TruServer, err error) {

	t = &TruServer{}

	// Use default teonet signal server if not specified
	if sigAddr == "" {
		sigAddr = defaultSignalServer
	}

	// Combine onOpenClose callbacks
	var onOpenClose []webrtc.OnOpenCloseType
	if connect != nil {
		onOpenClose = append(onOpenClose, connect)
	}
	if disconn != nil {
		onOpenClose = append(onOpenClose, disconn)
	}

	// Create WebRTC server
	if t.WebRTC, err = webrtc.New(
		sigAddr,                            // signal server address
		ownSign,                            // connect to remote signal server if false or create own if true
		name,                               // name of this server
		new(teogw.TeogwData).MarshalJson,   // marshal JSON data function
		new(teogw.TeogwData).UnmarshalJson, // unmarshal JSON data function
		onOpenClose...,                     // onOpenClose callback functions
	); err != nil {
		return
	}

	// Create Tru server
	if t.Tru, err = tru.New(port, params...); err != nil {
		return
	}

	// Register WebRTC commands
	t.webrtcCommands(appName, appVersion, appStart)

	return
}
