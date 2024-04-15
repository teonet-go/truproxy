// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Truproxy server WebRTC commands.

package server

import (
	"fmt"
	"time"

	w "github.com/teonet-go/teowebrtc_server"
)

// webrtcCommands adds WebRTC server commands.
func (t *TruServer) webrtcCommands(appName, appVersion string, appStart time.Time) {

	t.Commands.

		// Get hello (test command)
		Add("hello", func(dc w.DataChannel, gw w.WebRTCData) ([]byte, error) {
			return []byte("hello"), nil
		}).

		// Get hello-2 (test command)
		Add("hello-2", func(dc w.DataChannel, gw w.WebRTCData) ([]byte, error) {
			return []byte("hello-2"), nil
		}).

		// Get app name
		Add("name", func(dc w.DataChannel, gw w.WebRTCData) ([]byte, error) {
			return []byte(appName), nil
		}).

		// Get app version
		Add("version", func(dc w.DataChannel, gw w.WebRTCData) ([]byte, error) {
			return []byte(appVersion), nil
		}).

		// Get app uptime
		Add("uptime", func(dc w.DataChannel, gw w.WebRTCData) ([]byte, error) {
			// Time since app start
			d := time.Since(appStart)
			d = d.Round(time.Second)

			// Hours
			h := d / time.Hour
			d -= h * time.Hour

			// Minutes
			m := d / time.Minute
			d -= m * time.Minute

			// Seconds
			s := d / time.Second

			// Format time duration in hours, minutes and seconds
			uptime := fmt.Sprintf("%dh %02dm %02ds", h, m, s)
			return []byte(uptime), nil
		})
}
