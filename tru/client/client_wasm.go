// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Methods used in wasm build
//go:build wasm

package client

import (
	"fmt"
	"sync"
	"syscall/js"
	"time"
)

type Tru struct {
	global js.Value
	teoweb js.Value
	uuid   string
}

// New creates new Tru instance.
func New(port int, params ...interface{}) (t *Tru, err error) {

	t = new(Tru)

	// Define js variables
	t.global = js.Global()
	t.teoweb = t.global.Get("teoweb")

	// Get uuid addres from js function
	t.uuid = t.global.Call("uuidv4").String()

	// Add reader
	for _, p := range params {
		if r, ok := p.(func(ch *Channel, pac *Packet, err error) bool); ok {
			fmt.Println("add reader")
			t.teoweb.Call("addReader", js.FuncOf(
				func(this js.Value, args []js.Value) any {
					// fmt.Println("got in reader, command:", args[0].Get("command").String())
					switch args[0].Get("command").String() {
					case "name":
						t.global.Call("setIdText", "name", args[1])
					case "clients":
						t.global.Call("setIdText", "clients", args[1])
					case "version":
						t.global.Call("setIdText", "version", args[1])
					case "uptime":
						t.global.Call("setIdText", "uptime", args[1])
					case "data":
						if args[1].IsNull() {
							// fmt.Println("got empty data command")
							return nil
						}
						go r(&Channel{t: t}, &Packet{data: []byte(args[1].String())}, nil)
					}
					return nil
				}))
			break
		}
	}

	// On connect
	t.teoweb.Call("onOpen", js.FuncOf(func(this js.Value, args []js.Value) any {
		fmt.Println("onOpen")
		t.global.Call("setIdText", "online", true)
		t.teoweb.Call("sendCmd", "clients")
		t.teoweb.Call("subscribeCmd", "clients")
		t.teoweb.Call("sendCmd", "name")
		t.teoweb.Call("sendCmd", "uptime")
		t.teoweb.Call("sendCmd", "version")
		return nil
	}))

	// On disconnect
	t.teoweb.Call("onClose", js.FuncOf(func(this js.Value, args []js.Value) any {
		fmt.Println("onClose")
		t.global.Call("setIdText", "online", false)
		return nil
	}))

	return t, nil
}

// TODO:
func (t *Tru) Connect(addr string, reader ...ReaderFunc) (ch *Channel, err error) {

	fmt.Println("Connect", addr)

	// Connect to Teonet WebRTC server
	const url = "wss://signal.teonet.dev/signal"
	const peer = "tloop-server-1"
	t.teoweb.Call("connect", url, t.uuid, peer)

	if len(reader) > 0 {
		t.teoweb.Call("addReader", js.FuncOf(func(this js.Value, args []js.Value) any {
			// TODO: args[1] is js.Value with []byte
			// reader[0](&Channel{}, &Packet{data: args[1]}, nil)
			return nil
		}))
	}

	// Wait connected
	var wait = make(chan any)
	t.teoweb.Call("whenConnected", js.FuncOf(func(this js.Value, args []js.Value) any {
		wait <- nil
		return nil
	}))
	<-wait

	// TODO: Add connected channel to channels map

	ch = &Channel{t: t, peer: peer}

	return
}

// TODO:
func (t *Tru) LocalPort() int {
	return 0
}

func (t *Tru) ErrChannelDestroyed(err error) bool {
	return true
}

type Channel struct {
	t    *Tru   // Pointer to Tru
	peer string // Peer name
}

// TODO:
func (c *Channel) Addr() (addr string) {
	return c.peer
}

// TODO:
func (c *Channel) WriteTo(data []byte) (int, error) {
	fmt.Println("WriteTo", data)
	c.t.teoweb.Call("sendCmd", "data", string(data))
	return 0, nil
}

// TODO:
func (t *Tru) WriteToCh(data []byte, addr string) (int, error) {

	// The WriteToCh is not implemented in wasm version. This function is not
	// used in client mode.

	return 0, nil
}

// TODO:
func (c *Channel) Close() {
	c.t.teoweb.Call("close")
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
