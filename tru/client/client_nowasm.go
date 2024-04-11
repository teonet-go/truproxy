// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Methods used in nowasm build
//go:build !wasm

package client

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net"

	"github.com/teonet-go/teogw"
	w "github.com/teonet-go/teowebrtc_server"
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
				var dc w.DataChannel
				return r(&Channel{ch, dc}, &Packet{pac}, err)
			}
			params[i] = reader
		}
	}
	t.Tru, err = tru.New(port, append(params, tru.MaxDataLenType(1280))...)
	return
}

func (t *Tru) Connect(addr, peer string, reader ...ReaderFunc) (ch *Channel, 
	err error) {

	var r []tru.ReaderFunc
	if len(reader) > 0 {
		r = append(r, func(ch *tru.Channel, pac *tru.Packet, err error) bool {
			var dc w.DataChannel
			return reader[0](&Channel{ch, dc}, &Packet{pac}, err)
		})
	}

	c, err := t.Tru.Connect(addr, r...)
	if err != nil {
		return
	}
	var dc w.DataChannel
	ch = &Channel{c, dc}
	return
}

// WriteToCh write data to channel by address.
func (t *Tru) WriteToCh(data []byte, addr string) (id int, err error) {

	// The WriteToCh is not implemented in nowasm version. This function is not
	// used in client mode.

	// t.ForEachChannel(func(ch *tru.Channel) {
	// 	if ch.Addr().String() == addr {
	// 		id, err = ch.WriteTo(data)
	// 	}
	// })

	return
}

func (*Tru) ErrChannelDestroyed(err error) bool {
	return errors.Is(err, tru.ErrChannelDestroyed)
}

type Channel struct {
	*tru.Channel
	w.DataChannel
}

func (ch *Channel) WriteTo(data []byte, delivery ...interface{}) (id int, err error) {

	// Send to webrtc channel
	if _, ok := ch.Addr().(*net.UDPAddr); !ok {
		err = fmt.Errorf("looks like webrtc channel, try send to webrtc channel")
		fmt.Println(err)

		var gw teogw.TeogwData
		data = []byte(base64.StdEncoding.EncodeToString(data))
		data, err = gw.MarshalJson(gw, "data", data, nil)
		if err != nil {
			return
		}
		err = ch.DataChannel.Send(data)

		return
	}

	// Send to tru channel
	return ch.Channel.WriteTo(data, delivery...)
}

func (ch *Channel) String() string {
	if _, ok := ch.Addr().(*net.UDPAddr); ok {
		return ch.Addr().String()
	}
	if s, ok := ch.DataChannel.GetUser().(string); ok {
		return s
	}

	return ""
}

func (ch *Channel) Close() {

	// Close tru channel
	if _, ok := ch.Addr().(*net.UDPAddr); ok {
		ch.Channel.Close()
		return
	}

	// Close webrtc channel
	if _, ok := ch.DataChannel.GetUser().(string); ok {
		// ch.DataChannel.Close()
		// TODO: close webrtc channel
		return
	}
}

type Packet struct{ *tru.Packet }
