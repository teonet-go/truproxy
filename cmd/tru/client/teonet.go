// Copyright 2022 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Teonet module

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/teonet-go/teomon"
	"github.com/teonet-go/teonet"
)

type Teonet struct {
	*teonet.Teonet
	// fortune *teonet.APIClient
}

// newTeonet create new Teonet
func newTeonet() (teo *Teonet, err error) {

	teo = new(Teonet)

	// Create new teonet connector
	teo.Teonet, err = teonet.New(Params.appShort, Params.port,
		teonet.Log(), Params.loglevel, teonet.Logfilter(Params.logfilter),
		teonet.Stat(Params.stat), teonet.Hotkey(Params.hotkey),
	)
	if err != nil {
		teo.Log().Error.Println("can't init Teonet, error:", err)
		return
	}

	// Show this application private key
	if Params.showPrivate {
		fmt.Printf("%x\n", teo.GetPrivateKey())
		os.Exit(0)
	}

	// Connect to teonet
	for teo.Connect() != nil {
		teo.Log().Error.Println("can't connect to Teonet, try again...")
		time.Sleep(1 * time.Second)
	}

	// Print teonet address
	// teo.Log().Debug.Println("connected to teonet, address:", teo.Address())

	// Connet to fortune
	// if len(fortune) == 0 {
	// 	err = errors.New("can't connect to 'fortune', err: fortune address not set")
	// 	return
	// }
	// for i := 1; ; i++ {
	// 	err = teo.ConnectTo(fortune)
	// 	if err == nil {
	// 		break
	// 	}
	// 	err = errors.New("can't connect to 'fortune', err: " + err.Error())
	// 	teo.Log().Error.Println(err)
	// 	if i >= 5 {
	// 		return
	// 	}
	// }

	// Connet to fortune api
	// if teo.fortune, err = teo.NewAPIClient(fortune); err != nil {
	// 	err = errors.New("can't connect to 'fortune' api, err: " + err.Error())
	// 	return
	// }

	// Connect to monitor
	if len(monitor) == 0 {
		return
	}
	teomon.Connect(teo.Teonet, monitor, teomon.Metric{
		AppName:      appName,
		AppShort:     Params.appShort,
		AppVersion:   appVersion,
		TeoVersion:   teonet.Version,
		AppStartTime: appStartTime,
	})
	teo.Log().Debug.Println("connected to monitor")

	return
}

// Fortune get Fortune messsage from teofortune microservice
func (teo *Teonet) Fortune() (msg string, err error) {

	// Get fortune message from teofortune microservice
	// id, err := teo.fortune.SendTo("fortb", nil)
	// if err != nil {
	// 	return
	// }
	// data, err := teo.WaitFrom(fortune, uint32(id))
	// if err != nil {
	// 	return
	// }

	// msg = string(data)

	return
}

// Version return teonet version
func (teo *Teonet) Version() (msg string) {
	return teonet.Version
}
