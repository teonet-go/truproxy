// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Methods used in wasm build
//go:build wasm

package client

import (
	"encoding/base64"
	"fmt"
	"sync"
	"syscall/js"
	"time"

	"github.com/teonet-go/tru"
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
	var reader func(ch *Channel, pac *Packet, err error) bool
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
						if args[1].IsNull() || len(args[1].String()) == 0 {
							// Skip command without data
							return nil
						}
						s := args[1].String()
						d, err := base64.StdEncoding.DecodeString(s)
						if err != nil {
							fmt.Println("base64 decode error:", err, s)
							return nil
						}

						go r(&Channel{t: t}, &Packet{data: d}, nil)
					}
					return nil
				}),
			)
			reader = r
			break
		}
	}

	// On connect
	t.teoweb.Call("onOpen", js.FuncOf(func(this js.Value, args []js.Value) any {
		fmt.Println("client_wasm onOpen")
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
		fmt.Println("client_wasm onClose")
		t.global.Call("setIdText", "online", false)
		if reader != nil {
			reader(&Channel{t: t}, nil, tru.ErrChannelDestroyed)
		}
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
	t.teoweb.Call("connect", url, t.uuid, peer, false)

	// Common reader which sets in reader function argument
	if len(reader) > 0 {
		t.teoweb.Call("addReader", js.FuncOf(func(this js.Value, args []js.Value) any {
			// TODO: args[1] is js.Value with []byte
			// reader[0](&Channel{}, &Packet{data: args[1]}, nil)
			return nil
		}))
	}

	// Wait connected up to 10 seconds
	started := time.Now()
	for !t.teoweb.Call("connected").Bool() {
		if time.Since(started) > 10*time.Second {
			err = fmt.Errorf("timeout")
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Return connected channel
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
	data = []byte(base64.StdEncoding.EncodeToString(data))
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
	fmt.Println("---= wasm channel close =---")
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
