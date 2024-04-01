// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Methods used in nowasm build
//go:build !wasm

package client

import (
	"errors"
	"fmt"
	"net"

	"github.com/teonet-go/tru"
)

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
	// Convert reader in parameters to tru.ReaderFunc
	for i, p := range params {
		if r, ok := p.(func(ch *Channel, pac *Packet, err error) bool); ok {
			reader := func(ch *tru.Channel, pac *tru.Packet, err error) bool {
				return r(&Channel{ch}, &Packet{pac}, err)
			}
			params[i] = reader
		}
	}
	t.Tru, err = tru.New(port, append(params, tru.MaxDataLenType(1280))...)
	return
}

type Channel struct{ *tru.Channel }

func (t *Tru) Connect(addr string, reader ...ReaderFunc) (ch *Channel, err error) {

	var r []tru.ReaderFunc
	if len(reader) > 0 {
		r = append(r, func(ch *tru.Channel, pac *tru.Packet, err error) bool {
			return reader[0](&Channel{ch}, &Packet{pac}, err)
		})
	}

	c, err := t.Tru.Connect(addr, r...)
	if err != nil {
		return
	}
	ch = &Channel{c}
	return
}

// WriteToCh write data to channel by address
func (t *Tru) WriteToCh(data []byte, addr string) (id int, err error) {

	t.ForEachChannel(func(ch *tru.Channel) {
		if ch.Addr().String() == addr {
			id, err = ch.WriteTo(data)
		}
	})

	return
}

func (ch *Channel) WriteTo(data []byte, delivery ...interface{}) (id int, err error) {

	if _, ok := ch.Addr().(*net.UDPAddr); !ok {
		// if ch.Channel.Addr().String() == "" {
		err = fmt.Errorf("looks like webrtc channel, skip it")
		fmt.Println(err)
		return
	}

	return ch.Channel.WriteTo(data, delivery...)
}

func (*Tru) ErrChannelDestroyed(err error) bool {
	return errors.Is(err, tru.ErrChannelDestroyed)
}

type Packet struct{ *tru.Packet }
