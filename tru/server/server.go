// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Server package contains the server implementation for the Truproxy.
package server

import (
	"log"
	"sync"

	"github.com/teonet-go/teowebrtc_server"
	"github.com/teonet-go/tru"
)

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
	*sync.Mutex
	*teowebrtc_server.WebRTC
	*tru.Tru
	apiClients *APIClients
}

// New creates a new TruServer instance. It initializes the mutex, API clients,
// Teonet client, and websocket server. The appShort parameter specifies the
// application name. The monitor parameter optionally configures connecting to a
// Teonet monitor for metrics reporting. It returns the TruServer instance
// and any error. As an exported function, this serves as the main constructor for
// the TruServer type.
func New(appShort string) (teo *TruServer, err error) {
	teo = &TruServer{Mutex: new(sync.Mutex)}

	// Init api clients object
	teo.initAPIClients()

	return
}
