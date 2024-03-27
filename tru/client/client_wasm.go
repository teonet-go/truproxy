// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Methods used in wasm build
//go:build wasm

package client

import (
	"sync"
	"time"

	"github.com/teonet-go/tru"
)

type Tru struct {
}

// TODO:
func New(port int, params ...interface{}) (t *Tru, err error) {
	return nil, nil
}

// TODO:
func (t *Tru) Connect(addr string, reader ...tru.ReaderFunc) (ch *Channel, err error) {
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
func (c *Channel) Addr() string {
	return ""
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
