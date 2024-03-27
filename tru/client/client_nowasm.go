// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Methods used in nowasm build
//go:build !wasm

package client

import "github.com/teonet-go/tru"

type Tru struct {
	*tru.Tru
}

// New create new tru proxy object and start listen udp packets. Parameters by
// type:
//
//	int:                local port number, 0 for any
//	tru.ReaderFunc:     message receiver callback function
//	tru.ConnectFunc:    connect to server callback function
//	tru.PunchFunc:      punch callback function
//	*teolog.Teolog:     pointer to teolog
//	string:             loggers level
//	teolog.Filter:      loggers filter
//	tru.StartHotkey:    start hotkey meny
//	tru.ShowStat:       show statistic
//	tru.MaxDataLenType: max packet data length
func New(port int, params ...interface{}) (t *Tru, err error) {
	t = new(Tru)
	t.Tru, err = tru.New(port, append(params, tru.MaxDataLenType(1280))...)
	return
}
