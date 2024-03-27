// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Methods used in wasm build
//go:build wasm

package client

import (
	"net"
	"sync"
	"syscall/js"
	"time"
)

type Tru struct {
	global js.Value
	teoweb js.Value
	uuid   string
}

// TODO:
func New(port int, params ...interface{}) (t *Tru, err error) {

	t = new(Tru)

	// Define js variables
	t.global = js.Global()
	t.teoweb = t.global.Get("teoweb")

	// Get uuid addres from js function
	t.uuid = t.global.Call("uuidv4").String()

	// Add reader
	for _, p := range params {
		if _, ok := p.(func(ch *Channel, pac *Packet, err error) bool); ok {
			t.teoweb.Call("addReader", js.FuncOf(func(this js.Value, args []js.Value) any {
				// TODO: args[1] is js.Value with []byte
				// r(&Channel{}, &Packet{data: args[1]}, nil)
				return nil
			}))
			break
		}
	}

	return t, nil
}

// TODO:
func (t *Tru) Connect(addr string, reader ...ReaderFunc) (ch *Channel, err error) {

	// Connect to Teonet WebRTC server
	const url = "wss://signal.teonet.dev/signal"
	const server = "server-1"
	t.teoweb.Call("connect", url, t.uuid, server)

	if len(reader) > 0 {
		t.teoweb.Call("addReader", js.FuncOf(func(this js.Value, args []js.Value) any {
			// TODO: args[1] is js.Value with []byte
			// reader[0](&Channel{}, &Packet{data: args[1]}, nil)
			return nil
		}))
	}

	return nil, nil
}

// TODO:
func (t *Tru) LocalPort() int {
	return 0
}

func (t *Tru) ErrChannelDestroyed(err error) bool {
	return true
}

type Channel struct {
}

// TODO:
func (c *Channel) Addr() (addr net.Addr) {
	return
}

// TODO:
func (c *Channel) WriteTo(data []byte) (int, error) {
	return 0, nil
}

// TODO:
func (c *Channel) Close() {

}

// TODO:
func (c *Channel) String() string {
	return ""
}

// Packet struct (copy of tru.Packet)
type Packet struct {
	id                 uint32             // Packet ID
	status             uint8              // Packet Type
	data               []byte             // Packet Data
	time               time.Time          // Packet creating time
	retransmitTime     time.Time          // Packet retransmit time
	retransmitAttempts int                // Packet retransmit attempts
	delivery           PacketDeliveryFunc // Packet delivery callback function
	deliveryTimeout    time.Duration      // Packet delivery callback timeout
	deliveryTimer      time.Timer         // Packet delivery timeout timer
	sync.RWMutex
}
type PacketDeliveryFunc func(pac *Packet, err error)

func (p *Packet) Data() []byte {
	return p.data
}
