// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Common methods for wasm and nowasm build

package client

type ReaderFunc func(ch *Channel, pac *Packet, err error) (processed bool)
